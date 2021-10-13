package clux

import (
	"log"
	"os"
)

var _version_ = "debug"
var _branch_ = ""
var _commitId_ = ""
var _buildTime_ = ""
var _appName_ = ""
var ns string

func getEnv(name, def string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		v = def
	}
	return v
}
func init() {
	log.Println("["+_appName_+"]", "init ...")
	log.Println("["+_appName_+"]", "version", _version_)
	log.Println("["+_appName_+"]", "branch", _branch_)
	log.Println("["+_appName_+"]", "commit id", _commitId_)
	log.Println("["+_appName_+"]", "build time", _buildTime_)
	ns = getEnv("POD_NAMESPACE", "default")
	log.Println("["+_appName_+"]", "namespace", ns)
}
