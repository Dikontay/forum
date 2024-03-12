package main

import (
	"flag"
	"forum/configs"
	"forum/internal/app"
	"forum/internal/handlers"
	"forum/internal/render"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/pkg/client/sqlite"
	"log"
	"os"
)

func main() {
	log.Println("wait a minute...")

	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	cfg, err := configs.GetConfig(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	db, err := sqlite.OpenDB(cfg.DB.DSN)
	if err != nil {
		log.Println(err)
		return
	}

	repo := repository.NewRepository(db)
	//file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	// info := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	service := service.NewService(repo, logger)

	template, err := render.NewTemplateHTML(cfg.TemplateDir)
	if err != nil {
		log.Println(err)
		return
	}

	handler := handlers.NewHandler(service, template, cfg.GoogleConfig, cfg.GithubConfig)

	err = app.Server(cfg, handler.Routes())

	if err != nil {

		log.Println("Ooopss...\n", err)
		return
	}
}
