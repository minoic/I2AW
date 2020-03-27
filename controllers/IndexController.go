package controllers

import (
	"bytes"
	"github.com/MinoIC/I2AW/Database"
	"github.com/MinoIC/I2AW/image2ascii/convert"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"html/template"
	image2 "image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

type IndexController struct {
	beego.Controller
}

func init() {
	_ = os.Mkdir("./imgcache", os.ModePerm)
}

func (this *IndexController) Get() {
	this.TplName = "index.html"
	this.Data["xsrfData"] = template.HTML(this.XSRFFormHTML())
	method := this.GetString("method")
	if method == "ls" {
		DB := Database.GetDatabase()
		var items []Database.RgbItem
		DB.Find(&items, "session_id = ?", this.StartSession().SessionID())
		this.Data["json"] = items
		this.ServeJSON()
	} else if method == "size" {
		if size := this.StartSession().Get("size"); size == nil {
			_ = this.StartSession().Set("size", 120)
			Database.AddSessionAmount()
			_, _ = this.Ctx.ResponseWriter.Write([]byte("120"))
		} else {
			_, _ = this.Ctx.ResponseWriter.Write([]byte(strconv.Itoa(size.(int))))
		}
		return
	} else if method == "stats" {
		this.Data["json"] = Database.GetStats()
		this.ServeJSON()
	}
}

func (this *IndexController) Post() {
	if !this.CheckXSRFCookie() {
		this.Ctx.ResponseWriter.WriteHeader(401)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("xsrf 检测失败"))
		return
	}
	size, err := this.GetInt("size")
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(412)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("获取尺寸失败"))
		return
	}
	if size > 140 || size < 10 {
		this.Ctx.ResponseWriter.WriteHeader(412)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("输入了不合理的尺寸"))
		return
	}
	_ = this.StartSession().Set("size", size)
	// beego.Info(size)
	file, fileHeader, err := this.GetFile("img")
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(412)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("图片上传失败"))
		return
	}
	// beego.Info(file)
	key := RandKey(10)
	imgCache, err := os.Create("./imgcache/" + key + "_" + fileHeader.Filename)
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("服务器缓存图片失败"))
		return
	}
	_, err = io.Copy(imgCache, file)
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("服务器缓存图片失败"))
		return
	}
	imgSrc, err := os.Open("./imgcache/" + key + "_" + fileHeader.Filename)
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("服务器读取缓存失败"))
		return
	}
	img, _, err := image2.Decode(imgSrc)
	if err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("图片解析失败"))
		return
	}
	options := convert.Options{
		Ratio:           0,
		FixedWidth:      img.Bounds().Max.X * size * 2 / (img.Bounds().Max.Y + img.Bounds().Max.X),
		FixedHeight:     img.Bounds().Max.Y * size / (img.Bounds().Max.Y + img.Bounds().Max.X),
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
	// beego.Debug(matrix)
	/*	var value string*/
	var buf bytes.Buffer
	for x := range matrix {
		// beego.Debug(len(matrix[x]))
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
				for ; y < len(matrix[x])-1 && matrix[x][y+1] == matrix[x][y]; y = y + 1 {
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
		SessionID:  this.StartSession().SessionID(),
		Value:      template.HTML(buf.String()),
		SrcHeight:  img.Bounds().Max.Y,
		SrcWidth:   img.Bounds().Max.X,
		DstHeight:  options.FixedHeight,
		DstWidth:   options.FixedWidth,
	}
	// beego.Debug(buf.String())
	if err := DB.Create(&item).Error; err != nil {
		beego.Error(err)
		this.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = this.Ctx.ResponseWriter.Write([]byte("数据库处理失败"))
		return
	}
	Database.AddItemAmount()
	_, _ = this.Ctx.ResponseWriter.Write([]byte("处理成功！"))
	defer runtime.GC()
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
