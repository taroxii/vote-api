package config

import (
	"log"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/spf13/viper"
)

var Config *Configuration

type Configuration struct {
	ServerPort string         `json:"serverPort" mapstructure:"server-port"`
	JWTSecret  string         `json:"jwtSecret" mapstructure:"jwt-secret"`
	Timeout    int64          `json:"timeout" mapstructure:"timeout"`
	Postgres   PostgresConfig `json:"postgres" mapstructure:"postgres"`
	JWTConfig  echojwt.Config
}

type PostgresConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Database string `json:"database" mapstructure:"database"`
	CertPEM  string `json:"certPem" mapstructure:"ca-pem"`
	CAFile   string `json:"caFile" mapstructure:"ca-file"`
	SSL      bool   `json:"ssl" mapstructure:"ssl"`
}

// LoadConfig loads configuration from environment variables.
func NewConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read in viper config:%s", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

}
