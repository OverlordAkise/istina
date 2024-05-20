package main

import (
	"bytes"
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"sigs.k8s.io/yaml"
	"go.uber.org/zap"
)

var darkrpbodystring = `{"avgping":2.5,"players":[{"pingcur":0.0,"addonsize":-1.0,"screenmode":"","playtimel":0.0,"concommands":-1.0,"fpsavg":-1.0,"fpshigh":-1.0,"funccount":-1.0,"country":"","ip":"Error!","playtime":0.0,"job":"","jitver":"","hookhudpaintbackground":-1.0,"addoncount":-1.0,"screensize":"","nick":"Bot03","pingavg":-1.0,"hookhudpaint":-1.0,"hooktick":-1.0,"hookthink":-1.0,"online":true,"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","os":"","fpslow":-1.0,"hookpredrawhud":-1.0,"packetslost":0.0,"hookcreatemove":-1.0,"luarama":-1.0,"luaramb":-1.0,"steamid":"BOT"},{"pingcur":5.0,"addonsize":4024.0,"screenmode":"window","playtimel":49.0,"concommands":42.0,"fpsavg":-2048.0,"fpshigh":256.0,"funccount":3208.0,"country":"AT","ip":"loopback","playtime":53.0,"job":"Medic","jitver":"LuaJIT 2.0.4","hookhudpaintbackground":-1.0,"addoncount":106.0,"screensize":"1920x1080","nick":"Medicman","pingavg":8.0,"hookhudpaint":12.0,"hooktick":3.0,"hookthink":6.0,"online":true,"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","os":"windows","fpslow":165.0,"hookpredrawhud":-1.0,"packetslost":0.0,"hookcreatemove":-1.0,"luarama":29571.5,"luaramb":54228.3,"steamid":"STEAM_0:0:12345678"}],"jobs":[{"jobname":"Medic","switches":1.0,"playtime":43.0},{"jobname":"Citizen","playtime":59.0}],"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","luaramb":16712.677734375,"tickrateset":66.66666815678282,"plycount":2.0,"uptime":63.314998626708987,"gamemode":"darkrp","luarama":9758.3076171875,"deaths":4.0,"avgfps":-1024.5,"map":"gm_construct","tickratecur":66.70377620082623,"weaponkills":[{"wepclass":"m9k_acr","victim":"BOT","attacker":"STEAM_0:0:12345678"},{"wepclass":"m9k_acr","victim":"BOT","attacker":"STEAM_0:0:12345678"},{"wepclass":"m9k_ragingbull","victim":"BOT","attacker":"STEAM_0:0:12345678"},{"wepclass":"m9k_ragingbull","victim":"BOT","attacker":"STEAM_0:0:12345678"}],"entscount":3985.0}`
var darkrpstattestbody *bytes.Buffer = bytes.NewBufferString(darkrpbodystring)
// var jobstatbody *bytes.Buffer = bytes.NewBufferString(darkrpbodystring)

var tttstattestbody *bytes.Buffer = bytes.NewBufferString(`{"ainnocent":1.0,"innocent":1.0,"spectator":1.0,"plycount":2.0,"players":[{"pingavg":5.0,"screenmode":"window","fpsavg":295.0,"luarama":22215.064453125,"role":"innocent","concommands":39.0,"addoncount":50.0,"hookpredrawhud":-1.0,"nick":"Medicman","funccount":3224.0,"steamid":"STEAM_0:0:12345678","pingcur":5.0,"addonsize":1721.0,"country":"AT","ip":"loopback","playtime":204.0,"roundid":"20230330135514","jitver":"LuaJIT 2.0.4","hookhudpaintbackground":-1.0,"hookcreatemove":-1.0,"screensize":"1600x900","sv_allowcslua":"1","hookcount":172.0,"host_timescale":"1.0","fpslow":243.0,"luaramb":38595.09765625,"roundstate":4.0,"vcollide_wireframe":"0","hookhudpaint":5.0,"packetslost":0.0,"os":"windows","hookthink":4.0,"sv_cheats":"0","serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","hooktick":2.0,"fpshigh":285.0}],"tickratecur":66.63189814221598,"atraitor":0.0,"roundid":"20230330135514","traitor":0.0,"entscount":441.0,"avgping":5.0,"detective":0.0,"luarama":8679.4794921875,"tickrateset":66.66666815678282,"roundstate":4.0,"uptime":207.02999877929688,"gamemode":"terrortown","luaramb":11542.5341796875,"kills":[{"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","wepclass":"m9k_mossberg590","victim":"BOT","roundid":"20230330135514","attackerrole":"innocent","victimrole":"traitor","roundstate":3.0,"attacker":"STEAM_0:0:12345678"}],"avgfps":295.0,"adetective":0.0,"roundresult":3.0,"serverid":"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx","map":"ttt_lego"}
`)

type StringValue struct {
	Value string `json:"value" db:"value"`
}
type NumberValue struct {
	Value int `json:"value" db:"value"`
}

var r *gin.Engine

func TestMain(m *testing.M) {
	configfile, err := os.ReadFile("./configtest.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configfile, &config)
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard

	//Setup database before test
	InitDatabase(config.Mysql)
	r = SetupRouter(zap.NewNop())

	fmt.Println("Running tests...")
	m.Run()

	//After test
	db.MustExec("DROP TABLE linux")
	db.MustExec("DROP TABLE luaerror")
	db.MustExec("DROP TABLE rpserver")
	db.MustExec("DROP TABLE rpplayer")
	db.MustExec("DROP TABLE weaponkills")
	db.MustExec("DROP TABLE jobstats")
	db.MustExec("DROP TABLE plyjobtimes")
	db.MustExec("DROP TABLE bans")
	db.MustExec("DROP TABLE luctuslog")
	db.MustExec("DROP TABLE tttserver")
	db.MustExec("DROP TABLE tttplayer")
	db.MustExec("DROP TABLE tttkills")
	db.MustExec("DROP TABLE joinstats")
	db.MustExec("DROP TABLE playeravatar")
}

func TestWebroutes(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	//fmt.Println("Return:",w.Code,"Body:",w.Body.String())
	if w.Code != 200 {
		t.Fatal("/ had http code != 200")
	}
	if w.Body.String() != "OK" {
		t.Fatal("/ had body != 'OK'")
	}
}

func TestDarkRP(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/darkrpstat", darkrpstattestbody)
	r.ServeHTTP(w, req)
	//fmt.Println("Return:",w.Code,"Body:",w.Body.String())
	if w.Code != 200 {
		t.Fatal("/ had http code != 200")
	}
	if w.Body.String() != "OK" {
		t.Fatal("/ had body != 'OK'")
	}
	//Checking darkrp server
	var darkrpstat DarkRPStat
	err := db.Get(&darkrpstat, "SELECT * FROM rpserver LIMIT 1")
	if err != nil {
		t.Fatal(err)
	}
	if darkrpstat.Gamemode != "darkrp" {
		t.Fatal("Gamemode in DB is not the same as input!")
	}
	if darkrpstat.Plycount != 2 {
		t.Fatal("Plycount in DB is not the same as input!")
	}
	//Checking darkrp player
	var rpplayer DarkRPPlayer
	err = db.Get(&rpplayer, "SELECT * FROM rpplayer WHERE steamid != 'BOT'")
	if err != nil {
		t.Fatal(err)
	}
	if rpplayer.Nick != "Medicman" {
		t.Fatal("Player Nick in DB is not the same as input!")
	}
	if rpplayer.Pingcur != 5 {
		t.Fatal("Player Ping in DB is not the same as input!")
	}
	//Checking table kills
	var sv StringValue
	err = db.Get(&sv, "SELECT victim as value FROM weaponkills WHERE attacker = 'STEAM_0:0:12345678' AND weaponclass = 'm9k_acr'")
	if err != nil {
		t.Fatal(err)
	}
	if sv.Value != "BOT" {
		t.Fatal("DarkRPKills victim in DB is not the same as input!")
	}
	//Checking table jobstats
	var nv NumberValue
	err = db.Get(&nv, "SELECT timespent as value FROM jobstats WHERE jobname = 'Citizen'")
	if err != nil {
		t.Fatal(err)
	}
	if nv.Value != 59 {
		t.Fatal("Jobstats timespent in DB is not the same as input!")
	}
}

func TestTTT(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tttstat", tttstattestbody)
	r.ServeHTTP(w, req)
	//fmt.Println("Return:",w.Code,"Body:",w.Body.String())
	if w.Code != 200 {
		t.Fatal("/ had http code != 200")
	}
	if w.Body.String() != "OK" {
		t.Fatal("/ had body != 'OK'")
	}
	//Checking tttserver
	var tttstat TTTStat
	err := db.Get(&tttstat, "SELECT * FROM tttserver LIMIT 1")
	if err != nil {
		t.Fatal(err)
	}
	if tttstat.Gamemode != "terrortown" {
		t.Fatal("Gamemode in DB is not the same as input!")
	}
	if tttstat.Plycount != 2 {
		t.Fatal("Plycount in DB is not the same as input!")
	}
	//Checking tttplayer
	var tttplayer TTTPlayer
	err = db.Get(&tttplayer, "SELECT * FROM tttplayer WHERE steamid != 'BOT'")
	if err != nil {
		t.Fatal(err)
	}
	if tttplayer.Nick != "Medicman" {
		t.Fatal("Player Nick in DB is not the same as input!")
	}
	if tttplayer.Hooktick != 2 {
		t.Fatal("Player Hooktick in DB is not the same as input!")
	}
	//Checking weaponkills
	var sv StringValue
	err = db.Get(&sv, "SELECT wepclass as value FROM tttkills WHERE attacker = 'STEAM_0:0:12345678'")
	if err != nil {
		t.Fatal(err)
	}
	if sv.Value != "m9k_mossberg590" {
		t.Fatal("TTTKills wepclass in DB is not the same as input!")
	}
}
