package config

import (
	"github.com/mounlion/markdownwatcher/model"
	"flag"
)

var (
	Config = model.Config{}
	hoursUpdate []int
	botToken string
	dataSource string
	cities = map[string]string{
		"moscow": "Москва",
		"belgorod": "Белгород",
	}
)

func GetConfig()  {
	Debug := flag.Bool("debug", false, "Use debug mode for create and update Mark Down Watcher")
	Logger := flag.Bool("log", false, "Use log for view all processes")
	flag.Parse()

	Config.Debug = Debug
	Config.Logger = Logger

	if *Debug {
		botToken = "__BOT_TOKEN__"
		hoursUpdate = append(hoursUpdate,  9, 11, 12, 13, 14, 15, 16, 17, 18, 21, 22, 23, 0)
		dataSource = "__DB_PATH__"
	} else {
		botToken = "__BOT_TOKEN__"
		hoursUpdate = append(hoursUpdate,  8, 10, 12, 14, 17, 18, 19, 22, 23)
		dataSource = "__DB_PATH__"
	}

	Config.BotToken = &botToken
	Config.HoursUpdate = &hoursUpdate
	Config.DataSource = &dataSource
	Config.Cities = &cities
}