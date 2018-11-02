package browser

import (
	"bufio"
	"fmt"
	"github.com/gitaiqaq/Env.Beta/utils"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type FirefoxProfile struct {
	userjs      string
	preferences map[string]interface{}
}

func (p FirefoxProfile) set_preference(key string, value interface{}) {
	p.preferences[key] = value
}

func (p FirefoxProfile) _write_user_prefs() {
	fi, err := os.OpenFile(p.userjs, os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	for item := range p.preferences {
		value := p.preferences[item]
		_, err := fi.WriteString(`user_pref("` + item + `", `)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		switch value.(type) {
		case string:
			fi.WriteString("\"" + value.(string) + "\"")
			break
		default:
			fi.WriteString(fmt.Sprint(value))
		}
		fi.WriteString(`);
`)
	}

	fi.Sync()

	fi.Close()
}

func ParseValue(value string) interface{} {
	if strings.HasPrefix(value, "\"") {
		return strings.Trim(value, "\"")
	}
	b, err := strconv.ParseBool(value)
	if err == nil {
		return b
	}
	f, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return f
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return i
	}
	return value
}

func (p FirefoxProfile) _read_existing_userjs() {
	re := regexp.MustCompile(`user_pref\("(.*)",\s(.*)\)`)
	fi, err := os.Open(p.userjs)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		matches := re.FindSubmatch(a)
		if len(matches) == 3 {
			log.Println(string(matches[0]))
			log.Println(string(matches[1]), string(matches[2]))
			p.set_preference(string(matches[1]), ParseValue(string(matches[2])))
		}
	}
}

type Firefox struct {
	Name string

	path    string
	profile FirefoxProfile
	Host    string
	Port    uint64
}

func (c *Firefox) BaseArgs() string {
	return "-no-remote"
}

func (c *Firefox) ProgramDir() string {
	var programFile string
	switch runtime.GOOS {
	case "darwin":
		programFile = "/Applications/Firefox.app"
	case "windows":
		programFile = "C:\\Program Files\\"
		var tPath = "Mozilla Firefox\\"
		if !utils.Exists(programFile + tPath) {
			if runtime.GOARCH == "amd64" {
				programFile = "C:\\Program Files (x86)\\"
				if utils.Exists(programFile + tPath) {
					programFile = programFile + tPath
				} else {
					programFile = ""
				}
			} else {
				programFile = ""
			}
		} else {
			programFile = programFile + tPath
		}
	}
	return programFile
}

func (c *Firefox) Execable() string {
	return "firefox.exe"
}

func (c *Firefox) Profile() string {
	return "Env.Beta"
}

func (c *Firefox) ProfileDir() string {
	out := utils.CmdRun(c.ProgramDir()+c.Execable(), "-CreateProfile", c.Profile())

	var path string
	for _, token := range strings.Split(out, "'") {
		if strings.Contains(token, "prefs.js") {
			path = token
		}
	}

	if path != "" {
		c.path = path
	}
	return "--user-data-dir=" + path
}

func (c *Firefox) Incognito() string {
	return "-private"
}

func (c *Firefox) ProxyServer(address string) string {
	adds := strings.Split(address, ":")
	port, _ := strconv.ParseUint(adds[1], 10, 0)

	c.profile = FirefoxProfile{
		userjs:      c.path,
		preferences: make(map[string]interface{}),
	}
	// profile._read_existing_userjs()
	c.profile.set_preference("network.proxy.type", 1)

	c.profile.set_preference("network.proxy.http", adds[0])
	c.profile.set_preference("network.proxy.http_port", port)

	c.profile.set_preference("network.proxy.ssl", adds[0])
	c.profile.set_preference("network.proxy.ssl_port", port)

	c.profile.set_preference("network.proxy.socks", adds[0])
	c.profile.set_preference("network.proxy.socks_port", port)

	c.profile._write_user_prefs()

	return ""
}

func (c *Firefox) Tpl() string {
	if c.ProfileDir() != "" {
		c.ProxyServer("127.0.0.1:2654")
		return "/D {{Escape ProgramDir}} {{ProxyServer .ProxyServer}}{{Execable}} -P Env.Beta "
	}
	return "firefox -P Env.Beta {{ProxyServer .ProxyServer}}{{Incognito}} {{ BaseArgs }}"
}
