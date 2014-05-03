package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"typecho-app-store/typecho"
	"typecho-app-store/typecho/ziputil"
	"fmt"
	"encoding/json"
)

var (
	Db orm.Ormer
	storagePath	string
	repoPath	string
)

type Plugin struct {
	Id	int	`orm:"column(id);pk"`
	Package	string	`orm:"column(package);unique"`
	CreatedAt	time.Time	`orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdateAt	time.Time	`orm:"column(updated_at);auto_now;type(datetime)"`
	Versions	[]*Version	`orm:"reverse(many)"`
}

func (p *Plugin) TableName() string {
	return "plugins"
}

type Version struct {
	Id	int	`orm:"column(id);pk"`
	Name	string	`orm:"column(name)"`
	Plugin	*Plugin `orm:"rel(fk)"`
	Author	string	`orm:"column(author)"`
	VersionNo	string	`orm:"column(version)"`
	Description	string	`orm:"column(description)"`
	Link	string	`orm:"column(link)"`
	Require	string	`orm::column(require):`
	CreatedAt	time.Time	`orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdateAt	time.Time	`orm:"column(updated_at);auto_now;type(datetime)"`

}

func (v *Version) TableName() string {
	return "plugin_versions"
}

func init() {
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	dbUser := beego.AppConfig.String("mysqluser")
	dbPassword := beego.AppConfig.String("mysqlpass")
	dbHost := beego.AppConfig.String("mysqlhost")
	dbName := beego.AppConfig.String("mysqldb")
	orm.RegisterDataBase("default", "mysql", dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName+"?charset=utf8")
	orm.RegisterModel(new(Plugin), new(Version))
	Db = orm.NewOrm()
}

func UpdatePlugin(p *typecho.Plugin) {
	if IsPluginExisted(p.Package) {
		AppendVersion(p)
	} else {
		AddNewPlugin(p)
	}
}

func GetAllPlugins() []*Plugin {
	var plugins []*Plugin
	pluginModel := new(Plugin)
	_, err := Db.QueryTable(pluginModel).All(&plugins)
	if err != nil {
		return plugins
	}
	for _, plugin := range plugins {
		Db.LoadRelated(plugin, "Versions", 1, 0, 0, "-version")
	}

	return plugins
}

func AddNewPlugin(p *typecho.Plugin) {
	plugin := new(Plugin)
	plugin.Package = p.Package
	id, err := Db.Insert(plugin)
	if err != nil {
		fmt.Printf("Insert Fail %s \n", p)
	} else {
		fmt.Printf("Insert OK %s - %s\n", p, id)
		AppendVersion(p)
	}
}

func AppendVersion(p *typecho.Plugin) {
	plugin := FindPlugin(p.Package)
	if ! IsVersionExisted(plugin.Id, p.Version) {
		version := 	new(Version)
		version.Plugin = plugin
		version.Author = p.Author
		version.Name = p.Name
		version.Description = p.Description
		version.Link = p.Link
		version.Require = p.Require
		version.VersionNo = p.Version
		_, err := Db.Insert(version)
		if err == nil {
			zipPlugin(p.Package, p.Version, p.Source)
		}
	}
}

func FindPlugin(name string) *Plugin {
	plugin := new(Plugin)
	plugin.Package = name
	Db.Read(plugin, "Package")
	return plugin
}

func IsPluginExisted(name string) bool {
	plugin := new(Plugin)
	return Db.QueryTable(plugin).Filter("package", name).Exist()
}

func IsVersionExisted(pluginId int, v string) bool {
	version := new(Version)
	return Db.QueryTable(version).Filter("plugin_id", pluginId).Filter("version", v).Exist()
}

func SetStoragePath(path string) {
	storagePath = path
}

func SetRepoPath(path string) {
	repoPath = path
}

func zipPlugin(name, version, repo string) {
	zipFile := storagePath + "archive/" + name + "/" + name + "-" + version + ".zip"
	var directory string
	if name == "CommenToMail4BAE" {
		directory = repoPath + repo + "/" + name + "/CommentToMail/"
	} else if name == "Contribute" {
		directory = repoPath + repo + "/" + name + "/plugins/Contribute/"
	} else {
		directory = repoPath + repo + "/" + name + "/"
	}
	err := ziputil.Zip(zipFile, directory)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

type VersionJson struct {
	Version	string	`json:"version"`
	Author	string	`json:"author"`
	Description	string	`json:"description"`
	Link	string	`json:"link"`
	Require	string	`json:"require"`
	CreatedAt	time.Time	`json:"created_at"`
}

type PluginJson struct {
	Name	string	`json:"name"`
	Versions	[]VersionJson	`json:"versions"`
}

type PluginsJson struct {
	Packages	[]PluginJson	`json:"packages"`
}

func ToJson() ([]byte, error) {
	var pluginsJson PluginsJson
	plugins := GetAllPlugins()
	for _, plugin := range plugins {
		var pluginJson PluginJson
		pluginJson.Name = plugin.Package

		for _, version := range plugin.Versions {
			versionJson := VersionJson{}
			versionJson.Version = version.VersionNo
			versionJson.Author = version.Author
			versionJson.CreatedAt = version.CreatedAt
			versionJson.Description = version.Description
			versionJson.Link = version.Link
			versionJson.Require = version.Require

			pluginJson.Versions = append(pluginJson.Versions, versionJson)
		}

		pluginsJson.Packages = append(pluginsJson.Packages, pluginJson)
	}

	bytes, err := json.Marshal(pluginsJson)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
