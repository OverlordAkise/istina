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
        net.WriteInt(pingavg,12)
        net.WriteInt(fpsavg,12)
        net.WriteInt(fps_highest,12)
        net.WriteInt(fps_lowest,12)
        net.WriteFloat(luaramb)
        net.WriteFloat(luarama)
        net.WriteString(LuctusGetOS())
        net.WriteString(system.GetCountry())
        net.WriteString(ScrW().."x"..ScrH())
        net.WriteString(system.IsWindowed() and "window" or "fullscreen")
        net.WriteString(jit.version)
        local ht = hook.GetTable()
        if ht then
            net.WriteInt(ht["Think"] and table.Count(ht["Think"]) or -1,10)
            net.WriteInt(ht["Tick"] and table.Count(ht["Tick"]) or -1,10)
            net.WriteInt(ht["HUDPaint"] and table.Count(ht["HUDPaint"]) or -1,10)
            net.WriteInt(ht["HUDPaintBackground"] and table.Count(ht["HUDPaintBackground"]) or -1,10)
            net.WriteInt(ht["PreDrawHUD"] and table.Count(ht["PreDrawHUD"]) or -1,10)
            net.WriteInt(ht["CreateMove"] and table.Count(ht["CreateMove"]) or -1,10)
        else
            net.WriteInt(-1,10)
            net.WriteInt(-1,10)
            net.WriteInt(-1,10)
            net.WriteInt(-1,10)
            net.WriteInt(-1,10)
            net.WriteInt(-1,10)
        end
        local cts = concommand and concommand.GetTable() and table.Count(concommand.GetTable()) or -1
        net.WriteInt(cts,11)
        net.WriteInt(table.Count(_G),16)
        
        local together = 0
        local amount = 0
        for k,v in pairs(engine.GetAddons()) do
            if v["mounted"] then
                together = together + v["size"]
                amount = amount + 1
            end
        end
        net.WriteInt(amount,16)
        net.WriteInt(together/1024/1024,32)
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
    if not system.HasFocus() then return end --tabbed out
    --fps
    local fps = 1 / RealFrameTime()
    if fps < 1 then return end --too early
    fps_count = fps_count + 1
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

--Joinstats
hook.Add("InitPostEntity", "luctus_monitor_connecttime", function()
	net.Start("luctus_monitor_connecttime")
	net.SendToServer()
end)

--Avatar
net.Receive("luctus_istina_avatar",function()
    if IsValid(LuctusIstinaAvatar) then LuctusIstinaAvatar:Remove() end
    LuctusIstinaAvatar = vgui.Create("AvatarImage")
    LuctusIstinaAvatar:SetSize(256,256)
    LuctusIstinaAvatar:SetPos(4,ScrH()/2)
    LuctusIstinaAvatar:SetPlayer(LocalPlayer(),256)
    LuctusIstinaAvatar:SetAlpha(255)
    LuctusIstinaAvatar:SetDrawOnTop(true)
    hook.Add("PostRender","luctus_zima",function()
        local pic = util.Base64Encode(render.Capture({
                format = "jpg",
                x = 4,
                y = ScrH()/2,
                w = 256,
                h = 256,
        }))
        LuctusIstinaAvatar:Remove()
        net.Start("luctus_istina_avatar")
            net.WriteString(pic)
        net.SendToServer()
        hook.Remove("PostRender","luctus_zima")
    end)
end)

print("[luctus_monitor] cl loaded!")
