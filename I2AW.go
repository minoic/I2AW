package main

import (
	"github.com/MinoIC/I2AW/controllers"
	"github.com/astaxie/beego"
)

func main() {
	/*	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 31536000
		beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 31536000
		beego.BConfig.WebConfig.Session.SessionName = "5gEmGeJJFuYh7E"
		beego.BConfig.WebConfig.Session.SessionProvider = "file"
		beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
		beego.BConfig.WebConfig.Session.SessionOn = true*/
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/rgbvalue/:identifier", &controllers.RgbValueController{})
	beego.Run()
}
