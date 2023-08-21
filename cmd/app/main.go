package main

import (
	"sync"

	"github.com/BobaUbisoft17/chsuBot/internal/bot"
	chsuAPI "github.com/BobaUbisoft17/chsuBot/internal/chsuAPI"
	"github.com/BobaUbisoft17/chsuBot/internal/config"
	"github.com/BobaUbisoft17/chsuBot/internal/database"
	reload "github.com/BobaUbisoft17/chsuBot/internal/reloadSchedule"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.New()
	db := database.NewStorage(
		cfg.DatabaseURL,
		logger,
	)
	db.Start()

	API := chsuAPI.New(
		map[string]string{
			"username": "mobil",
			"password": "ds3m#2nn",
		},
		logger,
	)

	groupIDs, err := API.GroupsId()
	if err != nil {
		panic(err)
	}
	db.AddGroups(groupIDs)

	rel := reload.NewReloader(
		API,
		db,
		logger,
	)
	rel.ReloadSchedule(0)

	b := bot.New(
		API,
		db,
		logger,
		cfg.TelegramBotToken,
	)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.StartBot()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		rel.Start()
	}()
	wg.Wait()
}
