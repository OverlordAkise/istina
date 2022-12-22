package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type LuctusLinuxStat struct {
	Serverip        string  `json:"serverip"`
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

type LuctusLuaStat struct {
	Serverid    string   `json:"serverid" db:"serverid" binding:"required"`
	Map         string   `json:"map" db:"map"`
	Gamemode    string   `json:"gamemode" db:"gamemode"`
	Tickrateset float64  `json:"tickrateset" db:"tickrateset"`
	Tickratecur float64  `json:"tickratecur" db:"tickratecur"`
	Entscount   float64  `json:"entscount" db:"entscount"`
	Plycount    float64  `json:"plycount" db:"plycount"`
	Avgfps      float64  `json:"avgfps" db:"avgfps"`
	Avgping     float64  `json:"avgping" db:"avgping"`
	Luaramb     float64  `json:"luaramb" db:"luaramb"`
	Luarama     float64  `json:"luarama" db:"luarama"`
	Players     []Player `json:"players" db:"players"`
}

type Player struct {
	Serverid   string  `json:"serverid"  db:"serverid" binding:"required"`
	Steamid    string  `json:"steamid"  db:"steamid"`
	Nick       string  `json:"nick"  db:"nick"`
	Job        string  `json:"job"  db:"job"`
	Fpsavg     float64 `json:"fpsavg" db:"fpsavg"`
	Fpslow     float64 `json:"fpslow" db:"fpslow"`
	Fpshigh    float64 `json:"fpshigh" db:"fpshigh"`
	Pingavg    float64 `json:"pingavg" db:"pingavg"`
	Pingcur    float64 `json:"pingcur" db:"pingcur"`
	Luaramb    float64 `json:"luaramb" db:"luaramb"`
	Luarama    float64 `json:"luarama" db:"luarama"`
	Os         string  `json:"os" form:"os"`
	Country    string  `json:"country" db:"country"`
	Screensize string  `json:"screensize" db:"screensize"`
	Screenmode string  `json:"screenmode" db:"screenmode"`
	Jitver     string  `json:"jitver" db:"jitver"`
	Ip         string  `json:"ip" db:"ip"`
	Playtime   float64 `json:"playtime" db:"playtime"`
	Online     bool    `json:"online" db:"online"`
}

type LuctusLuaStatExtra struct {
	Serverid    string        `json:"serverid" db:"db" binding:"required"`
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

type LuctusJobSyncs struct {
	Serverid string   `json:"serverid" binding:"required"`
	Jobnames []string `json:"jobnames"`
}

type LuctusLuaError struct {
	Hash        string `json:"hash" form:"hash"`
	Error       string `json:"error" form:"error"`
	Stack       string `json:"stack" form:"stack"`
	Addon       string `json:"addon" form:"addon"`
	Gamemode    string `json:"gamemode" form:"gamemode"`
	Gameversion string `json:"gmv" form:"gmv"`
	Os          string `json:"os" form:"os"`
	Ds          string `json:"ds" form:"ds"`
	Realm       string `json:"realm" form:"realm"`
	Version     string `json:"v" form:"v"`
}

var LUCTUSDEBUG bool = false

func debugPrint(a ...any) (n int, err error) {
	if LUCTUSDEBUG == true {
		return fmt.Println(a...)
	} else {
		return 0, nil
	}
}

func main() {
	fmt.Println("Starting up...")
	gin.SetMode(gin.ReleaseMode)
	InitDatabase("USER:PASSWORD@tcp(localhost:3306)/DATABASENAME")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.GET("/debugon", func(c *gin.Context) {
		LUCTUSDEBUG = true
		c.String(200, "OK")
	})
	r.GET("/debugoff", func(c *gin.Context) {
		LUCTUSDEBUG = false
		c.String(200, "OK")
	})
	r.POST("/linuxstat", func(c *gin.Context) {
		var ls LuctusLinuxStat
		c.BindJSON(&ls)
		InsertLinuxStats(c.ClientIP(), ls)
		c.String(200, "OK")
	})
	r.POST("/luaerror", func(c *gin.Context) {
		var ls LuctusLuaError
		c.Bind(&ls)
		InsertLuaError(c.ClientIP(), ls)
		c.String(200, "OK")
	})
	r.POST("/luastat", func(c *gin.Context) {
		var ls LuctusLuaStat
		c.BindJSON(&ls)
		InsertLuaStat(c.ClientIP(), ls)
		c.String(200, "OK")
	})
	r.POST("/luastatextra", func(c *gin.Context) {
		var ls LuctusLuaStatExtra
		c.BindJSON(&ls)
		InsertLuaStatExtra(c.ClientIP(), ls)
		c.String(200, "OK")
	})
	r.POST("/luajobinit", func(c *gin.Context) {
		var ls LuctusJobSyncs
		c.BindJSON(&ls)
		InsertLuaJobSyncs(c.ClientIP(), ls)
		c.String(200, "OK")
	})
	fmt.Println("Running...")
	r.Run("0.0.0.0:7077")
}

func InitDatabase(conString string) {
	var err error
	db, err = sqlx.Open("mysql", conString)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	db.Ping()

	db.MustExec(`CREATE TABLE IF NOT EXISTS linux(
    id SERIAL,
    ts TIMESTAMP,
    serverip VARCHAR(50),
    realserverip VARCHAR(50),
    cpuidle DOUBLE,
    cpusteal DOUBLE,
    cpuiowait DOUBLE,
    ramtotal INT,
    ramused INT,
    ramfree INT,
    diskfree INT,
    diskused INT,
    disktotal INT,
    diskpercentused INT
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS luaerror(
    id SERIAL,
    ts TIMESTAMP,
    serverip VARCHAR(50),
    hash VARCHAR(20),
    error TEXT,
    stack TEXT,
    addon TEXT,
    gamemode TEXT,
    gameversion TEXT,
    os TEXT,
    ds TEXT,
    realm VARCHAR(8),
    version VARCHAR(8)
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS luastate(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    serverip VARCHAR(50),
    map VARCHAR(50),
	gamemode VARCHAR(20),
	tickrateset INT,
	tickratecur INT,
	entscount INT,
	plycount INT,
    avgfps INT,
    avgping INT,
    luaramb INT,
    luarama INT
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS luaplayer(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    serverip VARCHAR(50),
    steamid VARCHAR(50),
	nick VARCHAR(50),
	job VARCHAR(50),
	fpsavg INT,
	fpslow INT,
	fpshigh INT,
	pingavg INT,
	pingcur INT,
	luaramb INT,
	luarama INT,
    os VARCHAR(10),
    country VARCHAR(4),
    screensize VARCHAR(15),
    screenmode VARCHAR(15),
    jitver VARCHAR(20),
    ip VARCHAR(25),
	playtime INT,
    online BOOL
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS weaponkills(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    serverip VARCHAR(50),
    weaponclass VARCHAR(255),
    victim VARCHAR(50),
    attacker VARCHAR(50)
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS jobstats(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    serverip VARCHAR(50),
    jobname VARCHAR(255),
    switchedto BIGINT,
    timespent BIGINT,
    unique(serverid,jobname)
    )`)
	fmt.Println("DB initialized...")
}

func InsertLinuxStats(serverip string, ls LuctusLinuxStat) {
	debugPrint(">>> InsertLinuxStats")
	debugPrint(ls)
	stmt, err := db.Prepare("INSERT INTO linux(serverip,realserverip,cpuidle,cpusteal,cpuiowait,ramtotal,ramused,ramfree,diskfree,diskused,disktotal,diskpercentused) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(ls.Serverip, serverip, ls.CpuIdle, ls.CpuSteal, ls.CpuIowait, ls.RamTotal, ls.RamUsed, ls.RamFree, ls.DiskFree, ls.DiskUsed, ls.DiskTotal, ls.DiskPercentUsed); err != nil {
		panic(err)
	}
	debugPrint("<<< InsertLinuxStats")
}

func InsertLuaError(serverip string, ls LuctusLuaError) {
	debugPrint(">>> InsertLuaError")
	debugPrint(ls)
	stmt, err := db.Prepare("INSERT INTO luaerror(serverip,hash,error,stack,addon,gamemode,gameversion,os,ds,realm,version) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(serverip, ls.Hash, ls.Error, ls.Stack, ls.Addon, ls.Gamemode, ls.Gameversion, ls.Os, ls.Ds, ls.Realm, ls.Version); err != nil {
		panic(err)
	}
	debugPrint("<<< InsertLuaError")
}

func InsertLuaStat(serverip string, ls LuctusLuaStat) {
	debugPrint("["+ls.Serverid+"]", ">>> InsertLuaStat")
	debugPrint("["+ls.Serverid+"]", ls)
	stmt, err := db.Prepare("INSERT INTO luastate(serverid,serverip,map,gamemode,tickrateset,tickratecur,entscount,plycount,avgfps,avgping,luaramb,luarama) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(ls.Serverid, serverip, ls.Map, ls.Gamemode, ls.Tickrateset, ls.Tickratecur, ls.Entscount, ls.Plycount, ls.Avgfps, ls.Avgping, ls.Luaramb, ls.Luarama); err != nil {
		panic(err)
	}
	debugPrint("["+ls.Serverid+"]", "Current players:", len(ls.Players))
	if len(ls.Players) > 0 {
		_, err = db.NamedExec("INSERT INTO luaplayer (serverid,serverip,steamid,nick,job,fpsavg,fpslow,fpshigh,pingavg,pingcur,luaramb,luarama,os,country,screensize,screenmode,jitver,ip,playtime,online) VALUES (:serverid, '"+serverip+"', :steamid, :nick, :job, :fpsavg, :fpslow, :fpshigh, :pingavg, :pingcur, :luaramb, :luarama, :os, :country, :screensize, :screenmode, :jitver, :ip, :playtime, :online)", ls.Players)
		if err != nil {
			panic(err)
		}
	}
	debugPrint("["+ls.Serverid+"]", "<<< InsertLuaStat")
}

func InsertLuaStatExtra(serverip string, ls LuctusLuaStatExtra) {
	debugPrint("["+ls.Serverid+"]", ">>> InsertLuaStatExtra")
	debugPrint("["+ls.Serverid+"]", "--- Weaponkills:")
	tx := db.MustBegin()
	for _, v := range ls.Weaponkills {
		debugPrint("["+ls.Serverid+"]", "Inserting:", ls.Serverid, serverip, v.Wepclass, v.Victim, v.Attacker)
		tx.MustExec("INSERT IGNORE INTO weaponkills(serverid,serverip,weaponclass,victim,attacker) VALUES(?,?,?,?,?)", ls.Serverid, serverip, v.Wepclass, v.Victim, v.Attacker)
	}
	tx.Commit()

	debugPrint("["+ls.Serverid+"]", "--- Jobplaytimes:")
	tx = db.MustBegin()
	for _, v := range ls.Jobtimes {
		debugPrint("["+ls.Serverid+"]", "Inserting:", v.Time, v.Jobname, ls.Serverid)
		tx.MustExec("UPDATE jobstats SET timespent = timespent + ? WHERE jobname = ? and serverid = ?", v.Time, v.Jobname, ls.Serverid)
	}
	tx.Commit()

	debugPrint("["+ls.Serverid+"]", "--- Jobswitches:")
	tx = db.MustBegin()
	for _, v := range ls.Jobswitches {
		debugPrint("["+ls.Serverid+"]", "Inserting:", v.Amount, v.Jobname, ls.Serverid)
		tx.MustExec("UPDATE jobstats SET switchedto = switchedto + ? WHERE jobname = ? and serverid = ?", v.Amount, v.Jobname, ls.Serverid)
	}
	tx.Commit()
	debugPrint("["+ls.Serverid+"]", "<<< InsertLuaStatExtra")
}

func InsertLuaJobSyncs(serverip string, ls LuctusJobSyncs) {
	debugPrint("["+ls.Serverid+"]", ">>> InsertLuaJobSyncs")
	debugPrint("["+ls.Serverid+"]", "--- Jobs:")
	tx := db.MustBegin()
	for _, v := range ls.Jobnames {
		debugPrint("["+ls.Serverid+"]", "Inserting:", ls.Serverid, serverip, v, 0, 0)
		tx.MustExec("INSERT IGNORE INTO jobstats(serverid,serverip,jobname,switchedto,timespent) VALUES(?,?,?,?,?)", ls.Serverid, serverip, v, 0, 0)
	}
	tx.Commit()
	debugPrint("["+ls.Serverid+"]", "<<< InsertLuaJobSyncs")
}
