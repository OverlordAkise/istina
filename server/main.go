package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io"
	//Config file
	"io/ioutil"
	"sigs.k8s.io/yaml"
	//Discord Webhooks
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	//Logging
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Mysql   string `json:"mysql"`
	Port    string `json:"port"`
	Debug   bool   `json:"debug"`
	Logfile string `json:"logfile"`
}

var db *sqlx.DB
var config Config
var LUCTUSDEBUG bool = false
var httpclient = http.Client{}
var discordURLRegex = regexp.MustCompile("^https:\\/\\/discord.com\\/api\\/webhooks\\/\\d+\\/[-_a-zA-Z0-9]+$")

func SetupRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/metrics"},
	}))
	RegisterMetrics(r, db)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.GET("/debugon", func(c *gin.Context) {
		LUCTUSDEBUG = true
		logger.Warn("Enabled debug output")
		c.String(200, "OK")
	})
	r.GET("/debugoff", func(c *gin.Context) {
		LUCTUSDEBUG = false
		logger.Warn("Disabled debug output")
		c.String(200, "OK")
	})
	r.POST("/tttstat", func(c *gin.Context) {
		var ts TTTStat
		err := c.BindJSON(&ts)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.String(400, "INVALID DATA")
			return
		}
		ts.Serverip = c.ClientIP()
		InsertTTTStat(ts, logger)
		c.String(200, "OK")
	})
	r.POST("/linuxstat", func(c *gin.Context) {
		var ls LuctusLinuxStat
		err := c.BindJSON(&ls)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.String(400, "INVALID DATA")
			return
		}
		ls.Realserverip = c.ClientIP()
		InsertLinuxStats(ls)
		c.String(200, "OK")
	})
	r.POST("/luaerror", func(c *gin.Context) {
		var le LuctusLuaError
		err := c.BindJSON(&le)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.String(400, "INVALID DATA")
			return
		}
		le.Serverip = c.ClientIP()
		InsertLuaError(le)
		c.String(200, "OK")
	})
	r.POST("/darkrpstat", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		var ds DarkRPStat
		err := c.BindJSON(&ds)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.String("body", string(body)),
			)
			c.String(400, "INVALID DATA")
			return
		}
		ds.Serverip = c.ClientIP()
		InsertDarkRPStat(ds, logger)
		c.String(200, "OK")
	})
	r.POST("/luctuslogs", func(c *gin.Context) {
		var ll LuctusLogs
		err := c.BindJSON(&ll)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.String(400, "INVALID DATA")
			return
		}
		ll.Serverip = c.ClientIP()
		InsertLuctusLogs(ll, logger)
		c.String(200, "OK")
	})
	r.POST("/playeravatar", func(c *gin.Context) {
		var pa PlayerAvatar
		err := c.BindJSON(&pa)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.String(400, "INVALID DATA")
			return
		}
		InsertPlayerAvatar(pa)
		c.String(200, "OK")
	})
	r.POST("/discordmsg", func(c *gin.Context) {
		var dc DiscordMessage
		err := c.BindJSON(&dc)
		if err != nil {
			logger.Error("Couldn't bind JSON",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.String(400, "INVALID DATA")
		}
		if !discordURLRegex.MatchString(dc.Url) {
			logger.Error("Discord Regex Mismatch",
				zap.String("url", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.String("wurl", dc.Url),
			)
			c.String(400, "INVALID URL")
			return
		}
		logger.Info("Sending discord webhook",
			zap.String("url", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.String("msg", dc.Msg),
			zap.String("tag", dc.Tag),
			zap.String("wurl", dc.Url),
		)
		NotifyDiscordWebhook(dc)
		c.String(200, "OK")
	})
	return r
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

	//logger
	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	cfg.OutputPaths = []string{
		config.Logfile,
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	fmt.Println("Debug mode:", LUCTUSDEBUG)
	gin.SetMode(gin.ReleaseMode)
	InitDatabase(config.Mysql)
	r := SetupRouter(logger)
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	fmt.Println("Now listening on *:" + config.Port)
	logger.Info("Starting server on *:" + config.Port)
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

	db.MustExec(`CREATE TABLE IF NOT EXISTS rpserver(
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

	db.MustExec(`CREATE TABLE IF NOT EXISTS rpplayer(
    id SERIAL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    serverid VARCHAR(50),
    steamid VARCHAR(50),
    nick VARCHAR(250),
    job VARCHAR(250),
    rank VARCHAR(50),
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
    addonsize INT,
    warns INT,
    money BIGINT
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS weaponkills(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    weaponclass VARCHAR(255),
    victim VARCHAR(50),
    attacker VARCHAR(50)
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS jobstats(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    jobname VARCHAR(255),
    switchedto BIGINT,
    timespent BIGINT,
    unique(serverid,jobname)
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS plyjobtimes(
    id SERIAL,
    ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    serverid VARCHAR(50),
    steamid VARCHAR(50),
    jobname VARCHAR(255),
    timespent BIGINT,
    unique(serverid,steamid,jobname)
    )`)

	db.MustExec(`CREATE TABLE IF NOT EXISTS bans(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    admin VARCHAR(255),
    target VARCHAR(50),
    reason TEXT,
    bantime BIGINT,
    curtime BIGINT
    )`)

	///Logs

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
    vcollidewireframe VARCHAR(5),
    alive BOOL
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
    attackerrole VARCHAR(20),
    hitgroup INT
    )`)

	///Joinstats

	db.MustExec(`CREATE TABLE IF NOT EXISTS joinstats(
    id SERIAL,
    ts TIMESTAMP,
    serverid VARCHAR(50),
    steamid VARCHAR(50),
    jointime BIGINT,
    connected BOOL
    )`)

	///Playeravatar

	db.MustExec(`CREATE TABLE IF NOT EXISTS playeravatar(
    id SERIAL,
    ts TIMESTAMP,
    steamid VARCHAR(50),
    steamid64 VARCHAR(50),
    image TEXT,
    unique(steamid,steamid64)
    )`)

	fmt.Println("DB initialized!")
}

func InsertLinuxStats(ls LuctusLinuxStat) {
	_, err := db.NamedExec("INSERT INTO linux(serverip,realserverip,cpuidle,cpusteal,cpuiowait,ramtotal,ramused,ramfree,diskfree,diskused,disktotal,diskpercentused) VALUES(:serverip,:realserverip,:cpuidle,:cpusteal,:cpuiowait,:ramtotal,:ramused,:ramfree,:diskfree,:diskused,:disktotal,:diskpercentused)", ls)
	if err != nil {
		panic(err)
	}
}

func InsertLuaError(le LuctusLuaError) {
	_, err := db.NamedExec("INSERT INTO luaerror(serverip,hash,error,stack,addon,gamemode,gameversion,os,ds,realm,version) VALUES(:serverip,:hash,:error,:stack,:addon,:gamemode,:gameversion,:os,:ds,:realm,:version)", le)
	if err != nil {
		panic(err)
	}
}

func InsertDarkRPStat(ds DarkRPStat, logger *zap.Logger) {
	tx := db.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			logger.Error("Error during InsertDarkRPStat SQL",
				zap.String("error", err.Error()),
			)
			err = tx.Rollback()
			if err != nil {
				panic(err)
			}
		}
	}()

	_, err := tx.NamedExec("INSERT INTO rpserver(serverid,serverip,map,gamemode,tickrateset,tickratecur,entscount,plycount,uptime,avgfps,avgping,luaramb,luarama) VALUES(:serverid,:serverip,:map,:gamemode,:tickrateset,:tickratecur,:entscount,:plycount,:uptime,:avgfps,:avgping,:luaramb,:luarama)", ds)
	if err != nil {
		panic(err)
	}
	if len(ds.Players) > 0 {
		_, err = tx.NamedExec("INSERT INTO rpplayer (serverid,steamid,nick,job,rank,fpsavg,fpslow,fpshigh,pingavg,pingcur,luaramb,luarama,packetslost,os,country,screensize,screenmode,jitver,ip,playtime,playtimel,online,hookthink,hooktick,hookhudpaint,hookhudpaintbackground,hookpredrawhud,hookcreatemove,concommands,funccount,addoncount,addonsize, warns, money) VALUES (:serverid, :steamid, :nick, :job, :rank, :fpsavg, :fpslow, :fpshigh, :pingavg, :pingcur, :luaramb, :luarama, :packetslost, :os, :country, :screensize, :screenmode, :jitver, :ip, :playtime, :playtimel, :online, :hookthink, :hooktick, :hookhudpaint, :hookhudpaintbackground, :hookpredrawhud, :hookcreatemove, :concommands, :funccount, :addoncount, :addonsize, :warns, :money)", ds.Players)
		if err != nil {
			panic(err)
		}
	}

	for _, v := range ds.Weaponkills {
		tx.MustExec("INSERT IGNORE INTO weaponkills(serverid,weaponclass,victim,attacker) VALUES(?,?,?,?)", ds.Serverid, v.Wepclass, v.Victim, v.Attacker)
	}

	for _, v := range ds.Jobs {
		tx.MustExec("INSERT INTO jobstats(serverid,jobname,switchedto,timespent) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE switchedto=switchedto+?, timespent=timespent+?;", ds.Serverid, v.Jobname, v.Switches, v.Playtime, v.Switches, v.Playtime)
	}

	for _, v := range ds.Plyjobs {
		tx.MustExec("INSERT INTO plyjobtimes(serverid,steamid,jobname,timespent) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE timespent=timespent+?;", ds.Serverid, v.Steamid, v.Jobname, v.Playtime, v.Playtime)
	}

	for _, v := range ds.Joinstats {
		tx.MustExec("INSERT INTO joinstats(serverid,steamid,jointime,connected) VALUES(?,?,?,?)", ds.Serverid, v.Steamid, v.Jointime, v.Connected)
	}

	for _, v := range ds.Bans {
		tx.MustExec("INSERT IGNORE INTO bans(serverid,admin,target,reason,bantime,curtime) VALUES(?,?,?,?,?,?)", ds.Serverid, v.Admin, v.Target, v.Reason, v.Bantime, v.Curtime)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func InsertLuctusLogs(ll LuctusLogs, logger *zap.Logger) {
	tx := db.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			logger.Error("Error during InsertLuctusLogs SQL",
				zap.String("error", err.Error()),
			)
			err = tx.Rollback()
			if err != nil {
				panic(err)
			}
		}
	}()
	for _, v := range ll.Logs {
		tx.MustExec("INSERT IGNORE INTO luctuslog(serverid,serverip,date,cat,msg) VALUES(?,?,?,?,?)", ll.Serverid, ll.Serverip, v.Date, v.Cat, v.Msg)
	}
	err := tx.Commit()
	if err != nil {
		panic(err)
	}
}

func InsertPlayerAvatar(pa PlayerAvatar) {
	db.MustExec("INSERT INTO playeravatar(steamid,steamid64,image) VALUES(?,?,?) ON DUPLICATE KEY UPDATE image=?;", pa.Steamid, pa.Steamid64, pa.Image, pa.Image)
}

func InsertTTTStat(ts TTTStat, logger *zap.Logger) {
	tx := db.MustBegin()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			logger.Error("Error during InsertTTTStat SQL",
				zap.String("error", err.Error()),
			)
			err = tx.Rollback()
			if err != nil {
				panic(err)
			}
		}
	}()
	_, err := tx.NamedExec("INSERT INTO tttserver(serverid,serverip,map,gamemode,roundstate,roundid,roundresult,tickrateset,tickratecur,entscount,plycount,avgfps,avgping,luaramb,luarama,innocent,traitor,detective,spectator,ainnocent,atraitor,adetective) VALUES(:serverid,:serverip,:map,:gamemode,:roundstate,:roundid,:roundresult,:tickrateset,:tickratecur,:entscount,:plycount,:avgfps,:avgping,:luaramb,:luarama,:innocent,:traitor,:detective,:spectator,:ainnocent,:atraitor,:adetective)", ts)
	if err != nil {
		panic(err)
	}

	if len(ts.Players) > 0 {
		_, err = tx.NamedExec("INSERT INTO tttplayer (serverid,steamid,nick,role,roundstate,roundid,fpsavg,fpslow,fpshigh,pingavg,pingcur,luaramb,luarama,packetslost,os,country,screensize,screenmode,jitver,ip,playtime,hookcount,hookthink,hooktick,hookhudpaint,hookhudpaintbackground,hookpredrawhud,hookcreatemove,concommands,funccount,addoncount,addonsize,svcheats,hosttimescale,svallowcslua,vcollidewireframe,alive) VALUES (:serverid,:steamid,:nick,:role,:roundstate,:roundid,:fpsavg,:fpslow,:fpshigh,:pingavg,:pingcur,:luaramb,:luarama,:packetslost,:os,:country,:screensize,:screenmode,:jitver,:ip,:playtime,:hookcount,:hookthink,:hooktick,:hookhudpaint,:hookhudpaintbackground,:hookpredrawhud,:hookcreatemove,:concommands,:funccount,:addoncount,:addonsize,:svcheats,:hosttimescale,:svallowcslua,:vcollidewireframe,:alive)", ts.Players)
		if err != nil {
			panic(err)
		}
	}

	for _, v := range ts.Joinstats {
		tx.MustExec("INSERT INTO joinstats(serverid,steamid,jointime,connected) VALUES(?,?,?,?)", ts.Serverid, v.Steamid, v.Jointime, v.Connected)
	}

	if len(ts.Kills) > 0 {
		_, err = tx.NamedExec("INSERT INTO tttkills (serverid,roundstate,roundid,wepclass,victim,attacker,victimrole,attackerrole,hitgroup) VALUES (:serverid,:roundstate,:roundid,:wepclass,:victim,:attacker,:victimrole,:attackerrole,:hitgroup)", ts.Kills)
		if err != nil {
			panic(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func NotifyDiscordWebhook(dc DiscordMessage) {
	data := map[string]interface{}{
		"content": dc.Tag + dc.Msg,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", dc.Url, bytes.NewReader(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = httpclient.Do(req)
	if err != nil {
		panic(err)
	}
}
