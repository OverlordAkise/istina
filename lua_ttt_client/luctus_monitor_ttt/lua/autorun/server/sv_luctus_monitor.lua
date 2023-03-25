--Luctus Monitor (TTT)
--Made by OverlordAkise

--This script collects data about the server and players for analysis

LUCTUS_MONITOR_DEBUG = false

LUCTUS_MONITOR_URL = "http://localhost:7077/tttstat"


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


LUCTUS_MONITOR_PLAYERS = {}
LUCTUS_MONITOR_ROUNDID = os.date("%Y%m%d%H%M%S",os.time())

function LuctusMonitorStart(delay,result,roundstate)
    LuctusMonitorGetRoles()
    timer.Simple(delay,function()
        LuctusMonitorGetPlayers()
        GetCurrentTickrate()
    end)
    timer.Simple(delay+1,function()
        LuctusMonitorSend(result,roundstate)
    end)
end

function LuctusMonitorGetPlayers()
    net.Start("luctus_monitor_collect")
    net.Broadcast()
end

--Hooks for sending stats
hook.Add("TTTPrepareRound","luctus_monitor",function()
    LuctusMonitorStart(5,-1,2) --longer due to sweps spawning
end)

hook.Add("TTTBeginRound","luctus_monitor",function()
    LuctusMonitorStart(3,-1,3)
end)

hook.Add("TTTEndRound","luctus_monitor",function(result)
    LuctusMonitorStart(3,result,4)
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

--Get roles immediately
LUCTUS_MONITOR_ROLES = {}
LUCTUS_MONITOR_ROLES_ALIVE = {}

function LuctusResetRoleCounter()
    LUCTUS_MONITOR_ROLES = {
        ["traitor"] = 0,
        ["detective"] = 0,
        ["innocent"] = 0,
        ["spectator"] = 0,
    }
    LUCTUS_MONITOR_ROLES_ALIVE = {
        ["traitor"] = 0,
        ["detective"] = 0,
        ["innocent"] = 0,
    }
end

LuctusResetRoleCounter()

function LuctusMonitorGetRoles()
    for k,v in pairs(player.GetAll()) do
        if v:IsSpec() then 
            LUCTUS_MONITOR_ROLES["spectator"] = LUCTUS_MONITOR_ROLES["spectator"] + 1
            continue
        end
        if v:Alive() then
            LUCTUS_MONITOR_ROLES_ALIVE[v:GetRoleString()] = LUCTUS_MONITOR_ROLES_ALIVE[v:GetRoleString()] + 1
        end
        LUCTUS_MONITOR_ROLES[v:GetRoleString()] = LUCTUS_MONITOR_ROLES[v:GetRoleString()] + 1
    end
end

function LuctusMonitorSend(roundResult,roundstate)
    if GetRoundState() == 2 then
        LUCTUS_MONITOR_ROUNDID = os.date("%Y%m%d%H%M%S",os.time())
    end
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
        v["serverid"] = LUCTUS_MONITOR_SERVER_ID
        v["roundstate"] = roundstate
        v["roundid"] = LUCTUS_MONITOR_ROUNDID
        table.insert(data["players"],v)
    end
    data["gamemode"] = engine.ActiveGamemode()
    data["map"] = game.GetMap()
    data["roundstate"] = roundstate
    data["tickrateset"] = 1/engine.TickInterval()
    data["tickratecur"] = LUCTUS_MONITOR_CURRENT_TICKRATE
    data["entscount"] = #ents.GetAll()
    data["plycount"] = #player.GetAll()
    data["uptime"] = CurTime()
    data["serverid"] = LUCTUS_MONITOR_SERVER_ID
    data["roundid"] = LUCTUS_MONITOR_ROUNDID
    data["roundresult"] = roundResult
    if server_avgping_c == 0 then
        data["avgfps"] = 0
        data["avgping"] = 0
    else
        data["avgfps"] = server_avgfps/server_avgfps_c
        data["avgping"] = server_avgping/server_avgping_c
    end
    data["innocent"] = LUCTUS_MONITOR_ROLES["innocent"]
    data["traitor"] = LUCTUS_MONITOR_ROLES["traitor"]
    data["detective"] = LUCTUS_MONITOR_ROLES["detective"]
    data["spectator"] = LUCTUS_MONITOR_ROLES["spectator"]
    data["ainnocent"] = LUCTUS_MONITOR_ROLES_ALIVE["innocent"]
    data["atraitor"] = LUCTUS_MONITOR_ROLES_ALIVE["traitor"]
    data["adetective"] = LUCTUS_MONITOR_ROLES_ALIVE["detective"]
    
    data["kills"] = LUCTUS_MONITOR_KILLS
    
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
            LuctusDebugPrint("[luctus_monitor] Do Sync successfull!")
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
    LUCTUS_MONITOR_PLAYERS = {}
    LUCTUS_MONITOR_KILLS = {}
    LuctusResetRoleCounter()
end

hook.Add("PlayerInitialSpawn","luctus_monitor_ply_init",function(ply)
    local jobname = ""
    if ply.getJobTable and ply:getJobTable() and ply:getJobTable().name then
        jobname = ply:getJobTable().name
    end
    LUCTUS_MONITOR_PLAYERS[ply:SteamID()] = {
        ["steamid"] = ply:SteamID(),
        ["nick"] = ply:Nick(),
        ["role"] = ply:GetRoleString(),
        ["roundstate"] = GetRoundState(),
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
        ["hookcount"] = -1,
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
        ["sv_cheats"] = "-1",
        ["host_timescale"] = "-1",
        ["sv_allowcslua"] = "-1",
        ["vcollide_wireframe"] = "-1",
    }
    ply.lmonplaytime = CurTime()
end)

net.Receive("luctus_monitor_collect",function(len,ply)
    LuctusDebugPrint("Got stats for new player:")
    LuctusDebugPrint(ply:Nick().."//"..ply:SteamID())
    LUCTUS_MONITOR_PLAYERS[ply:SteamID()] = {
        ["steamid"] = ply:SteamID(),
        ["nick"] = ply:Nick(),
        ["role"] = ply:GetRoleString(),
        ["roundstate"] = GetRoundState(),
        ["pingcur"] = ply:Ping(),
        ["pingavg"] = net.ReadInt(12),
        ["fpsavg"] = net.ReadInt(12),
        ["fpshigh"] = net.ReadInt(12),
        ["fpslow"] = net.ReadInt(12),
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
        ["hookcount"] = net.ReadInt(16),
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
        ["sv_cheats"] = net.ReadString(),
        ["host_timescale"] = net.ReadString(),
        ["sv_allowcslua"] = net.ReadString(),
        ["vcollide_wireframe"] = net.ReadString(),
    }
end)

--Kill counting
LUCTUS_MONITOR_KILLS = {}
hook.Add("PlayerDeath","luctus_monitor_kills",function(victim,inflictor,attacker)
    LuctusDebugPrint("Death occured, saving...")
    local wepstat = {
        ["serverid"] = LUCTUS_MONITOR_SERVER_ID,
        ["roundid"] = LUCTUS_MONITOR_ROUNDID,
        ["roundstate"] = GetRoundState(),
        ["victim"] = victim:SteamID(),
        ["victimrole"] = victim:GetRoleString(),
        ["wepclass"] = "",
        ["attacker"] = "",
        ["attackerrole"] = "",
    }
    if IsValid(attacker) and attacker:IsPlayer() and not victim:IsSpec() then
        LuctusDebugPrint("Logging killer of death")
        wepstat["attacker"] = attacker:SteamID()
        wepstat["attackerrole"] = attacker:GetRoleString()
        if attacker:GetActiveWeapon() and IsValid(attacker:GetActiveWeapon()) then
            wepstat["wepclass"] = attacker:GetActiveWeapon():GetClass()
        end
    end
    table.insert(LUCTUS_MONITOR_KILLS,wepstat)
end)

print("[luctus_monitor] sv loaded!")
