--Luctus Monitor
--Made by OverlordAkise

--This script collects data about the server and players for analysis

LUCTUS_MONITOR_DEBUG = false

LUCTUS_MONITOR_URL = "http://localhost:7077/luastat"
LUCTUS_MONITOR_URL_EXTRA = "http://localhost:7077/luastatextra"
LUCTUS_MONITOR_URL_INIT = "http://localhost:7077/luajobinit"


function debugPrint(text)
    if LUCTUS_MONITOR_DEBUG then
        print(text)
    end
end

function debugPrintTable(tab)
    if LUCTUS_MONITOR_DEBUG then
        PrintTable(tab)
    end
end

function LuctusCreateUuid()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and math.random(0, 0xf) or math.random(8, 0xb)
        return string.format('%x', v)
    end)
end

LUCTUS_MONITOR_SERVER_ID = LuctusCreateUuid()
if file.Exists("data/luctus_monitor.txt","GAME") then
    print("[luctus_monitor] Found server ID, loading...")
    LUCTUS_MONITOR_SERVER_ID = file.Read("data/luctus_monitor.txt","GAME")
else
    print("[luctus_monitor] No ID found, creating...")
    file.Write("luctus_monitor.txt",LUCTUS_MONITOR_SERVER_ID)
end
print("[luctus_monitor] ServerID: ",LUCTUS_MONITOR_SERVER_ID)

--for getting id remote
concommand.Add("m_monitor", function(ply,cmd,args)
    ply:PrintMessage(HUD_PRINTCONSOLE,LUCTUS_MONITOR_SERVER_ID)
end)

util.AddNetworkString("luctus_monitor_collect")

timer.Create("luctus_monitor_autorestart",10,0,function()
    if not timer.Exists("luctus_monitor_timer") then
        print("[luctus_monitor] Starting Monitor timer")
        LuctusMonitorStart()
    end
end)

LUCTUS_MONITOR_PLAYERS = {}

function LuctusMonitorStart()
    timer.Create("luctus_monitor_timer",300,0,function()
        --This takes time, so run it first, then send it
        GetCurrentTickrate()
        LuctusMonitorCollectPlayers()
        
        timer.Simple(5,function()
            LuctusMonitorDo()
            LuctusMonitorDoExtras()
        end)
    end)
end


--Monitor deaths
local lm_deaths = 0
hook.Add("PostPlayerDeath","luctus_monitor_stat",function(ply)
    lm_deaths = lm_deaths + 1
end)

--Monitor Tickrate
LUCTUS_MONITOR_CURRENT_TICKRATE = 0
function GetCurrentTickrate()
    local abc = {}
    local LastCapture = 0
    local tickDelta = 0
    local lastTick = nil
    hook.Add("Tick", "average_tickrate_calc", function()
        local sysTime = SysTime()
        if not lastTick then lastTick = SysTime() return end
        tickDelta = sysTime - lastTick
        abc[#abc+1] = tickDelta
        if #abc >= 100 then
            local all = 0
            for k,v in pairs(abc) do
                all = all + v
            end
            LUCTUS_MONITOR_CURRENT_TICKRATE = 1/(all/#abc)
            hook.Remove("Tick", "average_tickrate_calc")
        end
        lastTick = sysTime
    end)
end

function LuctusMonitorDo()
    local data = {["players"] = {}}
    local server_avgfps = 0
    local server_avgfps_c = 0
    local server_avgping = 0
    local server_avgping_c = 0
    
    for k,v in pairs(LUCTUS_MONITOR_PLAYERS) do
        
        server_avgping = server_avgping + v.pingcur
        server_avgfps = server_avgfps + v.fpsavg
        server_avgping_c = server_avgping_c + 1
        server_avgfps_c = server_avgfps_c + 1
        if not player.GetBySteamID(k) then
            v.online = false --left already
        end
        table.insert(data["players"],v)
    end
    data["gamemode"] = engine.ActiveGamemode()
    data["map"] = game.GetMap()
    data["tickrateset"] = 1/engine.TickInterval()
    data["tickratecur"] = LUCTUS_MONITOR_CURRENT_TICKRATE
    data["entscount"] = #ents.GetAll()
    data["plycount"] = #player.GetAll()
    data["uptime"] = CurTime()
    data["serverid"] = LUCTUS_MONITOR_SERVER_ID
    if server_avgping_c == 0 then
        data["avgfps"] = 0
        data["avgping"] = 0
    else
        data["avgfps"] = server_avgfps/server_avgfps_c
        data["avgping"] = server_avgping/server_avgping_c
    end
    data["deaths"] = lm_deaths
    
    data["luaramb"] = collectgarbage("count")
    collectgarbage("collect")
    data["luarama"] = collectgarbage("count")
    
    local ret = HTTP({
        failed = function(failMessage)
            print("[luctus_monitor] FAILED TO POST STATS!")
            print("[luctus_monitor]",os.date("%H:%M:%S - %d/%m/%Y",os.time()))
            print(failMessage)
        end,
        success = function(httpcode,body,headers)
            debugPrint("[luctus_monitor] Do Sync successfull!")
        end, 
        method = "POST",
        url = LUCTUS_MONITOR_URL,
        body = util.TableToJSON(data),
        type = "application/json; charset=utf-8",
        timeout = 10
    })
    
    debugPrint("(Do) Sent the following:")
    debugPrint("Table:")
    debugPrintTable(data)
    debugPrint("Json:")
    debugPrint(util.TableToJSON(data))
    LUCTUS_MONITOR_PLAYERS = {}
    lm_deaths = 0
end

hook.Add("PlayerInitialSpawn","luctus_monitor_ply_init",function(ply)
    local jobname = ""
    if ply.getJobTable and ply:getJobTable() and ply:getJobTable().name then
        jobname = ply:getJobTable().name
    end
    LUCTUS_MONITOR_PLAYERS[ply:SteamID()] = {
        ["steamid"] = ply:SteamID(),
        ["nick"] = ply:Nick(),
        ["job"] = jobname,
        ["pingcur"] = ply:Ping(),
        ["pingavg"] = -1,
        ["fpsavg"] = -1,
        ["fpshigh"] = -1,
        ["fpslow"] = -1,
        ["packetslost"] = ply:PacketLoss(),
        ["luaramb"] = -1,
        ["luarama"] = -1,
        ["os"] = "",
        ["country"] = "",
        ["screensize"] = "",
        ["screenmode"] = "",
        ["jitver"] = "",
        ["ip"] = ply:IPAddress(),
        ["serverid"] = LUCTUS_MONITOR_SERVER_ID,
        ["playtime"] = 0,
        ["online"] = true
    }
    ply.lmonplaytime = CurTime()
end)

function LuctusMonitorCollectPlayers()
    net.Start("luctus_monitor_collect")
    net.Broadcast()
end

net.Receive("luctus_monitor_collect",function(len,ply)
    debugPrint("Got stats for new player:")
    debugPrint(ply:Nick().."//"..ply:SteamID())
    LUCTUS_MONITOR_PLAYERS[ply:SteamID()] = {
        ["steamid"] = ply:SteamID(),
        ["nick"] = ply:Nick(),
        ["job"] = DarkRP and ply:getJobTable().name or "",
        ["pingcur"] = ply:Ping(),
        ["pingavg"] = net.ReadInt(32),
        ["fpsavg"] = net.ReadInt(32),
        ["fpshigh"] = net.ReadInt(32),
        ["fpslow"] = net.ReadInt(32),
        ["packetslost"] = ply:PacketLoss(),
        ["luaramb"] = net.ReadFloat(),
        ["luarama"] = net.ReadFloat(),
        ["os"] = net.ReadString(),
        ["country"] = net.ReadString(),
        ["screensize"] = net.ReadString(),
        ["screenmode"] = net.ReadString(),
        ["jitver"] = net.ReadString(),
        ["ip"] = ply:IPAddress(),
        ["serverid"] = LUCTUS_MONITOR_SERVER_ID,
        ["playtime"] = math.Round(CurTime() - ply.lmonplaytime),
        ["online"] = true
    }
end)


--Extras (weaponkills, jobtime)

local weaponkills = {}
hook.Add("PlayerDeath","luctus_monitor_extra",function(victim,inflictor,attacker)
    if IsValid(attacker) and attacker:IsPlayer() and attacker:GetActiveWeapon() and IsValid(attacker:GetActiveWeapon()) then
        table.insert(weaponkills,{
            ["wepclass"] = attacker:GetActiveWeapon():GetClass(),
            ["attacker"] = attacker:SteamID(),
            ["victim"] = victim:SteamID()
        })
    end
end)

hook.Add("postLoadCustomDarkRPItems","luctus_monitor_extra",function()
    local alljobs = {}
    for k,v in pairs(RPExtraTeams) do
        table.insert(alljobs,v.command)
    end
    local data = {
        ["jobnames"] = alljobs,
        ["serverid"] = LUCTUS_MONITOR_SERVER_ID,
    }

    local ret = HTTP({
        failed = function(failMessage)
            print("[luctus_monitor] FAILED TO POST INITJOBS!")
            print("[luctus_monitor]",os.date("%H:%M:%S - %d/%m/%Y",os.time()))
            print(failMessage)
        end,
        success = function(httpcode,body,headers)
            debugPrint("[luctus_monitor] DoExtras Init successfull!")
        end, 
        method = "POST",
        url = LUCTUS_MONITOR_URL_INIT,
        body = util.TableToJSON(data),
        type = "application/json; charset=utf-8",
        timeout = 10
    })
end)

local jobtimes = {}
local jobswitches = {}

hook.Add("PlayerInitialSpawn","luctus_monitor_extra",function(ply)
    ply.switchedJob = CurTime()
end)

hook.Add("OnPlayerChangedTeam","luctus_monitor_extra",function(ply,before,after)
    local beforeName = team.GetName(before)
    local afterName = team.GetName(after)
    --switched
    if not jobswitches[afterName] then
        jobswitches[afterName] = 1
    else
        jobswitches[afterName] = jobswitches[afterName] + 1
    end
    --jobtime
    if not jobtimes[beforeName] then
        jobtimes[beforeName] = 1
    end
    jobtimes[beforeName] = jobtimes[beforeName] + math.Round(CurTime()-ply.switchedJob)
    ply.switchedJob = CurTime()
end)

hook.Add("PlayerDisconnect","luctus_monitor_extra",function(ply)
    if IsValid(ply) then
        local jobname = team.GetName(ply:Team())
        if not jobtimes[jobname] then
            jobtimes[jobname] = 0
        end
        jobtimes[jobname] = jobtimes[jobname] + math.Round(CurTime()-ply.switchedJob)
    end
end)


function LuctusMonitorDoExtras()
    
    for k,v in pairs(player.GetAll()) do
        local jobname = team.GetName(v:Team())
        if not jobtimes[jobname] then
            jobtimes[jobname] = 0
        end
        jobtimes[jobname] = jobtimes[jobname] + math.Round(CurTime()-v.switchedJob)
        v.switchedJob = CurTime()
    end
    
    local jsonJobTimes = {}
    for k,v in pairs(jobtimes) do
        table.insert(jsonJobTimes,{
            ["jobname"] = k,
            ["time"] = v,
        })
    end

    local jsonJobSwitches = {}
    for k,v in pairs(jobswitches) do
        table.insert(jsonJobSwitches,{
            ["jobname"] = k,
            ["amount"] = v,
        })
    end

    local data = {
        ["serverid"] = LUCTUS_MONITOR_SERVER_ID,
        ["weaponkills"] = weaponkills,
        ["jobtimes"] = jsonJobTimes,
        ["jobswitches"] = jsonJobSwitches,
    }
    
    local ret = HTTP({
        failed = function(failMessage)
            print("[luctus_monitor] FAILED TO POST STATS!")
            print("[luctus_monitor]",os.date("%H:%M:%S - %d/%m/%Y",os.time()))
            print(failMessage)
        end,
        success = function(httpcode,body,headers)
            debugPrint("[luctus_monitor] DoExtras Sync successfull!")
        end, 
        method = "POST",
        url = LUCTUS_MONITOR_URL_EXTRA,
        body = util.TableToJSON(data),
        type = "application/json; charset=utf-8",
        timeout = 10
    })
    
    --reset
    weaponkills = {}
    jobtimes = {}
    jobswitches = {}
    
    debugPrint("(DoExtras) Sent the following:")
    debugPrint("Table:")
    debugPrintTable(data)
    debugPrint("Json:")
    debugPrint(util.TableToJSON(data))
end

print("[luctus_monitor] sv loaded!")

