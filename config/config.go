package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logs_format string `mapstructure:"logs_format"`
	Listen      struct {
		BindIp       string `mapstructure:"bind_ip"`
		Port         string `mapstructure:"port"`
		WriteTimeout int    `mapstructure:"write_timeout"`
		ReadTimeout  int    `mapstructure:"read_timeout"`
	} `mapstructure:"listen"`
	Database struct {
		Host     string `mapstructure:"db_host"`
		Port     string `mapstructure:"db_port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"db_name"`
		SSLMode  string `mapstructure:"ssl_mode"`
	} `mapstructure:"database"`
	Redis struct {
		RedisHost string `mapstructure:"redis_host"`
		RedisPort string `mapstructure:"redis_port"`
		RedisDB   int    `mapstructure:"redis_db"`
		LifeTime  int    `mapstructure:"life_time"`
	} `mapstructure:"cache"`
}

func GetConfig() *Config {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}

func (cfg *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.DBName,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)
}

func (cfg *Config) GetAdress() string {
	return fmt.Sprintf(
		"%s:%s",
		cfg.Listen.BindIp,
		cfg.Listen.Port,
	)
}
