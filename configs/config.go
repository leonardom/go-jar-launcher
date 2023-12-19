package configs

import "github.com/spf13/viper"

type Config struct {
	JavaHome   string
	JVMOptions []string
	JARFile    string
	Args       []string
}

func LoadConfig(file string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg *Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
