package controllers

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"typecho-app-store/models"
	"typecho-app-store/typecho/logger"
)

type PluginController struct {
	beego.Controller
}

func (this *PluginController) Get() {
	logger.Log(this.Ctx.Request.Host, this.Ctx.Request.Referer(), "packages", "")
	this.Ctx.Output.Header("Content-Type", "application/json")
	json, err := ioutil.ReadFile(beego.AppConfig.String("storagePath")+"packages/packages.json")
	if err != nil {
		this.Ctx.Output.Body([]byte(`{"packages": []}`))
	} else {
		this.Ctx.Output.Body(json)
	}
	this.StopRun()
}

func (this *PluginController) Post() {
	name := this.GetString("name")
	logger.Log(this.Ctx.Request.Host, this.Ctx.Request.Referer(), "check", name)
	result := struct {
		Existed	bool	`json:"existed"`
		}{true}
	if name == "" {
		result.Existed = true
	} else {
		if models.IsPluginExisted(name) {
			result.Existed = true
		} else {
			result.Existed = false
		}
	}
	this.Data["json"] = result
	this.ServeJson()
	this.StopRun()
}
