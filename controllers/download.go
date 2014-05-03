package controllers

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"typecho-app-store/typecho/logger"
)

type DownloadController struct {
	beego.Controller
}

func (this *DownloadController) Get() {
	plugin := this.GetString(":plugin")
	version := this.GetString(":version")
	file := beego.AppConfig.String("storagePath")+"/archive/"+plugin+"/"+plugin+"-"+version+".zip"
	logger.Log(this.Ctx.Request.Host, this.Ctx.Request.Referer(), "download", plugin+"/"+version)
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		this.Ctx.Output.SetStatus(404)
		this.StopRun()
	}
	this.Ctx.Output.Header("Content-Disposition", `attachment; filename="`+plugin+"-"+version+`.zip"`)
	this.Ctx.Output.Body(bytes)
	this.StopRun()
}
