package main

import (
    "net/http"
    "net/http/httptest"
	"fmt"
	"testing"
    "bytes"
    //db,config
    "github.com/gin-gonic/gin"
	"io/ioutil"
	"sigs.k8s.io/yaml"
)



var darkrpstattestbody *bytes.Buffer = bytes.NewBufferString( `{"avgping":2.5,"players":[{"pingcur":0.0,"addonsize":-1.0,"screenmode":"","playtimel":0.0,"concommands":-1.0,"fpsavg":-1.0,"fpshigh":-1.0,"funccount":-1.0,"country":"","ip":"Error!","playtime":0.0,"job":"","jitver":"","hookhudpaintbackground":-1.0,"addoncount":-1.0,"screensize":"","nick":"Bot03","pingavg":-1.0,"hookhudpaint":-1.0,"hooktick":-1.0,"hookthink":-1.0,"online":true,"serverid":"fe7af896-0f56-4d93-9b5f-bbc26a6792c2","os":"","fpslow":-1.0,"hookpredrawhud":-1.0,"packetslost":0.0,"hookcreatemove":-1.0,"luarama":-1.0,"luaramb":-1.0,"steamid":"BOT"},{"pingcur":5.0,"addonsize":4024.0,"screenmode":"window","playtimel":49.0,"concommands":42.0,"fpsavg":-2048.0,"fpshigh":256.0,"funccount":3208.0,"country":"AT","ip":"loopback","playtime":53.0,"job":"Medic","jitver":"LuaJIT 2.0.4","hookhudpaintbackground":-1.0,"addoncount":106.0,"screensize":"1920x1080","nick":"Medicman","pingavg":8.0,"hookhudpaint":12.0,"hooktick":3.0,"hookthink":6.0,"online":true,"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","os":"windows","fpslow":165.0,"hookpredrawhud":-1.0,"packetslost":0.0,"hookcreatemove":-1.0,"luarama":29571.5,"luaramb":54228.3,"steamid":"STEAM_0:0:12345678"}],"jobs":[{"jobname":"Medic","switches":1.0,"playtime":43.0},{"jobname":"Citizen","playtime":59.0}],"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","luaramb":16712.677734375,"tickrateset":66.66666815678282,"plycount":2.0,"uptime":63.314998626708987,"gamemode":"darkrp","luarama":9758.3076171875,"deaths":4.0,"avgfps":-1024.5,"map":"gm_construct","tickratecur":66.70377620082623,"weaponkills":[{"wepclass":"m9k_acr","victim":"BOT","attacker":"STEAM_0:0:12345678"},{"wepclass":"m9k_acr","victim":"BOT","attacker":"STEAM_0:0:12345678"},{"wepclass":"m9k_ragingbull","victim":"BOT","attacker":"STEAM_0:0:12345678"},{"wepclass":"m9k_ragingbull","victim":"BOT","attacker":"STEAM_0:0:12345678"}],"entscount":3985.0}`)


func TestMain(m *testing.M) {
    configfile, err := ioutil.ReadFile("./configtest.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configfile, &config)
	if err != nil {
		panic(err)
	}
    gin.SetMode(gin.ReleaseMode)
    //Setup database before test
	InitDatabase(config.Mysql)
    
    m.Run()
    
    //After test
    db.MustExec("DROP TABLE linux")
    db.MustExec("DROP TABLE luaerror")
    db.MustExec("DROP TABLE luastate")
    db.MustExec("DROP TABLE luaplayer")
    db.MustExec("DROP TABLE weaponkills")
    db.MustExec("DROP TABLE jobstats")
    db.MustExec("DROP TABLE luctuslog")
    db.MustExec("DROP TABLE tttserver")
    db.MustExec("DROP TABLE tttplayer")
    db.MustExec("DROP TABLE tttkills")
}

func TestWebroutes(t *testing.T) {
    r := SetupRouter()
    w := httptest.NewRecorder()
    req,_ := http.NewRequest("GET","/",nil)
    r.ServeHTTP(w,req)
    fmt.Println("Return:",w.Code,"Body:",w.Body.String())
    if w.Code != 200 {
        t.Fatal("/ had http code != 200")
    }
    if w.Body.String() != "OK" {
        t.Fatal("/ had body != 'OK'")
    }
}

func TestDarkRP(t *testing.T) {
	r := SetupRouter()
    w := httptest.NewRecorder()
    req,_ := http.NewRequest("POST","/darkrpstat",darkrpstattestbody)
    r.ServeHTTP(w,req)
    fmt.Println("Return:",w.Code,"Body:",w.Body.String())
    if w.Code != 200 {
        t.Fatal("/ had http code != 200")
    }
    if w.Body.String() != "OK" {
        t.Fatal("/ had body != 'OK'")
    }
    var darkrpstat DarkRPStat
    err := db.Get(&darkrpstat, "SELECT * FROM luastate LIMIT 1")
    if err != nil {
        t.Fatal(err)
    }
    //Main stat string
    if darkrpstat.Gamemode != "darkrp" {
        t.Fatal("Gamemode in DB is not the same as input!")
    }
    //Main stat float64
    if darkrpstat.Plycount != 2 {
        t.Fatal("Plycount in DB is not the same as input!")
    }
}
