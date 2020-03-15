package main

import (
	"github.com/MinoIC/I2AW/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/value/:identifier", &controllers.ValueController{})
	beego.Run()
}
