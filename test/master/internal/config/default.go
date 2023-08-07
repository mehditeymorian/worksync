package config

import "github.com/mehditeymorian/worksync/test/master/internal/db"

func Default() Config {
	return Config{
		Database: db.Config{
			Database: "",
			Host:     "",
			Port:     3306,
			Username: "",
			Password: "",
		},
		Jobs: make([]Job, 0),
	}
}
