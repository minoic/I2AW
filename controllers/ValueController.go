package controllers

import (
	"github.com/MinoIC/I2AW/Database"
	"github.com/astaxie/beego"
)

type ValueController struct {
	beego.Controller
}

func (this *ValueController) Get() {
	identifier := this.Ctx.Input.Param(":identifier")
	DB := Database.GetDatabase()
	var item Database.Item
	DB.First(&item, "identifier = ?", identifier)
	_, _ = this.Ctx.ResponseWriter.Write([]byte(item.Value))
}
