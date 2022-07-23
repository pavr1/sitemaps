package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ExecutionMode int32  `env:"EXECUTION_MODE" envDefault:"1"`
	AppName       string `env:"APP_NAME" envDefault:"Sitemaps"`
	HttpPort      int    `env:"HTTP_PORT" envDefault:"8080"`
	UrlRegExpr    string `env:"URL_REG_EXPR" envDefault:"<a.*?href=\"(.*?)\""`
	OutputFile    string `env:"OUTPUT_FILE" envDefault:"%s.xml"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
