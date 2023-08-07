package db

import "fmt"

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
}

func (c Config) DNS() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?interpolateParams=false&parseTime=true&charset=utf8mb4",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
