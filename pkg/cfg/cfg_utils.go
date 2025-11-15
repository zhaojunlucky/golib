package cfg

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func GetCfgPath(appName string, cfgFile string) []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home dir: %v", err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current dir: %v", err)
	}
	var cfgPaths []string = []string{
		filepath.Join("/etc", appName, cfgFile),
		filepath.Join(homeDir, ".config", appName, cfgFile),
		filepath.Join(curDir, cfgFile),
	}

	return cfgPaths
}
