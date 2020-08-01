package controllers

import (
	"github.com/MinoIC/I2AW/database"
	"github.com/astaxie/beego"
)

type RgbValueController struct {
	beego.Controller
}

func (this *RgbValueController) Get() {
	this.TplName = "rgbvalue.html"
	identifier := this.Ctx.Input.Param(":identifier")
	DB := database.GetDatabase()
	var item database.RgbItem
	DB.First(&item, "identifier = ?", identifier)
	this.Data["value"] = item.Value
}
