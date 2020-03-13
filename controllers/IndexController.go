package controllers

import "github.com/astaxie/beego"
import _ "image2ascii/convert"
import _ "image/png"
import _ "image/jpeg"

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {

}
