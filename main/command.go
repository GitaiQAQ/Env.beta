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
	"github.com/gitaiqaq/Env.Beta/utils"
	"net"
	"runtime"
	"strings"
	"syscall"
	"text/template"
	"github.com/gitaiqaq/Env.Beta/browser"
	"github.com/kballard/go-shellquote"
	"log"
	"os"
)

// https://peter.sh/experiments/chromium-command-line-switches/

type Command struct {
	Execer      []string
	Browser     browser.IBrowser
	ProxyServer string
}

func (c *Command) SetExecer() {
	switch runtime.GOOS {
	case "darwin":
		c.Execer = []string{"open"}
	case "windows":
		c.Execer = []string{"cmd", "/c", "start"}
	default:
		c.Execer = []string{"xdg-open"}
	}
}

func (c *Command) SetBrowser(program string) {
	switch program {
	case "chrome":
		c.Browser = &browser.Chrome{}
		break
	case "firefox":
		c.Browser = &browser.Firefox{}
		break
	}
	log.Println("Use browser " + program)
}

func (c *Command) Init(program string) {
	c.SetExecer()
	c.SetBrowser(program)
}

func (c *Command) SetProxyServer(proxy interface{}) {
	switch proxy.(type) {
	case string:
		c.ProxyServer = proxy.(string)
	case net.Addr:
		c.ProxyServer = proxy.(net.Addr).String()
	}
}

func (c *Command) String(appTpl string) string {
	if appTpl == "" {
		appTpl = c.Browser.Tpl()
	}

	var funs = template.FuncMap{
		"StringsJoin": strings.Join,
		"Escape":      syscall.EscapeArg,
		"BaseArgs":		c.Browser.BaseArgs,
		"ProgramDir":		c.Browser.ProgramDir,
		"Execable":		c.Browser.Execable,
		"Profile":		c.Browser.Profile,
		"ProfileDir":		c.Browser.ProfileDir,
		"Incognito":		c.Browser.Incognito,
		"ProxyServer":		c.Browser.ProxyServer,
	}

	tmpl, err := template.New("command").Funcs(funs).Parse(appTpl)

	if err != nil {
		panic(err)
	}
	var b = &strings.Builder{}

	err = tmpl.Execute(b, c)
	if err != nil {
		panic(err)
	}

	return b.String()
}

func (c *Command) Start(urls []string) {
	var args = c.Execer
	cmd := c.String("")
	words, err := shellquote.Split(cmd)
	if err != nil {
		log.Println("Parse cmd error: " + cmd)
		os.Exit(1)
	}
	args = append(args, words...)
	utils.CmdRun(args[0], append(args[1:], urls...)...)
}

func test() {
	var command = Command{}
	var program = "chrome"
	command.Init(program)
	command.SetProxyServer("localhost:8080")
	command.Start([]string{"www.baidu.com"})
}
