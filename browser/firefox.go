package browser

import (
	"github.com/gitaiqaq/Env.Beta/utils"
	"github.com/shibukawa/configdir"
	"runtime"
)


type Firefox struct {
	name string
}

func (c Firefox) BaseArgs() string {
	return "-no-remote"
}

func (c Firefox) ProgramDir() string {
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

func (c Firefox) Execable() string {
	return "firefox.exe"
}

func (c Firefox) Profile() string {
	return "Env.Beta"
}

func (c Firefox) ProfileDir() string {
	configDirs := configdir.New("Mozilla", "Firefox")
	cache := configDirs.QueryCacheFolder()
	if !cache.Exists("Profiles") {
		cache.CreateParentDir("Profiles")
	}
	utils.CmdRun(c.ProgramDir() + c.Execable(), "-CreateProfile", c.Profile())
	return "--user-data-dir=" + cache.Path
}

func (c Firefox) Incognito() string {
	return "-private"
}

func (c Firefox) ProxyServer(address string) string {
	return "--proxy-server=" + address
}

func (c Firefox) Tpl() string {
	if c.ProfileDir() != "" {
		return "/d {{Escape ProgramDir}} {{Execable}} {{Incognito}} {{ BaseArgs }}"
	}
	return "firefox {{Incognito}} {{ BaseArgs }}"
}