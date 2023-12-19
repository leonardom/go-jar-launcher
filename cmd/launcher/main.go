package main

import (
	"fmt"
	"os"

	"github.com/leonardom/go-jar-launcher/configs"
	"github.com/leonardom/go-jar-launcher/internal"
	"moul.io/banner"
)

const CONFIG_FILE = "app.yaml"

func main() {
	args := os.Args[1:]
	fmt.Println(banner.Inline("go jar launcher"))
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nWorking Dir: %v\n", wd)
	configFile := CONFIG_FILE
	if len(args) > 0 {
		configFile = args[0]
	}
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		fmt.Printf("ERROR: Missing config file \"%v\"!\n", configFile)
		os.Exit(1)
	}
	config, err := configs.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	launcher := internal.NewLauncher(config)
	err = launcher.Execute()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
