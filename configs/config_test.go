package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	wd, err := os.Getwd()
	assert.Nil(t, err)
	cfg, err := LoadConfig(wd + "/../test/config/app.yaml")
	assert.Nil(t, err)
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "./jre", cfg.JavaHome)
		assert.ElementsMatch(t, []string{"-Dfoo=bar", "-Xmx1g"}, cfg.JVMOptions)
		assert.Equal(t, "example.jar", cfg.JARFile)
		assert.ElementsMatch(t, []string{"arg1", "arg2", "arg3"}, cfg.Args)
	}
}
