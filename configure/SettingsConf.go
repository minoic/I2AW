package configure

import (
	"github.com/astaxie/beego/config"
)

var conf config.Configer

var (
	WebApplicationName string
	WebHostName        string
	AdminAddress       string
)

func init() {
	var err error
	conf, err = config.NewConfig("ini", "conf/settings.conf")
	if err != nil {
		panic("cant get settings.conf: " + err.Error())
	}
	ReloadConfig()
}

func ReloadConfig() {
	var err error
	WebApplicationName = conf.String("WebApplicationName")
	AdminAddress = conf.String("WebAdminAddress")
	secure, err := conf.Bool("WebSecure")
	if err != nil {
		panic(err)
	}
	if secure {
		WebHostName = "https://" + conf.String("WebHostName")
	} else {
		WebHostName = "http://" + conf.String("WebHostName")
	}
}

func GetConf() config.Configer {
	return conf
}
