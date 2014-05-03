package typecho

import (
	"regexp"
	"io/ioutil"
	"strings"
	"os"
)

type Plugin struct {
	Package	string
	Name	string
	Description	string
	Author	string
	Version	string
	Link	string
	Require	string
	Source	string
}

func Parse(path, packageName, repo string, retry bool) Plugin{
	plugin := Plugin{"", "", "", "", "", "", "*", ""}
	pluginContent, err := ioutil.ReadFile(path)
	if err != nil {
		if retry {
			os.Rename(strings.Replace(path, "Plugin.php", "plugin.php", 1), path)
			plugin = Parse(path, packageName, repo, false)
		}
		return plugin
	}
	reString := `/\*\*([\s\S]*?)\*/`
	re, _ := regexp.Compile(reString)
	matches := re.FindAllString(string(pluginContent), -1)
	wantedMatch := ""
	for _, match := range matches {
		if strings.Contains(match, "@package") {
			wantedMatch = match
			break
		}
	}

	lines := strings.Split(wantedMatch, "\n")
	for _, line := range lines {
		line = strings.Replace(line, "/**", "", -1)
		line = strings.Replace(line, "*/", "", 1)
		line = strings.Replace(line, "*", "", 1)
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "@package") {
			plugin.Name = strings.Trim(strings.Replace(line, "@package", "", 1), " ")
			continue
		}
		if strings.HasPrefix(line, "@author") {
			plugin.Author = strings.Trim(strings.Replace(line, "@author", "", 1), " ")
			continue
		}
		if strings.HasPrefix(line, "@version") {
			plugin.Version = strings.Trim(strings.Replace(line, "@version", "", 1), " ")
			continue
		}
		if strings.HasPrefix(line, "@link") {
			plugin.Link = strings.Trim(strings.Replace(line, "@link", "", 1), " ")
			continue
		}
		if strings.HasPrefix(line, "@dependence") {
			plugin.Require = strings.Trim(strings.Replace(line, "@dependence", "", 1), " ")
			continue
		}
		if ! strings.HasPrefix(line, "@") {
			plugin.Description = plugin.Description + line
		}
	}
	plugin.Source = repo
	plugin.Package = packageName
	return plugin
}
