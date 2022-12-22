--Luctus Monitor
--Made by OverlordAkise

local fps_all = 0
local fps_count = 0
local fps_lowest = 999
local fps_highest = 0
local ping_all = 0
local ping_count = 0

net.Receive("luctus_monitor_collect",function()
    SendStatistics()
end)

hook.Add("InitPostEntity", "luctus_monitor_collect_init", function()
    SendStatistics()
end)

function SendStatistics()
    if not LocalPlayer() or not IsValid(LocalPlayer()) then return end
    local fpsavg = math.Round(fps_all/fps_count)
    local pingavg = math.Round(ping_all/ping_count)
    local luaramb = collectgarbage("count")
    collectgarbage("collect")
    local luarama = collectgarbage("count")
    
    net.Start("luctus_monitor_collect")
        net.WriteInt(pingavg,32)
        net.WriteInt(fpsavg,32)
        net.WriteInt(fps_highest,32)
        net.WriteInt(fps_lowest,32)
        net.WriteFloat(luaramb)
        net.WriteFloat(luarama)
        net.WriteString(LuctusGetOS())
        net.WriteString(system.GetCountry())
        net.WriteString(ScrW().."x"..ScrH())
        net.WriteString(system.IsWindowed() and "window" or "fullscreen")
        net.WriteString(jit.version)
    net.SendToServer()
    
    --reset fps
    fps_all = 0
    fps_count = 0
    fps_lowest = 999
    fps_highest = 0
    --reset ping
    ping_all = 0
    ping_count = 0
end

timer.Create("luctus_monitor_timer",10,0,function()
    if not LocalPlayer() or not IsValid(LocalPlayer()) then return end
    --fps
    fps_count = fps_count + 1
    local fps = 1 / RealFrameTime()
    fps_all = fps_all + fps
    if fps > fps_highest then
        fps_highest = fps
    end
    if fps < fps_lowest then
        fps_lowest = fps
    end
    --ping
    
    if not LocalPlayer().Ping then return end
    ping_count = ping_count + 1
    local ping = LocalPlayer():Ping()
    ping_all = ping_all + ping
end)

function LuctusGetOS()
    if system.IsLinux() then
        return "linux"
    end
    if system.IsWindows() then
        return "windows"
    end
    if system.IsOSX() then
        return "osx"
    end
    return "unknown"
end

print("[luctus_monitor] cl loaded!")
