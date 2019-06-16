package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gitlab.com/wisdomvast/AiskaraokeServerGolang/api/v6"
	"gitlab.com/wisdomvast/AiskaraokeServerGolang/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	nsV6 := beego.NewNamespace("/v6",
		beego.NSBefore(FilterDebug),
		beego.NSRouter("/debug", &v6.GlobalApi{}, "get,post:Test"),
		beego.NSNamespace("/content",
			beego.NSRouter("/appConfig", &v6.GlobalApi{}, "get,post:AppConfig"),
			beego.NSRouter("/newFeatureImage", &v6.GlobalApi{}, "get,post:NewFeatureImage"),
			beego.NSRouter("/popup/api", &v6.GlobalApi{}, "get,post:ApiPopup"),
			beego.NSRouter("/popup/view", &v6.GlobalApi{}, "get,post:WebviewPopup"),

			beego.NSRouter("/user/profile/all", &v6.GlobalApi{}, "get,post:AllProfile"),
			beego.NSRouter("/user/profile/personal", &v6.GlobalApi{}, "get,post:UserProfile"),
			beego.NSRouter("/user/profile/recentlySing", &v6.GlobalApi{}, "get,post:UserRecently"),
			beego.NSRouter("/user/profile/edit", &v6.GlobalApi{}, "get,post:EditProfile"),

			beego.NSRouter("/page/home", &v6.GlobalApi{}, "get,post:HomePage"),
			beego.NSRouter("/page/listCategory", &v6.GlobalApi{}, "get,post:ListCategoryPage"),
			beego.NSRouter("/page/categoryDetail", &v6.GlobalApi{}, "get,post:CategoryDetailPage"),
			beego.NSRouter("/page/publicPlaylistDetail", &v6.GlobalApi{}, "get,post:PublicPlaylistDetail"),
			beego.NSRouter("/page/songDetail", &v6.GlobalApi{}, "get,post:SongDetail"),
		),

		beego.NSNamespace("/service",
			beego.NSRouter("/keepDevice", &v6.GlobalApi{}, "get,post:KeepDevice"),
			beego.NSRouter("/subscribeMobile", &v6.GlobalApi{}, "get,post:SubscribeMobile"),
			beego.NSRouter("/song/session", &v6.GlobalApi{}, "get,post:UserPlaySession"),
			beego.NSRouter("/song/record/play", &v6.GlobalApi{}, "get,post:StatsPlaySong"),

			beego.NSRouter("/popup/view", &v6.GlobalApi{}, "get,post:WebviewPopup"),
		),
	)
	beego.AddNamespace(nsV6)

}

var FilterDebug = func(ctx *context.Context) {

	beego.Debug("body:", string(ctx.Input.RequestBody))
	beego.Debug("params:", ctx.Input.Params())
	beego.Debug("form:", ctx.Request.Form)
	beego.Debug("postform:", ctx.Request.PostForm)
	beego.Debug("RequestURI", ctx.Request.RequestURI)
	beego.Debug("method", ctx.Request.Method)
	for name, value := range ctx.Request.Header {
		beego.Debug(name, ":", value)
	}
	return
}
