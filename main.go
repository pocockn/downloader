package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/pocockn/downloader/config"
	"github.com/pocockn/downloader/handlers"
	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/store"
	"github.com/pocockn/downloader/watcher"
	"github.com/pocockn/downloader/worker"
)

const cfgPath = "config.yaml"

func main() {
	cfg, err := config.New(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(
		middleware.RequestID(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
		middleware.Recover(),
		middleware.Logger(),
	)

	db, err := store.ConnectBolt(cfg.TableName)
	if err != nil {
		log.Fatal(err)
	}

	urlChan := make(chan models.URL)

	pool := worker.NewPool(cfg.Workers, db, urlChan)
	watch := watcher.New(cfg.WatchInterval, db)

	go pool.Run()
	go watch.Process()

	h := handlers.New(db, pool)
	e.POST("store", h.URLStore)
	e.GET("urls", h.URLs)

	log.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
