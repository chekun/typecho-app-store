package logger

import (
	"github.com/astaxie/beego/logs"
)

var log *logs.BeeLogger

func init() {
	log = logs.NewLogger(10000)
	log.SetLogger("file", `{
		"filename": "storage/logs/access.log",
		"maxdays": 365
	}`)
}

func Log(ip, referer, category, extra string) {
	log.Info("["+category+"] ["+ip+"] ["+referer+"]" + " " + extra)
}

