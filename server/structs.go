package main

type LuctusLinuxStat struct {
	Serverip        string  `json:"serverip"`
	Realserverip    string  `json:"realserverip"`
	CpuIdle         int     `json:"cpuidle"`
	CpuSteal        float64 `json:"cpusteal"`
	CpuIowait       float64 `json:"cpuiowait"`
	RamTotal        int     `json:"ramtotal"`
	RamUsed         int     `json:"ramused"`
	RamFree         int     `json:"ramfree"`
	DiskTotal       int     `json:"disktotal"`
	DiskUsed        int     `json:"diskused"`
	DiskFree        int     `json:"diskfree"`
	DiskPercentUsed int     `json:"diskpercentused"`
}

type DarkRPStat struct {
	Serverid    string         `json:"serverid" db:"serverid" binding:"required"`
	Serverip    string         `json:"serverip" db:"serverip"`
	Map         string         `json:"map" db:"map"`
	Gamemode    string         `json:"gamemode" db:"gamemode"`
	Tickrateset float64        `json:"tickrateset" db:"tickrateset"`
	Tickratecur float64        `json:"tickratecur" db:"tickratecur"`
	Entscount   float64        `json:"entscount" db:"entscount"`
	Plycount    float64        `json:"plycount" db:"plycount"`
	Avgfps      float64        `json:"avgfps" db:"avgfps"`
	Avgping     float64        `json:"avgping" db:"avgping"`
	Luaramb     float64        `json:"luaramb" db:"luaramb"`
	Luarama     float64        `json:"luarama" db:"luarama"`
	Players     []DarkRPPlayer `json:"players" db:"players"`
}

type DarkRPPlayer struct {
	Serverid               string  `json:"serverid"  db:"serverid" binding:"required"`
	Steamid                string  `json:"steamid"  db:"steamid"`
	Nick                   string  `json:"nick"  db:"nick"`
	Job                    string  `json:"job"  db:"job"`
	Fpsavg                 float64 `json:"fpsavg" db:"fpsavg"`
	Fpslow                 float64 `json:"fpslow" db:"fpslow"`
	Fpshigh                float64 `json:"fpshigh" db:"fpshigh"`
	Pingavg                float64 `json:"pingavg" db:"pingavg"`
	Pingcur                float64 `json:"pingcur" db:"pingcur"`
	Luaramb                float64 `json:"luaramb" db:"luaramb"`
	Luarama                float64 `json:"luarama" db:"luarama"`
	Packetslost            float64 `json:"packetslost" db:"packetslost"`
	Os                     string  `json:"os" form:"os"`
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
}

type DarkRPExtraStat struct {
	Serverid    string        `json:"serverid" db:"serverid" binding:"required"`
	Serverip    string        `json:"serverip" db:"serverip"`
	Weaponkills []WeaponKills `json:"weaponkills" db:"weaponkills"`
	Jobtimes    []Jobtimes    `json:"jobtimes" db:"jobtimes"`
	Jobswitches []Jobswitches `json:"jobswitches" db:"jobswitches"`
}
type WeaponKills struct {
	Wepclass string `json:"wepclass"`
	Victim   string `json:"victim"`
	Attacker string `json:"attacker"`
}
type Jobtimes struct {
	Jobname string  `json:"jobname"`
	Time    float64 `json:"time"`
}
type Jobswitches struct {
	Jobname string  `json:"jobname"`
	Amount  float64 `json:"amount"`
}

type DarkRPJobSync struct {
	Serverid string   `json:"serverid" binding:"required"`
	Serverip string   `json:"serverip" db:"serverip"`
	Jobnames []string `json:"jobnames"`
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
