package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"tg-bot/pkg/db"
	"tg-bot/pkg/producer"
	"tg-bot/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	config := viper.GetViper()
	salt := config.GetString("env.salt")

	database, err := db.NewDataBase(config)
	if err != nil {
		log.Panicf("DB error: %s", err.Error())
	}

	producer, err := producer.NewProducer(*config)
	if err != nil {
		log.Panic(err)
	}

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", routes.Ping)
		v1.POST("/kick-worker", routes.KickWorker(producer))
	}
	bot := router.Group("/bot")
	{
		bot.GET("/start", routes.Start(salt, config.GetString("tg.name"), database))
	}

	// .Use(middleware.BasicAuth(database))

	// start server
	server := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	// handle OS interrupt signal to gracefully shutdown server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	if err := server.Shutdown(nil); err != nil {
		log.Fatalf("Failed to shutdown server: %s", err)
	}
}
