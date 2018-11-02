package utils

import (
	"log"
	"os"
	"os/exec"
)

func CmdRun(path string, args ...string) string {
	cmd := exec.Command(path, args...)
	log.Println(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Run failed with %s", err)
		os.Exit(1)
	}
	log.Println(string(out))
	return string(out)
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
