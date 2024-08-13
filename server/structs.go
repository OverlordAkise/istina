package main

type DarkRPStat struct {
	Id          float64         `json:"id" db:"id"`
	Ts          string          `json:"ts" db:"ts"`
	Serverid    string          `json:"serverid" db:"serverid" binding:"required"`
	Serverip    string          `json:"serverip" db:"serverip"`
	Map         string          `json:"map" db:"map"`
	Gamemode    string          `json:"gamemode" db:"gamemode"`
	Uptime      float64         `json:"uptime" db:"uptime"`
	Tickrateset float64         `json:"tickrateset" db:"tickrateset"`
	Tickratecur float64         `json:"tickratecur" db:"tickratecur"`
	Entscount   float64         `json:"entscount" db:"entscount"`
	Plycount    float64         `json:"plycount" db:"plycount"`
	Avgfps      float64         `json:"avgfps" db:"avgfps"`
	Avgping     float64         `json:"avgping" db:"avgping"`
	Luaramb     float64         `json:"luaramb" db:"luaramb"`
	Luarama     float64         `json:"luarama" db:"luarama"`
	Players     []DarkRPPlayer  `json:"players" db:"players"`
	Jobs        []DarkRPJobstat `json:"jobs" db:"jobs"`
	Plyjobs     []DarkRPPlyjob  `json:"plyjobs" db:"plyjobs"`
	Weaponkills []DarkRPKills   `json:"weaponkills" db:"weaponkills"`
	Joinstats   []Joinstatistic `json:"joinstats" db:"joinstats"`
	Bans        []UlxBan        `json:"bans" db:"bans"`
	Warns       []Warn          `json:"warns" db:"warns"`
}

type DarkRPPlayer struct {
	Id                     float64 `json:"id" db:"id"`
	Ts                     string  `json:"ts" db:"ts"`
	Serverid               string  `json:"serverid" db:"serverid" binding:"required"`
	Steamid                string  `json:"steamid" db:"steamid"`
	Nick                   string  `json:"nick" db:"nick"`
	Job                    string  `json:"job" db:"job"`
	Rank                   string  `json:"rank" db:"rank"`
	Fpsavg                 float64 `json:"fpsavg" db:"fpsavg"`
	Fpslow                 float64 `json:"fpslow" db:"fpslow"`
	Fpshigh                float64 `json:"fpshigh" db:"fpshigh"`
	Pingavg                float64 `json:"pingavg" db:"pingavg"`
	Pingcur                float64 `json:"pingcur" db:"pingcur"`
	Luaramb                float64 `json:"luaramb" db:"luaramb"`
	Luarama                float64 `json:"luarama" db:"luarama"`
	Packetslost            float64 `json:"packetslost" db:"packetslost"`
	Os                     string  `json:"os" db:"os"`
	Country                string  `json:"country" db:"country"`
	Screensize             string  `json:"screensize" db:"screensize"`
	Screenmode             string  `json:"screenmode" db:"screenmode"`
	Jitver                 string  `json:"jitver" db:"jitver"`
	Ip                     string  `json:"ip" db:"ip"`
	Playtime               float64 `json:"playtime" db:"playtime"`
	Playtimel              float64 `json:"playtimel" db:"playtimel"`
	Online                 bool    `json:"online" db:"online"`
	Hookthink              float64 `json:"hookthink" db:"hookthink"`
	Hooktick               float64 `json:"hooktick" db:"hooktick"`
	Hookhudpaint           float64 `json:"hookhudpaint" db:"hookhudpaint"`
	Hookhudpaintbackground float64 `json:"hookhudpaintbackground" db:"hookhudpaintbackground"`
	Hookpredrawhud         float64 `json:"hookpredrawhud" db:"hookpredrawhud"`
	Hookcreatemove         float64 `json:"hookcreatemove" db:"hookcreatemove"`
	Concommands            float64 `json:"concommands" db:"concommands"`
	Funccount              float64 `json:"funccount" db:"funccount"`
	Addoncount             float64 `json:"addoncount" db:"addoncount"`
	Addonsize              float64 `json:"addonsize" db:"addonsize"`
	Warns                  float64 `json:"warns" db:"warns"`
	Money                  float64 `json:"money" db:"money"`
}

type DarkRPKills struct {
	Wepclass string `json:"wepclass"`
	Victim   string `json:"victim"`
	Attacker string `json:"attacker"`
}
type DarkRPJobstat struct {
	Jobname  string  `json:"jobname" db:"jobname"`
	Switches float64 `json:"switches" db:"switches"`
	Playtime float64 `json:"playtime" db:"playtime"`
}
type DarkRPPlyjob struct {
	Steamid  string  `json:"steamid" db:"steamid"`
	Jobname  string  `json:"jobname" db:"jobname"`
	Playtime float64 `json:"playtime" db:"playtime"`
}
type Joinstatistic struct {
	Steamid   string  `json:"steamid" db:"steamid"`
	Jointime  float64 `json:"jointime" db:"jointime"`
	Connected bool    `json:"connected" db:"connected"`
}

type UlxBan struct {
	Admin   string  `json:"admin" db:"admin"`
	Target  string  `json:"target" db:"target"`
	Reason  string  `json:"reason" db:"reason"`
	Bantime float64 `json:"bantime" db:"bantime"`
	Curtime float64 `json:"curtime" db:"curtime"`
}

type Warn struct {
	Admin   string  `json:"admin" db:"admin"`
	Target  string  `json:"target" db:"target"`
	Reason  string  `json:"reason" db:"reason"`
}

type LuctusLuaError struct {
	Serverip    string `json:"serverip" form:"serverip" db:"serverip"`
	Hash        string `json:"hash" form:"hash" db:"hash"`
	Error       string `json:"error" form:"error" db:"error"`
	Stack       string `json:"stack" form:"stack" db:"stack"`
	Addon       string `json:"addon" form:"addon" db:"addon"`
	Gamemode    string `json:"gamemode" form:"gamemode" db:"gamemode"`
	Gameversion string `json:"gmv" form:"gmv" db:"gmv"`
	Os          string `json:"os" form:"os" db:"os"`
	Ds          string `json:"ds" form:"ds" db:"ds"`
	Realm       string `json:"realm" form:"realm" db:"realm"`
	Version     string `json:"v" form:"v" db:"v"`
}

///// Discord Webhook Message

type DiscordMessage struct {
	Url string `json:"url" form:"url"`
	Tag string `json:"tag" form:"tag"`
	Msg string `json:"msg" form:"msg"`
}

///// Logs

type LuctusLog struct {
	Msg  string `json:"msg" form:"msg"`
	Date string `json:"date" form:"date"`
	Cat  string `json:"cat" form:"cat"`
}

type LuctusLogs struct {
	Serverid string      `json:"serverid" db:"serverid"`
	Serverip string      `json:"serverip" db:"serverip"`
	Logs     []LuctusLog `json:"logs" form:"logs" db:"logs"`
}

///// Playeravatar

type PlayerAvatar struct {
	Steamid   string `json:"steamid" db:"steamid"`
	Steamid64 string `json:"steamid64" db:"steamid64"`
	Image     string `json:"image" db:"image"`
}

///// TTT

type TTTStat struct {
	Id          float64         `json:"id" db:"id"`
	Ts          string          `json:"ts" db:"ts"`
	Serverid    string          `json:"serverid" db:"serverid" binding:"required"`
	Serverip    string          `json:"serverip" db:"serverip"`
	Map         string          `json:"map" db:"map"`
	Gamemode    string          `json:"gamemode" db:"gamemode"`
	Roundstate  float64         `json:"roundstate" db:"roundstate"`
	Roundid     string          `json:"roundid" db:"roundid"`
	Roundresult float64         `json:"roundresult" db:"roundresult"`
	Tickrateset float64         `json:"tickrateset" db:"tickrateset"`
	Tickratecur float64         `json:"tickratecur" db:"tickratecur"`
	Entscount   float64         `json:"entscount" db:"entscount"`
	Plycount    float64         `json:"plycount" db:"plycount"`
	Avgfps      float64         `json:"avgfps" db:"avgfps"`
	Avgping     float64         `json:"avgping" db:"avgping"`
	Luaramb     float64         `json:"luaramb" db:"luaramb"`
	Luarama     float64         `json:"luarama" db:"luarama"`
	Innocent    float64         `json:"innocent" db:"innocent"`
	Traitor     float64         `json:"traitor" db:"traitor"`
	Detective   float64         `json:"detective" db:"detective"`
	Spectator   float64         `json:"spectator" db:"spectator"`
	Ainnocent   float64         `json:"ainnocent" db:"ainnocent"`
	Atraitor    float64         `json:"atraitor" db:"atraitor"`
	Adetective  float64         `json:"adetective" db:"adetective"`
	Players     []TTTPlayer     `json:"players" db:"players"`
	Kills       []TTTKills      `json:"kills" db:"kills"`
	Joinstats   []Joinstatistic `json:"joinstats" db:"joinstats"`
}

type TTTPlayer struct {
	Id                     float64 `json:"id" db:"id"`
	Ts                     string  `json:"ts" db:"ts"`
	Serverid               string  `json:"serverid" db:"serverid" binding:"required"`
	Steamid                string  `json:"steamid" db:"steamid"`
	Nick                   string  `json:"nick" db:"nick"`
	Role                   string  `json:"role" db:"role"`
	Roundid                string  `json:"roundid" db:"roundid"`
	Roundstate             float64 `json:"roundstate" db:"roundstate"`
	Fpsavg                 float64 `json:"fpsavg" db:"fpsavg"`
	Fpslow                 float64 `json:"fpslow" db:"fpslow"`
	Fpshigh                float64 `json:"fpshigh" db:"fpshigh"`
	Pingavg                float64 `json:"pingavg" db:"pingavg"`
	Pingcur                float64 `json:"pingcur" db:"pingcur"`
	Luaramb                float64 `json:"luaramb" db:"luaramb"`
	Luarama                float64 `json:"luarama" db:"luarama"`
	Packetslost            float64 `json:"packetslost" db:"packetslost"`
	Os                     string  `json:"os" db:"os"`
	Country                string  `json:"country" db:"country"`
	Screensize             string  `json:"screensize" db:"screensize"`
	Screenmode             string  `json:"screenmode" db:"screenmode"`
	Jitver                 string  `json:"jitver" db:"jitver"`
	Ip                     string  `json:"ip" db:"ip"`
	Playtime               float64 `json:"playtime" db:"playtime"`
	Hookcount              float64 `json:"hookcount" db:"hookcount"`
	Hookthink              float64 `json:"hookthink" db:"hookthink"`
	Hooktick               float64 `json:"hooktick" db:"hooktick"`
	Hookhudpaint           float64 `json:"hookhudpaint" db:"hookhudpaint"`
	Hookhudpaintbackground float64 `json:"hookhudpaintbackground" db:"hookhudpaintbackground"`
	Hookpredrawhud         float64 `json:"hookpredrawhud" db:"hookpredrawhud"`
	Hookcreatemove         float64 `json:"hookcreatemove" db:"hookcreatemove"`
	Concommands            float64 `json:"concommands" db:"concommands"`
	Funccount              float64 `json:"funccount" db:"funccount"`
	Addoncount             float64 `json:"addoncount" db:"addoncount"`
	Addonsize              float64 `json:"addonsize" db:"addonsize"`
	Svcheats               float64 `json:"svcheats" db:"svcheats"`
	Hosttimescale          float64 `json:"hosttimescale" db:"hosttimescale"`
	Svallowcslua           float64 `json:"svallowcslua" db:"svallowcslua"`
	Vcollidewireframe      float64 `json:"vcollidewireframe" db:"vcollidewireframe"`
	Alive                  bool    `json:"alive" db:"alive"`
}

type TTTKills struct {
	Serverid     string  `json:"serverid"  db:"serverid" binding:"required"`
	Roundid      string  `json:"roundid" db:"roundid"`
	Roundstate   float64 `json:"roundstate" db:"roundstate"`
	Wepclass     string  `json:"wepclass" db:"wepclass"`
	Victim       string  `json:"victim" db:"victim"`
	Attacker     string  `json:"attacker" db:"attacker"`
	Victimrole   string  `json:"victimrole" db:"victimrole"`
	Attackerrole string  `json:"attackerrole" db:"attackerrole"`
	Hitgroup     float64 `json:"hitgroup" db:"hitgroup"`
}
