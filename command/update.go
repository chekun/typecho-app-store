package main

import (
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"strings"
	"log"
	"github.com/astaxie/beego"
	"typecho-app-store/typecho"
	"typecho-app-store/models"
)

func runCommand(cmd *exec.Cmd) {
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%s\n", string(stdout))
}

func main() {

	//获取github库名
	repositories := strings.Split(beego.AppConfig.String("repositories"), "|")
	//获取仓库存储目录
	repositoryPath := beego.AppConfig.String("repositoryPath")
	models.SetRepoPath(repositoryPath)
	//获取压缩包存储目录
	storagePath := beego.AppConfig.String("storagePath")
	models.SetStoragePath(storagePath)

	os.Chdir(repositoryPath)
	for _, repo := range repositories {
		repoFolder := strings.Replace(repo, "/", "-", -1)
		_, err := os.Stat(repoFolder)
		if (err != nil) {
			if os.IsNotExist(err) {
				//git clone
				cmd := exec.Command("git", "clone", "https://github.com/"+repo+".git", repoFolder)
				runCommand(cmd)
				fmt.Println("git clone " + repoFolder)
				//处理万恶的gitsubmodule
				os.Chdir(repoFolder)
				cmd = exec.Command("git", "submodule", "update", "--init")
				runCommand(cmd)
				fmt.Println("git submodule update --init")
				os.Chdir("../")
			}
		}
		os.Chdir(repoFolder)
		//到这里就已经有了仓库目录了，那么那么进行一次pull操作
		cmd := exec.Command("git", "stash")
		runCommand(cmd)
		cmd = exec.Command("git", "pull", "origin", "master")
		runCommand(cmd)
		fmt.Println("git pull " + repoFolder)
		//更新万恶的gitsubmodule
		cmd = exec.Command("git", "submodule", "update", "--init")
		runCommand(cmd)
		cmd = exec.Command("git", "submodule", "foreach", "git", "pull", "origin", "master")
		runCommand(cmd)
		fmt.Println("git submodule foreach git pull origin master")
		os.Chdir("../")
		//到这里我们就有了所有的插件了
		files, _ := ioutil.ReadDir(repoFolder)
		for _, file := range files {
			fileName := file.Name()
			if fileName == ".gitignore" || fileName == ".git" || fileName == ".gitattributes" || fileName == ".gitmodules" {
				continue
			}
			if file.IsDir() {
				//常规插件
				//解析Plugin.php中的信息，然后打包待下载
				fmt.Printf("standard plugin %s \n", fileName)
				plugin := typecho.Plugin{}
				if fileName == "CommenToMail4BAE" {
					plugin = typecho.Parse(repoFolder+"/"+fileName+"/CommentToMail/Plugin.php", fileName, repoFolder, true)
				} else if fileName == "Contribute" {
					plugin = typecho.Parse(repoFolder+"/"+fileName+"/plugins/Contribute/Plugin.php", fileName, repoFolder, true)
				} else {
					plugin = typecho.Parse(repoFolder+"/"+fileName+"/Plugin.php", fileName, repoFolder, true)
				}
				fmt.Printf("Detail Info: %s \n", plugin)
				models.UpdatePlugin(&plugin)
				continue
			}
			if strings.Contains(fileName, ".php") {
				//这里是单个文件形式的插件
				//要建成目录，并重命名为Plugin.php
				pluginName := strings.Replace(fileName, ".php", "", 1)
				os.Mkdir(repoFolder+"/"+pluginName, 0755)
				os.Rename(repoFolder+"/"+fileName, repoFolder+"/"+pluginName+"/Plugin.php")
				plugin := typecho.Parse(repoFolder + "/" + pluginName + "/Plugin.php", pluginName, repoFolder, true)
				models.UpdatePlugin(&plugin)
			}
		}
	}

	//获取所有的数据，拼接成json
	json, err := models.ToJson()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	ioutil.WriteFile(storagePath + "packages/packages.json", json, 0755)

	fmt.Println("Task Finish!\n")

}

