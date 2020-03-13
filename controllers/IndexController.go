package controllers

import (
	"github.com/astaxie/beego"
	image2 "image"
	_ "image/jpeg"
	_ "image/png"
	"image2ascii/convert"
	"net/http"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Get() {
	options := convert.Options{
		Ratio:           0,
		FixedWidth:      160,
		FixedHeight:     60,
		FitScreen:       false,
		StretchedScreen: false,
		Colored:         false,
		Reversed:        false,
	}
	converter := convert.NewImageConverter()
	imageLink := `https://img.ntmc.tech/images/2020/03/13/mfddpItIKoiH3F0q.jpg`
	imageSrc, _ := http.Get(imageLink)
	image, _, _ := image2.Decode(imageSrc.Body)
	ret := converter.Image2ASCIIString(image, &options)
	_, _ = this.Ctx.ResponseWriter.Write([]byte(ret))
}
