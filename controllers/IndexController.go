package controllers

import (
	"bytes"
	"github.com/MinoIC/I2AW/Database"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"google.golang.org/appengine/runtime"
	"html/template"
	image2 "image"
	_ "image/jpeg"
	_ "image/png"
	"image2ascii/convert"
	"io"
	"math/rand"
	"os"
	runtime2 "runtime"
	"strconv"
	"time"
)

type IndexController struct {
	beego.Controller
}

func (this *IndexController) Prepare() {
	this.TplName = "index.html"
	this.Data["xsrfData"] = template.HTML(this.XSRFFormHTML())
}

func (this *IndexController) Get() {}

func (this *IndexController) Post() {
	if !this.CheckXSRFCookie() {
		_, _ = this.Ctx.ResponseWriter.Write([]byte("xsrf 检测失败"))
		return
	}
	file, fileHeader, err := this.GetFile("img")
	if err != nil {
		beego.Error(err)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("图片上传失败"))
		return
	}
	beego.Info(file)
	key := RandKey(10)
	imgCache, err := os.Create("./imgcache/" + key + "_" + fileHeader.Filename)
	if err != nil {
		beego.Error(err)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("服务器缓存图片失败"))
		return
	}
	_, err = io.Copy(imgCache, file)
	if err != nil {
		beego.Error(err)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("服务器缓存图片失败"))
		return
	}
	imgSrc, err := os.Open("./imgcache/" + key + "_" + fileHeader.Filename)
	if err != nil {
		beego.Error(err)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("服务器读取缓存失败"))
		return
	}
	img, _, err := image2.Decode(imgSrc)
	if err != nil {
		beego.Error(err)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("图片解析失败"))
		return
	}
	options := convert.Options{
		Ratio:           0,
		FixedWidth:      img.Bounds().Max.X * 240 / (img.Bounds().Max.Y + img.Bounds().Max.X),
		FixedHeight:     img.Bounds().Max.Y * 120 / (img.Bounds().Max.Y + img.Bounds().Max.X),
		FitScreen:       false,
		StretchedScreen: false,
		Colored:         false,
		Reversed:        false,
	}
	beego.Info("new post: ", fileHeader.Filename, "size: ", img.Bounds().Max.X, "x", img.Bounds().Max.Y, " -> ",
		options.FixedWidth, "x", options.FixedHeight)
	DB := Database.GetDatabase()
	converter := convert.NewImageConverter()
	/*	item := Database.Item{
		Model:      gorm.Model{},
		FileName:   key + "_" + fileHeader.Filename,
		Identifier: key,
		Value:      converter.Image2ASCIIString(img, &options),
	}*/
	matrix := converter.Image2CharPixelMatrix(img, &options)
	//beego.Debug(matrix)
	/*	var value string*/
	var buf bytes.Buffer
	for x := range matrix {
		//beego.Debug(len(matrix[x]))
		for y := 0; y < len(matrix[x]); y = y + 1 {
			/*beego.Debug(x,y,matrix[x][y])*/
			if matrix[x][y].Char == ' ' {
				buf.WriteString(`<font style="color:rgb(`)
				buf.WriteString(strconv.Itoa(int(matrix[x][y].R)))
				buf.WriteString(",")
				buf.WriteString(strconv.Itoa(int(matrix[x][y].G)))
				buf.WriteString(",")
				buf.WriteString(strconv.Itoa(int(matrix[x][y].B)))
				buf.WriteString(`)">`)
				buf.WriteString(`&ensp;`)
				for ; y < len(matrix[x])-1 && matrix[x][y+1].Char == ' '; y = y + 1 {
					buf.WriteString(`&ensp;`)
				}
				buf.WriteString(`</font>`)
				/*				value = value + `<font style="color:rgb(` + strconv.Itoa(int(matrix[x][y].R)) + "," +
								"" + strconv.Itoa(int(matrix[x][y].G)) + "," + strconv.Itoa(int(matrix[x][y].B)) + `)">` + "&ensp;" + `</font>`*/
			} else {
				buf.WriteString(`<font style="color:rgb(`)
				buf.WriteString(strconv.Itoa(int(matrix[x][y].R)))
				buf.WriteString(",")
				buf.WriteString(strconv.Itoa(int(matrix[x][y].G)))
				buf.WriteString(",")
				buf.WriteString(strconv.Itoa(int(matrix[x][y].B)))
				buf.WriteString(`)">`)
				buf.WriteByte(matrix[x][y].Char)
				for ; y < len(matrix[x])-1 && matrix[x][y+1].Char == matrix[x][y].Char; y = y + 1 {
					buf.WriteByte(matrix[x][y].Char)
				}
				buf.WriteString(`</font>`)
				/*				value = value + `<font style="color:rgb(` + strconv.Itoa(int(matrix[x][y].R)) + "," +
								"" + strconv.Itoa(int(matrix[x][y].G)) + "," + strconv.Itoa(int(matrix[x][y].B)) + `)">` + string(matrix[x][y].Char) + `</font>`*/
			}
		}
		buf.WriteString(`<br>`)
	}
	item := Database.RgbItem{
		Model:      gorm.Model{},
		FileName:   key + "_" + fileHeader.Filename,
		Identifier: key,
		Value:      template.HTML(buf.String()),
	}
	beego.Debug(buf.String())
	if err := DB.Create(&item).Error; err != nil {
		beego.Error(err)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("图片处理失败"))
		return
	}
	this.Redirect("/rgbvalue/"+item.Identifier, 302)

	runtime.Stats()

}

func (this *IndexController) CheckXSRFCookie() bool {
	if !this.EnableXSRF {
		return true
	}
	token := this.Ctx.Input.Query("_xsrf")
	if token == "" {
		token = this.Ctx.Request.Header.Get("X-Xsrftoken")
	}
	if token == "" {
		token = this.Ctx.Request.Header.Get("X-Csrftoken")
	}
	if token == "" {
		return false
	}
	if this.XSRFToken() != token {
		return false
	}
	return true
}

func RandKey(keyLength int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := []byte(str)
	var ret []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= keyLength; i++ {
		ret = append(ret, b[r.Intn(len(str))])
	}
	return string(ret)
}
