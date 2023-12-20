package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/leonardom/go-jar-launcher/configs"
	"github.com/leonardom/go-jar-launcher/internal"
	"moul.io/banner"
)

func main() {
	appName := filepath.Base(os.Args[0])
	deleteStaleBackups(appName)
	args := os.Args[1:]
	log.Println(banner.Inline(appName))
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Printf("Working Dir: %v\n", wd)
	configFile := fmt.Sprintf("%v.yaml", appName)
	if len(args) > 0 {
		configFile = args[0]
	}
	config := loadConfigs(configFile)
	checkUpdate(appName, config)
	launchApp(config)
}

func deleteStaleBackups(appName string) {
	files, err := filepath.Glob("backup-" + appName + "*.zip")
	if err != nil {
		return
	}
	if len(files) > 2 {
		sort.Strings(files)
		files = files[:len(files)-2]
		fmt.Println(files)
		for _, f := range files {
			os.Remove(f)
		}
	}
}

func loadConfigs(configFile string) *configs.Config {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		log.Printf("ERROR: Missing config file \"%v\"!\n", configFile)
		os.Exit(1)
	}
	config, err := configs.LoadConfig(configFile)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	return config
}

func checkUpdate(appName string, config *configs.Config) {
	if config.CheckUpdate == "" {
		return
	}
	updater := internal.NewUpdater(appName, config)
	updater.CheckUpdate()
}

func launchApp(config *configs.Config) {
	launcher := internal.NewLauncher(config)
	err := launcher.Execute()
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
