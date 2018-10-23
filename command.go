// Copyright 2017 Gitai<i@gitai.me> All rights reserved.
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall
// be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
// ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"runtime"
	"github.com/shibukawa/configdir"
	"os"
	"text/template"
	"strings"
	"log"
	"os/exec"
	"path/filepath"
	"net"
)

// https://peter.sh/experiments/chromium-command-line-switches/

type Command struct {
	Program		string
	Execer      []string
	BrowserPath string
	UserDataDir string
	ProxyServer string
	Lang string
}

func getDataDir(program string) string {
	configDirs := configdir.New("Env.Beta", program)
	cache := configDirs.QueryCacheFolder()
	if !cache.Exists("First Run") {
		cache.CreateParentDir("First Run")
	}
	return cache.Path
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (c *Command) SetUserDataDir() {
	c.UserDataDir = getDataDir(c.Program)
}

func (c *Command) SetExecer() {
	switch runtime.GOOS {
	case "darwin":
		c.Execer = []string{"open"}
	case "windows":
		if c.BrowserPath != c.Program {
			c.Execer = []string{"cmd", "/c", "start", "/d", filepath.Dir(c.BrowserPath) }
			c.BrowserPath = filepath.Base(c.BrowserPath)
		} else {
			c.Execer = []string{"cmd", "/c", "start"}
		}
	default:
		c.Execer = []string{"xdg-open"}
	}
}

func (c *Command) SetBrowser(program string) {
	c.Program = program
}

func (c *Command) Build(program string)  {
	c.SetBrowser(program)
	c.SetBrowserPath()
	c.SetExecer()
	c.SetUserDataDir()
}

func _GetBrowserPathWindows(arch string, program string) string{
	var programFile = "C:\\Program Files\\"
	if arch == "amd64" {
		programFile = "C:\\Program Files (x86)\\"
	}

	switch program {
	case "chrome":
		programFile = programFile + "Google\\Chrome\\Application\\chrome.exe"
		break
	case "firefox":
		programFile = programFile + "Mozilla Firefox\\firefox.exe"
		break
	}

	return programFile
}

func _GetBrowserPathMac(program string) string {
	var programFile = ""
	switch program {
	case "chrome":
		programFile = "/Applications/Google\\ Chrome.app"
		break
	case "firefox":
		programFile = "/Applications/Firefox.app"
	}
	return programFile
}

func _GetBrowserPathLinux(arch string, program string) {

}

func _GetBrowserPath(os string, arch string, program string) string {
	var programFile = ""
	switch os {
	case "darwin" :
		programFile = _GetBrowserPathMac(program)
	case "windows":
		programFile = _GetBrowserPathWindows(arch, program)

		if arch == "amd64" && !Exists(programFile) {
			programFile = _GetBrowserPathWindows("386", program)
		}

		if program == "chrome" && !Exists(programFile) {
			programFile = "chrome"
		}
	}
	return programFile
}

func (c *Command) SetBrowserPath() {
	c.BrowserPath = _GetBrowserPath(runtime.GOOS, runtime.GOARCH, c.Program)
}

func (c *Command) SetProxyServer(proxy interface{}) {
	switch proxy.(type) {
	case string:
		c.ProxyServer = proxy.(string)
	case net.Addr:
		c.ProxyServer = proxy.(net.Addr).String()
	}
}

func (c *Command) SetLang(lang string)  {
	c.Lang = lang
}

func (c *Command) String(appTpl string) string {
	switch appTpl {
	case "chrome":
		appTpl = "{{ StringsJoin .Execer \" \" }} {{ .BrowserPath }} --user-data-dir={{ .UserDataDir }} --proxy-server={{ .ProxyServer }} --lang={{ .Lang }}"
	case "firefox":
		appTpl = "{{ StringsJoin .Execer \" \" }} {{ .BrowserPath }} --user-data-dir={{ .UserDataDir }} --proxy-server={{ .ProxyServer }} --lang={{ .Lang }}"
	}

	tmpl, err := template.New("command").Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
	}).Parse(appTpl)

	if err != nil {
		panic(err)
	}
	var b = &strings.Builder{}

	err = tmpl.Execute(b, c)
	if err != nil {
		print(err)
	}

	return b.String()
}

func (c *Command) Start(urls []string)  {
	var args = c.Execer

	switch c.Program {
	case "chrome":
		args = append(args, c.BrowserPath, "--user-data-dir=" + c.UserDataDir, "--proxy-server=" + c.ProxyServer, "--in-process-plugins", "-incognito", "--allow-running-insecure-content", "--lang=" + c.Lang)
	}

	cmd := exec.Command(args[0], append(args[1:], urls...)...)

	log.Println(cmd)
	cmd.Start()
}

func test() {
	var command = Command{}
	var program = "chrome"
	command.Build(program)
	command.SetProxyServer("localhost:8080")
	command.SetLang("local")
	command.Start([]string{"about:version"})
}