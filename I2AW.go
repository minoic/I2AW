package main

import (
	"github.com/MinoIC/I2AW/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/rgbvalue/:identifier", &controllers.RgbValueController{})
	beego.Run()
}
