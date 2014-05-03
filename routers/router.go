package routers

import (
	"typecho-app-store/controllers"
	"github.com/astaxie/beego"
	"html/template"
	"net/http"
)

func page_not_found(rw http.ResponseWriter, r *http.Request){
	t, _ := template.New("redirect.tpl").ParseFiles(beego.ViewsPath+"/redirect.tpl")
	data := new(interface {})
	t.Execute(rw, data)
}

func init() {
	beego.Errorhandler("404", page_not_found)
    beego.Router("/", &controllers.MainController{})
	beego.Router("/packages.json", &controllers.PluginController{}, "get:Get")
	beego.Router("/plugins", &controllers.PluginController{}, "post:Post")
	beego.Router("/archive/:plugin([A-Za-z0-9]+)/:version(.*)", &controllers.DownloadController{})
}
