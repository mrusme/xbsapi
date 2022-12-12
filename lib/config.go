package lib

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	StatusOnline     int8 = iota
	StatusOffline         = 2
	StatusNoNewSyncs      = 3
)

type ServiceStatus int8

type Config struct {
	Debug string

	Database struct {
		Type       string
		Connection string
	}

	Server struct {
		BindIP string
		Port   string

		Prefork                 bool
		ServerHeader            string
		StrictRouting           bool
		CaseSensitive           bool
		ETag                    bool
		Concurrency             int
		ProxyHeader             string
		EnableTrustedProxyCheck bool
		TrustedProxies          []string
		DisableStartupMessage   bool
		AppName                 string
		ReduceMemoryUsage       bool
		Network                 string
		EnablePrintRoutes       bool
	}

	Service struct {
		Status      ServiceStatus
		Message     string
		MaxSyncSize int
	}
}

func Cfg() (Config, error) {
	viper.SetDefault("Debug", "true")

	viper.SetDefault("Database.Type", "sqlite3")
	viper.SetDefault("Database.Connection", "file:ent?mode=memory&cache=shared&_fk=1")

	viper.SetDefault("Server.BindIP", "0.0.0.0")
	viper.SetDefault("Server.Port", "8000")

	viper.SetDefault("Server.Prefork", "false")
	viper.SetDefault("Server.ServerHeader", "")
	viper.SetDefault("Server.StrictRouting", "false")
	viper.SetDefault("Server.CaseSensitive", "false")
	viper.SetDefault("Server.ETag", "false")
	viper.SetDefault("Server.Concurrency", strconv.Itoa(256*1024))
	viper.SetDefault("Server.ProxyHeader", "")
	viper.SetDefault("Server.EnableTrustedProxyCheck", "false")
	viper.SetDefault("Server.TrustedProxies", "")
	viper.SetDefault("Server.DisableStartupMessage", "true")
	viper.SetDefault("Server.AppName", "xbsapi")
	viper.SetDefault("Server.ReduceMemoryUsage", "false")
	viper.SetDefault("Server.Network", "tcp")
	viper.SetDefault("Server.EnablePrintRoutes", "false")

	viper.SetDefault("Service.Status", strconv.Itoa(int(StatusOnline)))
	viper.SetDefault("Service.Message", "It really whips the llama's ass")
	viper.SetDefault("Service.MaxSyncSize", strconv.Itoa(204800))

	viper.SetConfigName("xbsapi.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$XDG_CONFIG_HOME/")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("xbsapi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	config = *ParseDatabaseURL(&config)

	return config, nil
}

func ParseDatabaseURL(config *Config) *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return config
	}

	dbURL, err := url.Parse(databaseURL)
	if err != nil {
		return config
	}

	host, port, _ := net.SplitHostPort(dbURL.Host)
	dbname := strings.TrimLeft(dbURL.Path, "/")
	user := dbURL.User.Username()
	password, _ := dbURL.User.Password()

	switch dbURL.Scheme {
	case "postgresql", "postgres":
		if port == "" {
			port = "5432"
		}
		config.Database.Type = "postgres"
		config.Database.Connection = fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s",
			host, port, dbname, user, password,
		)
	case "mysql":
		if port == "" {
			port = "3306"
		}
		config.Database.Type = "mysql"
		config.Database.Connection = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=True",
			user, password, host, port, dbname,
		)
	}

	return config
}
