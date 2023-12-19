package internal

import (
	"os"
	"testing"

	"github.com/leonardom/go-jar-launcher/configs"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	cfg, err := configs.LoadConfig(wd + "/../test/config/app.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, cfg)
	launcher := NewLauncher(cfg)
	assert.NotNil(t, launcher)
	cmd := launcher.getCommand()
	expectedArgs := []string{"-Dfoo=bar", "-Xmx1g", "-jar", "example.jar", "arg1", "arg2", "arg3"}
	expectedCmd := Command{
		Name: "./jre/bin/java",
		Args: expectedArgs,
	}
	assert.Equal(t, expectedCmd, cmd)
}
