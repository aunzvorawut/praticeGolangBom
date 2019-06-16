package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "gitlab.com/wisdomvast/AiskaraokeServerGolang/routers"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	maxIdle, _ := beego.AppConfig.Int("maxIdle")
	maxConn, _ := beego.AppConfig.Int("maxConn")

	orm.RegisterDataBase("default", beego.AppConfig.String("mysqldriver"),
		beego.AppConfig.String("mysqluser")+":"+
			beego.AppConfig.String("mysqlpass")+"@tcp("+beego.AppConfig.String("mysqlurls")+":"+
			beego.AppConfig.String("mysqlport")+")/"+beego.AppConfig.String("mysqldb")+"?charset=utf8&loc=Asia%2FBangkok&parseTime=true",
		maxIdle, maxConn)

	orm.Debug, _ = beego.AppConfig.Bool("mysqldebug")

	name := "default" // Database alias. for many db
	force := false    // Drop table and re-create.
	verbose := true   // Print log.

	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		beego.Error(err)
	}

}

func main() {

	debug, err := beego.AppConfig.Bool("debug")
	if err != nil {
		debug = false
	}

	info, err := beego.AppConfig.Bool("info")
	if err != nil {
		info = false
	}

	if beego.BConfig.RunMode == "dev" {

		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.SetLevel(beego.LevelDebug)

	} else if beego.BConfig.RunMode == "test" {

		beego.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/app.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.SetLevel(beego.LevelDebug)

	} else if beego.BConfig.RunMode == "prd" {

		beego.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/app.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
		beego.SetLevel(beego.LevelError)
		if info {
			beego.SetLevel(beego.LevelInformational)
		}
		if debug {
			beego.SetLevel(beego.LevelDebug)
		}

	} else {

		beego.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/app.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
		beego.SetLevel(beego.LevelError)
		if info {
			beego.SetLevel(beego.LevelInformational)
		}
		if debug {
			beego.SetLevel(beego.LevelDebug)
		}

	}

	beego.BConfig.WebConfig.StaticDir["/static"] = "static"

	orm.RunCommand()
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Length", "Content-Type", "X-Atmosphere-tracking-id", "X-Atmosphere-Framework", "X-Cache-Dat", "X-Access-Token", "X-Nonce"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.Run()
}