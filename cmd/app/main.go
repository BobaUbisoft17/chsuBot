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

	groupDb := database.NewGroupStorage(
		cfg.DatabaseURL,
		logger,
	)

	userDb := database.NewUserStorage(
		cfg.DatabaseURL,
		logger,
	)

	if err := groupDb.Start(); err != nil {
		panic(err)
	}
	if err := userDb.Start(); err != nil {
		panic(err)
	}

	api := chsuAPI.New(
		map[string]string{
			"username": "mobil",
			"password": "ds3m#2nn",
		},
		logger,
	)

	groupIDs, err := api.GroupsId()
	if err != nil {
		panic(err)
	}

	groupDb.AddGroups(groupIDs)

	rel := reload.NewReloader(
		api,
		groupDb,
		logger,
	)
	rel.ReloadSchedule(0)

	b := bot.New(
		api,
		groupDb,
		userDb,
		logger,
		cfg.AdminId,
		cfg.TelegramBotToken,
		cfg.TypeStart,
		cfg.WebhookURL,
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
