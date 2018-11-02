package browser

import (
	"github.com/gitaiqaq/Env.Beta/utils"
	"github.com/shibukawa/configdir"
	"path/filepath"
	"runtime"
)

type Chrome struct {
	name string
}

func (c *Chrome) BaseArgs() string {
	return "--lang=local"
}

func (c *Chrome) ProgramDir() string {
	var programFile string
	switch runtime.GOOS {
	case "darwin":
		programFile = "/Applications/Google\\ Chrome.app"
	case "windows":
		programFile = "C:\\Program Files\\"
		var tPath = "Google\\Chrome\\Application\\"
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
		}
	}
	return filepath.Dir(programFile)
}

func (c *Chrome) Execable() string {
	return "chrome.exe"
}

func (c *Chrome) Profile() string {
	return ""
}

func (c *Chrome) ProfileDir() string {
	configDirs := configdir.New("Env.Beta", "Chrome")
	cache := configDirs.QueryCacheFolder()
	if !cache.Exists("First Run") {
		cache.CreateParentDir("First Run")
	}
	return "--user-data-dir=\"" + cache.Path + "\""
}

func (c *Chrome) Incognito() string {
	return "-incognito"
}

func (c *Chrome) ProxyServer(address string) string {
	return "--proxy-server=" + address
}

func (c *Chrome) Tpl() string {
	if c.ProfileDir() != "" {
		return "/D {{Escape ProgramDir}} {{Execable}} {{ProfileDir}} {{Incognito}} {{ProxyServer .ProxyServer}} {{ BaseArgs }}"
	}
	return "chrome {{ProfileDir}} {{Incognito}} {{ProxyServer .ProxyServer}} {{ BaseArgs }}"
}
