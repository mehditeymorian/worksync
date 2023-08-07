package config

import "github.com/mehditeymorian/worksync/test/master/internal/db"

func Default() Config {
	return Config{
		Name: "",
		Database: db.Config{
			Database: "",
			Host:     "",
			Port:     3306,
			Username: "",
			Password: "",
		},
		Job: Job{
			Name: "",
			Cron: "",
		},
	}
}
