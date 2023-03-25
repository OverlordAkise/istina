package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	//Config file
	"io/ioutil"
	"sigs.k8s.io/yaml"
)

type Config struct {
	Mysql string `json:"mysql"`
	Port  string `json:"port"`
	Debug bool   `json:"debug"`
}

var db *sqlx.DB
var config Config
var LUCTUSDEBUG bool = false

func debugPrint(a ...any) (n int, err error) {
	if LUCTUSDEBUG == true {
		return fmt.Println(a...)
	} else {
		return 0, nil
	}
}

func main() {
	fmt.Println("Starting!")
	configfile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configfile, &config)
	if err != nil {
		panic(err)
	}
	LUCTUSDEBUG = config.Debug
	fmt.Println("Debug mode:", LUCTUSDEBUG)
	gin.SetMode(gin.ReleaseMode)
	InitDatabase(config.Mysql)
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
	r.POST("/tttstat", func(c *gin.Context) {
		var data TTTStat
		c.BindJSON(&data)
		data.Serverip = c.ClientIP()
		InsertTTTStat(data)
		c.String(200, "OK")
	})
	r.POST("/linuxstat", func(c *gin.Context) {
		var ls LuctusLinuxStat
		c.BindJSON(&ls)
		ls.Realserverip = c.ClientIP()
		InsertLinuxStats(ls)
		c.String(200, "OK")
	})
	r.POST("/luaerror", func(c *gin.Context) {
		var ls LuctusLuaError
		c.Bind(&ls)
		ls.Serverip = c.ClientIP()
		InsertLuaError(ls)
		c.String(200, "OK")
	})
	r.POST("/darkrpstat", func(c *gin.Context) {
		var ls DarkRPStat
		c.BindJSON(&ls)
		ls.Serverip = c.ClientIP()
		InsertDarkRPStat(ls)
		c.String(200, "OK")
	})
	r.POST("/luctuslogs", func(c *gin.Context) {
		var ll LuctusLogs
		c.BindJSON(&ll)
		ll.Serverip = c.ClientIP()
		InsertLuctusLogs(ll)
		c.String(200, "OK")
	})
	fmt.Println("Now listening on *:" + config.Port)
	r.Run("0.0.0.0:" + config.Port)
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
    uptime INT,
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
    packetslost INT,
    os VARCHAR(10),
    country VARCHAR(4),
    screensize VARCHAR(15),
    screenmode VARCHAR(15),
    jitver VARCHAR(20),
    ip VARCHAR(25),
    playtime INT,
    playtimel INT,
    online BOOL,
    hookthink INT,
    hooktick INT,
    hookhudpaint INT,
    hookhudpaintbackground INT,
    hookpredrawhud INT,
    hookcreatemove INT,
    concommands INT,
    funccount INT,
    addoncount INT,
    addonsize INT
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

	db.MustExec(`CREATE TABLE IF NOT EXISTS luctuslog(
    id SERIAL,
    ts TIMESTAMP,
    date VARCHAR(24),
    serverid VARCHAR(50),
    serverip VARCHAR(50),
    cat VARCHAR(255),
    msg TEXT
    )`)

	///TTT

	db.MustExec(`CREATE TABLE IF NOT EXISTS tttserver(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    serverip VARCHAR(50),
    map VARCHAR(50),
    gamemode VARCHAR(20),
    roundstate INT,
    roundid VARCHAR(20),
    roundresult INT,
    tickrateset INT,
    tickratecur INT,
    entscount INT,
    plycount INT,
    avgfps INT,
    avgping INT,
    luaramb INT,
    luarama INT,
    innocent INT,
    traitor INT,
    detective INT,
    spectator INT,
    ainnocent INT,
    atraitor INT,
    adetective INT
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS tttplayer(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    steamid VARCHAR(50),
    nick VARCHAR(50),
    role VARCHAR(20),
    roundstate INT,
    roundid VARCHAR(20),
    fpsavg INT,
    fpslow INT,
    fpshigh INT,
    pingavg INT,
    pingcur INT,
    luaramb INT,
    luarama INT,
    packetslost INT,
    os VARCHAR(10),
    country VARCHAR(4),
    screensize VARCHAR(15),
    screenmode VARCHAR(15),
    jitver VARCHAR(20),
    ip VARCHAR(25),
    playtime INT,
    hookcount INT,
    hookthink INT,
    hooktick INT,
    hookhudpaint INT,
    hookhudpaintbackground INT,
    hookpredrawhud INT,
    hookcreatemove INT,
    concommands INT,
    funccount INT,
    addoncount INT,
    addonsize INT,
    svcheats VARCHAR(5),
    hosttimescale VARCHAR(5),
    svallowcslua VARCHAR(5),
    vcollidewireframe VARCHAR(5)
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS tttkills(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    roundstate INT,
    roundid VARCHAR(20),
    wepclass VARCHAR(255),
    victim VARCHAR(50),
    attacker VARCHAR(50),
    victimrole VARCHAR(20),
    attackerrole VARCHAR(20)
    )`)

	fmt.Println("DB initialized!")
}

func InsertLinuxStats(ls LuctusLinuxStat) {
	debugPrint(">>> InsertLinuxStats")
	debugPrint(ls)
	_, err := db.NamedExec("INSERT INTO linux(serverip,realserverip,cpuidle,cpusteal,cpuiowait,ramtotal,ramused,ramfree,diskfree,diskused,disktotal,diskpercentused) VALUES(:serverip,:realserverip,:cpuidle,:cpusteal,:cpuiowait,:ramtotal,:ramused,:ramfree,:diskfree,:diskused,:disktotal,:diskpercentused)", ls)
	if err != nil {
		panic(err)
	}
	debugPrint("<<< InsertLinuxStats")
}

func InsertLuaError(ls LuctusLuaError) {
	debugPrint(">>> InsertLuaError")
	debugPrint(ls)
	_, err := db.NamedExec("INSERT INTO luaerror(serverip,hash,error,stack,addon,gamemode,gameversion,os,ds,realm,version) VALUES(:serverip,:hash,:error,:stack,:addon,:gamemode,:gameversion,:os,:ds,:realm,:version)", ls)
	if err != nil {
		panic(err)
	}
	debugPrint("<<< InsertLuaError")
}

func InsertDarkRPStat(ls DarkRPStat) {
	debugPrint("["+ls.Serverid+"]", ">>> InsertLuaStat")
	debugPrint("["+ls.Serverid+"]", ls)
	tx := db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO luastate(serverid,serverip,map,gamemode,tickrateset,tickratecur,entscount,plycount,uptime,avgfps,avgping,luaramb,luarama) VALUES(:serverid,:serverip,:map,:gamemode,:tickrateset,:tickratecur,:entscount,:plycount,:uptime,:avgfps,:avgping,:luaramb,:luarama)", ls)
	if err != nil {
		panic(err)
	}
	debugPrint("["+ls.Serverid+"]", "--- DarkRP players", len(ls.Players))
	if len(ls.Players) > 0 {
		_, err = tx.NamedExec("INSERT INTO luaplayer (serverid,serverip,steamid,nick,job,fpsavg,fpslow,fpshigh,pingavg,pingcur,luaramb,luarama,packetslost,os,country,screensize,screenmode,jitver,ip,playtime,playtimel,online,hookthink,hooktick,hookhudpaint,hookhudpaintbackground,hookpredrawhud,hookcreatemove,concommands,funccount,addoncount,addonsize) VALUES (:serverid, '"+ls.Serverip+"', :steamid, :nick, :job, :fpsavg, :fpslow, :fpshigh, :pingavg, :pingcur, :luaramb, :luarama, :packetslost, :os, :country, :screensize, :screenmode, :jitver, :ip, :playtime, :playtimel, :online, :hookthink, :hooktick, :hookhudpaint, :hookhudpaintbackground, :hookpredrawhud, :hookcreatemove, :concommands, :funccount, :addoncount, :addonsize)", ls.Players)
		if err != nil {
			panic(err)
		}
	}

	debugPrint("["+ls.Serverid+"]", "--- Weaponkills:")
	for _, v := range ls.Weaponkills {
		debugPrint("["+ls.Serverid+"]", "Inserting:", ls.Serverid, ls.Serverip, v.Wepclass, v.Victim, v.Attacker)
		tx.MustExec("INSERT IGNORE INTO weaponkills(serverid,serverip,weaponclass,victim,attacker) VALUES(?,?,?,?,?)", ls.Serverid, ls.Serverip, v.Wepclass, v.Victim, v.Attacker)
	}

	debugPrint("["+ls.Serverid+"]", "--- Jobstats:")
	for _, v := range ls.Jobs {
		debugPrint("["+ls.Serverid+"]", "Inserting:", v)
		tx.MustExec("INSERT INTO jobstats(serverid,serverip,jobname,switchedto,timespent) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE switchedto=switchedto+?, timespent=timespent+?;", ls.Serverid, ls.Serverip, v.Jobname, v.Switches, v.Playtime, v.Switches, v.Playtime)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	debugPrint("["+ls.Serverid+"]", "<<< InsertLuaStat")
}

func InsertLuctusLogs(ll LuctusLogs) {
	debugPrint("["+ll.Serverid+"]", ">>> InsertLuctusLogs")
	debugPrint("["+ll.Serverid+"]", "--- LogLines:", len(ll.Logs))
	tx := db.MustBegin()
	for _, v := range ll.Logs {
		debugPrint("["+ll.Serverid+"]", "Inserting:", ll.Serverid, ll.Serverip, v.Date, v.Cat, v.Msg)
		tx.MustExec("INSERT IGNORE INTO luctuslog(serverid,serverip,date,cat,msg) VALUES(?,?,?,?,?)", ll.Serverid, ll.Serverip, v.Date, v.Cat, v.Msg)
	}
	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	debugPrint("["+ll.Serverid+"]", "<<< InsertLuctusLogs")
}

func InsertTTTStat(data TTTStat) {
	debugPrint("["+data.Serverid+"]", ">>> InsertTTTStat")
	debugPrint("["+data.Serverid+"]", "All data:", data)
	tx := db.MustBegin()
	_, err := tx.NamedExec("INSERT INTO tttserver(serverid,serverip,map,gamemode,roundstate,roundid,roundresult,tickrateset,tickratecur,entscount,plycount,avgfps,avgping,luaramb,luarama,innocent,traitor,detective,spectator,ainnocent,atraitor,adetective) VALUES(:serverid,:serverip,:map,:gamemode,:roundstate,:roundid,:roundresult,:tickrateset,:tickratecur,:entscount,:plycount,:avgfps,:avgping,:luaramb,:luarama,:innocent,:traitor,:detective,:spectator,:ainnocent,:atraitor,:adetective)", data)
	if err != nil {
		panic(err)
	}

	debugPrint("["+data.Serverid+"]", "Current players:", len(data.Players))
	if len(data.Players) > 0 {
		_, err = tx.NamedExec("INSERT INTO tttplayer (serverid,steamid,nick,role,roundstate,roundid,fpsavg,fpslow,fpshigh,pingavg,pingcur,luaramb,luarama,packetslost,os,country,screensize,screenmode,jitver,ip,playtime,hookcount,hookthink,hooktick,hookhudpaint,hookhudpaintbackground,hookpredrawhud,hookcreatemove,concommands,funccount,addoncount,addonsize,svcheats,hosttimescale,svallowcslua,vcollidewireframe) VALUES (:serverid,:steamid,:nick,:role,:roundstate,:roundid,:fpsavg,:fpslow,:fpshigh,:pingavg,:pingcur,:luaramb,:luarama,:packetslost,:os,:country,:screensize,:screenmode,:jitver,:ip,:playtime,:hookcount,:hookthink,:hooktick,:hookhudpaint,:hookhudpaintbackground,:hookpredrawhud,:hookcreatemove,:concommands,:funccount,:addoncount,:addonsize,:svcheats,:hosttimescale,:svallowcslua,:vcollidewireframe)", data.Players)
		if err != nil {
			panic(err)
		}
	}

	debugPrint("["+data.Serverid+"]", "Kills:", len(data.Kills))
	if len(data.Kills) > 0 {
		_, err = tx.NamedExec("INSERT INTO tttkills (serverid,roundstate,roundid,wepclass,victim,attacker,victimrole,attackerrole) VALUES (:serverid,:roundstate,:roundid,:wepclass,:victim,:attacker,:victimrole,:attackerrole)", data.Kills)
		if err != nil {
			panic(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	debugPrint("["+data.Serverid+"]", "<<< InsertTTTStat")
}
