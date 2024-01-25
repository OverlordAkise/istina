--Luctus Monitor
--Made by OverlordAkise

--This script collects data about the server and players for analysis

LUCTUS_MONITOR_DEBUG = false

LUCTUS_MONITOR_URL = "http://localhost:7077/darkrpstat"
LUCTUS_MONITOR_URL_AVATAR = "http://localhost:7077/playeravatar"


function LuctusDebugPrint(text)
    if LUCTUS_MONITOR_DEBUG then
        print(text)
    end
end

function LuctusDebugPrintTable(tab)
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
util.AddNetworkString("luctus_monitor_connecttime")

timer.Create("luctus_monitor_autorestart",15,0,function()
    if not timer.Exists("luctus_monitor_timer") then
        print("[luctus_monitor] Starting Monitor timer")
        LuctusMonitorStart()
    end
end)

LUCTUS_MONITOR_PLAYERS = {}
local jobtimes = {}
local jobswitches = {}
local weaponkills = {}
local luctusJoinCache = {}
local luctusJoins = {}
local luctusBans = {}
local plyjobstats = {} --v2 jobstats

function LuctusMonitorStart()
    timer.Create("luctus_monitor_timer",180,0,function()
        --This takes time, so run it first, then send it
        GetCurrentTickrate()
        LuctusMonitorCollectPlayers()
        
        timer.Simple(5,function()
            LuctusMonitorDo()
        end)
    end)
end

--Monitor deaths
local lm_deaths = 0
hook.Add("PostPlayerDeath","luctus_monitor_stat",function(ply)
    lm_deaths = lm_deaths + 1
end)
hook.Add("OnNPCKilled", "luctus_monitor_stat",function(npc, attacker, inflictor)
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

function LuctusMonitorGetActiveWarns(ply)
    if LuctusWarnGetCount then
        return LuctusWarnGetCount(ply:SteamID())
    end
    if AWarn and AWarn.GetPlayerActiveWarnings then
        return AWarn:GetPlayerActiveWarnings(ply)
    end
    return -1
end

function LuctusMonitorDo()
    local data = {["players"] = {}}
    local server_avgfps = 0
    local server_avgfps_c = 0
    local server_avgping = 0
    local server_avgping_c = 0
    
    for k,v in pairs(LUCTUS_MONITOR_PLAYERS) do
        if v.fpsavg > 1 and v.fpsavg ~= 17 then
            server_avgfps = server_avgfps + v.fpsavg
            server_avgfps_c = server_avgfps_c + 1
        end
        server_avgping = server_avgping + v.pingcur
        server_avgping_c = server_avgping_c + 1
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
        data["avgping"] = 0
    else
        data["avgping"] = server_avgping/server_avgping_c
    end
    if server_avgfps_c == 0 or server_avgfps == 0 then
        data["avgfps"] = 0
    else
        data["avgfps"] = server_avgfps/server_avgfps_c
    end
    data["deaths"] = lm_deaths
    
    data["luaramb"] = collectgarbage("count")
    collectgarbage("collect")
    data["luarama"] = collectgarbage("count")
    
    
    --Jobtimes, weaponkills
    for k,v in pairs(player.GetAll()) do
        local jobname = team.GetName(v:Team())
        if not jobtimes[jobname] then
            jobtimes[jobname] = 0
        end
        jobtimes[jobname] = jobtimes[jobname] + math.Round(CurTime()-v.switchedJob)
        v.switchedJob = CurTime()
    end
    
    local jobStats = {}
    for k,v in pairs(jobtimes) do
        table.insert(jobStats,{
            ["jobname"] = k,
            ["playtime"] = v,
            ["switches"] = jobswitches[k],
        })
    end
    data["weaponkills"] = weaponkills
    data["jobs"] = jobStats
    data["plyjobs"] = LuctusMonitorGetJobStatsV2()
    data["joinstats"] = luctusJoins
    data["bans"] = luctusBans
    
    --Sending
    local ret = HTTP({
        failed = function(failMessage)
            print("[luctus_monitor] FAILED TO POST STATS!")
            print("[luctus_monitor]",os.date("%H:%M:%S - %d/%m/%Y",os.time()))
            print(failMessage)
        end,
        success = function(httpcode,body,headers)
            LuctusDebugPrint("[luctus_monitor] Do Sync successfull!")
            LuctusDebugPrint("[luctus_monitor] HTTP code:",httpcode)
            LuctusDebugPrint("[luctus_monitor] Body:",body)
        end, 
        method = "POST",
        url = LUCTUS_MONITOR_URL,
        body = util.TableToJSON(data),
        type = "application/json; charset=utf-8",
        timeout = 10
    })
    
    LuctusDebugPrint("(Do) Sent the following:")
    LuctusDebugPrint("Table:")
    LuctusDebugPrintTable(data)
    LuctusDebugPrint("Json:")
    LuctusDebugPrint(util.TableToJSON(data))
    --reset
    weaponkills = {}
    jobtimes = {}
    jobswitches = {}
    luctusJoins = {}
    luctusBans = {}
    plyjobstats = {}
    for k,ply in ipairs(plyjobstats) do
        plyjobstats[ply:SteamID()] = {}
    end
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
        ["rank"] = ply:GetUserGroup(),
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
        ["playtimel"] = 0,
        ["online"] = true,
        ["hookthink"] = -1,
        ["hooktick"] = -1,
        ["hookhudpaint"] = -1,
        ["hookhudpaintbackground"] = -1,
        ["hookpredrawhud"] = -1,
        ["hookcreatemove"] = -1,
        ["concommands"] = -1,
        ["funccount"] = -1,
        ["addoncount"] = -1,
        ["addonsize"] = -1,
        ["warns"] = -1,
        ["money"] = -1,
    }
    ply.lmonplaytime = CurTime()
    ply.lmonplaytimel = CurTime()
end)

function LuctusMonitorCollectPlayers()
    net.Start("luctus_monitor_collect")
    net.Broadcast()
end

net.Receive("luctus_monitor_collect",function(len,ply)
    LuctusDebugPrint("Got stats for new player:")
    LuctusDebugPrint(ply:Nick().."//"..ply:SteamID())
    LUCTUS_MONITOR_PLAYERS[ply:SteamID()] = {
        ["steamid"] = ply:SteamID(),
        ["nick"] = ply:Nick(),
        ["job"] = DarkRP and ply:getJobTable().name or "",
        ["rank"] = ply:GetUserGroup(),
        ["pingcur"] = ply:Ping(),
        ["pingavg"] = net.ReadInt(12),
        ["fpsavg"] = net.ReadInt(12),
        ["fpshigh"] = net.ReadInt(12),
        ["fpslow"] = net.ReadInt(12),
        ["packetslost"] = ply:PacketLoss(),
        ["luaramb"] = math.Clamp(net.ReadFloat(),-2,2047483647),
        ["luarama"] = math.Clamp(net.ReadFloat(),-2,2047483647),
        ["os"] = string.sub(net.ReadString(),1,10),
        ["country"] = string.sub(net.ReadString(),1,4),
        ["screensize"] = string.sub(net.ReadString(),1,15),
        ["screenmode"] = string.sub(net.ReadString(),1,15),
        ["jitver"] = string.sub(net.ReadString(),1,20),
        ["ip"] = ply:IPAddress(),
        ["serverid"] = LUCTUS_MONITOR_SERVER_ID,
        ["playtime"] = math.Round(CurTime() - ply.lmonplaytime),
        ["playtimel"] = math.Round(CurTime() - ply.lmonplaytimel),
        ["online"] = true,
        ["hookthink"] = net.ReadInt(10),
        ["hooktick"] = net.ReadInt(10),
        ["hookhudpaint"] = net.ReadInt(10),
        ["hookhudpaintbackground"] = net.ReadInt(10),
        ["hookpredrawhud"] = net.ReadInt(10),
        ["hookcreatemove"] = net.ReadInt(10),
        ["concommands"] = net.ReadInt(11),
        ["funccount"] = net.ReadInt(16),
        ["addoncount"] = net.ReadInt(16),
        ["addonsize"] = net.ReadInt(32),
        ["warns"] = LuctusMonitorGetActiveWarns(ply),
        ["money"] = math.Clamp(ply:getDarkRPVar("money") or -1,-2,9223372036854775807),
    }
    ply.lmonplaytimel = CurTime()
end)


--Weaponkills

hook.Add("PlayerDeath","luctus_monitor_extra",function(victim,inflictor,attacker)
    if IsValid(attacker) and attacker:IsPlayer() and attacker:GetActiveWeapon() and IsValid(attacker:GetActiveWeapon()) then
        table.insert(weaponkills,{
            ["wepclass"] = attacker:GetActiveWeapon():GetClass(),
            ["attacker"] = attacker:SteamID(),
            ["victim"] = victim:SteamID()
        })
    end
end)
hook.Add("OnNPCKilled", "luctus_monitor_extra",function(npc, attacker, inflictor)
    if IsValid(attacker) and attacker:IsPlayer() and attacker:GetActiveWeapon() and IsValid(attacker:GetActiveWeapon()) then
        table.insert(weaponkills,{
            ["wepclass"] = attacker:GetActiveWeapon():GetClass(),
            ["attacker"] = attacker:SteamID(),
            ["victim"] = "NPC"
        })
    end
end)


--Jobstats (+v2 for now)

hook.Add("PlayerInitialSpawn","luctus_monitor_extra",function(ply)
    ply.switchedJob = CurTime()
    ply.timeInCurJob = CurTime()
    plyjobstats[ply:SteamID()] = {}
end)

hook.Add("OnPlayerChangedTeam","luctus_monitor_extra",function(ply,before,after)
    local beforeName = team.GetName(before)
    local afterName = team.GetName(after)
    local steamid = ply:SteamID()
    --switches
    if not jobswitches[afterName] then
        jobswitches[afterName] = 1
    else
        jobswitches[afterName] = jobswitches[afterName] + 1
    end
    --jobtimes
    if not jobtimes[beforeName] then
        jobtimes[beforeName] = 1
    end
    jobtimes[beforeName] = jobtimes[beforeName] + math.Round(CurTime()-ply.switchedJob)
    ply.switchedJob = CurTime()
    --jobtimes v2
    if not plyjobstats[steamid] then plyjobstats[steamid] = {} end
    if not plyjobstats[steamid][beforeName] then
        plyjobstats[steamid][beforeName] = 1
    end
    plyjobstats[steamid][beforeName] = plyjobstats[steamid][beforeName] + math.Round(CurTime()-ply.timeInCurJob)
    ply.timeInCurJob = CurTime()
end)

hook.Add("PlayerDisconnect","luctus_monitor_extra",function(ply)
    if not IsValid(ply) then return end
    local jobname = team.GetName(ply:Team())
    local steamid = ply:SteamID()
    if not jobtimes[jobname] then
        jobtimes[jobname] = 0
    end
    jobtimes[jobname] = jobtimes[jobname] + math.Round(CurTime()-ply.switchedJob)
    --jobtimes v2
    if not plyjobstats[steamid][jobname] then
        plyjobstats[steamid][jobname] = 1
    end
    plyjobstats[steamid][jobname] = plyjobstats[steamid][jobname] + math.Round(CurTime()-ply.timeInCurJob)
end)

function LuctusMonitorGetJobStatsV2()
    local steamid = ""
    local jobname = ""
    for k,ply in ipairs(player.GetHumans()) do
        if not ply.timeInCurJob then continue end
        steamid = ply:SteamID()
        jobname = team.GetName(ply:Team())
	if not plyjobstats[steamid] then plyjobstats[steamid] = {} end
        if not plyjobstats[steamid][jobname] then
            plyjobstats[steamid][jobname] = 0
        end
        plyjobstats[steamid][jobname] = plyjobstats[steamid][jobname] + math.Round(CurTime()-ply.timeInCurJob)
        ply.timeInCurJob = CurTime()
    end
    local returnTab = {}
    for steamid,jobs in pairs(plyjobstats) do
        for job,playtime in pairs(jobs) do
            table.insert(returnTab,{
                steamid = steamid,
                jobname = job,
                playtime = playtime,
            })
        end
    end
    return returnTab
end


--Joinstats

gameevent.Listen("player_connect")
hook.Add("player_connect", "luctus_monitor_connecttime", function(data)
	luctusJoinCache[data.networkid] = CurTime()
end)
net.Receive("luctus_monitor_connecttime",function(len,ply)
    local sid = ply:SteamID()
    if luctusJoinCache[sid] then
        table.insert(luctusJoins,{
            ["steamid"] = sid,
            ["jointime"] = CurTime() - luctusJoinCache[sid],
            ["connected"] = true
        })
        luctusJoinCache[sid] = nil
    end
end)
hook.Add("PlayerDisconnected","luctus_monitor_connecttime",function(ply)
    local sid = ply:SteamID()
    if luctusJoinCache[sid] then
        table.insert(luctusJoins,{
            ["steamid"] = sid,
            ["jointime"] = CurTime() - luctusJoinCache[sid],
            ["connected"] = false
        })
        luctusJoinCache[sid] = nil
    end
end)


-- Bans

hook.Add("ULibPlayerBanned","luctus_monitor_bans",function(steamid,bandata)
    local callerName = "<server>"
    if bandata and bandata.admin then
        callerName = bandata.admin
    end
    local bantime = 0
    if bandata.unban != 0 then
        bantime = bandata.unban - bandata.time
    end
    table.insert(luctusBans,{
        ["admin"] = callerName,
        ["target"] = steamid,
        ["reason"] = bandata.reason,
        ["bantime"] = bantime,
        ["curtime"] = tonumber(bandata.time), --why is time a string?
    })
end)

hook.Add("SAM.BannedPlayer", "luctus_monitor_bans", function(ply, unban_date, reason, admin_steamid)
    local steamid = "UNKNOWN"
    if IsValid(ply) then
        steamid = ply:SteamID()
    end
    local nowtime = os.time()
    local unbantime = unban_date
    if unbantime > 0 then
        unbantime = unbantime-nowtime
    end
    table.insert(luctusBans,{
        ["admin"] = admin_steamid,
        ["target"] = steamid,
        ["reason"] = reason,
        ["bantime"] = unbantime,
        ["curtime"] = nowtime,
    })
end)
hook.Add("SAM.BannedSteamID", "luctus_monitor_bans", function(steamid, unban_date, reason, admin_steamid)
    local nowtime = os.time()
    local unbantime = unban_date
    if unbantime > 0 then
        unbantime = unbantime-nowtime
    end
    table.insert(luctusBans,{
        ["admin"] = admin_steamid,
        ["target"] = steamid,
        ["reason"] = reason,
        ["bantime"] = unbantime,
        ["curtime"] = nowtime,
    })
end)

hook.Add("Gextension_Ban","luctus_log_gexban",function(steamid64, length, reason, steamid64_admin, time)
    local victimid = util.SteamIDFrom64(steamid64)
    local adminid = util.SteamIDFrom64(steamid64_admin or "")
    
    table.insert(luctusBans,{
        ["admin"] = adminid,
        ["target"] = victimid,
        ["reason"] = reason,
        ["bantime"] = length,
        ["curtime"] = os.time(), --why is time a string?
    })
end)

--Avatar
util.AddNetworkString("luctus_istina_avatar")
hook.Add("OnPlayerChangedTeam","luctus_istina_avatar",function(ply,bt,at)
    if not ply.liaSynced then
        net.Start("luctus_istina_avatar")
        net.Send(ply)
        ply.liaSynced = true
        ply.liaDidAsk = true
    end
end)
net.Receive("luctus_istina_avatar",function(len,ply)
    if not ply.liaDidAsk then return end
    local pic = net.ReadString()
    local data = {}
    data["steamid"] = ply:SteamID()
    data["steamid64"] = ply:SteamID64()
    data["image"] = pic
    local ret = HTTP({
        failed = function(failMessage)
            print("[luctus_monitor] FAILED TO UPDATE PLAYER AVATAR!")
            print("[luctus_monitor]",os.date("%H:%M:%S - %d/%m/%Y",os.time()))
            print(failMessage)
        end,
        success = function(httpcode,body,headers)
            LuctusDebugPrint("[luctus_monitor] Playeravatar successfull!")
            LuctusDebugPrint("[luctus_monitor] HTTP code:",httpcode)
            LuctusDebugPrint("[luctus_monitor] Body:",body)
        end, 
        method = "POST",
        url = LUCTUS_MONITOR_URL_AVATAR,
        body = util.TableToJSON(data),
        type = "application/json; charset=utf-8",
        timeout = 10
    })
    --file.Write("server.jpg",util.Base64Decode(pic))
    ply.liaDidAsk = nil
end)


print("[luctus_monitor] sv loaded!")
