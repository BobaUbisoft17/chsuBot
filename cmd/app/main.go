package main

import (
	"database/sql"
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

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	groupDb := database.NewGroupStorage(
		db,
		logger,
	)

	userDb := database.NewUserStorage(
		db,
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

	if err = groupDb.AddGroups(groupIDs); err != nil {
		panic(err)
	}

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
