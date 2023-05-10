package main


import (
	"os"
	"os/signal"
	"tg-bot/pkg/bot"
	"tg-bot/pkg/consumer"
	"tg-bot/pkg/db"
	"tg-bot/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	errChan := make(chan error)
	userIdChan := make(chan int)
	notActiveUserIdChan := make(chan []int)

	defer close(errChan)
	defer close(userIdChan)
	defer close(notActiveUserIdChan)

	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		logger.WithFields(logrus.Fields{
			"type": "config-creation",
		}).Fatal(err)
	}

	config := viper.GetViper()
	database, err := db.NewDataBase(config)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"type": "database-creation",
		}).Fatal(err)
	}
	logger.Info("Successful connect to DB")

	tgbot, err := bot.NewBot(config, logger, database)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"type":     "bot-creation",
			"bot-type": "group-bot",
		}).Fatal(err)
	}
	logger.Info("Successful bot creation")

	cons, err := consumer.NewConsumer(*config, logger)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"type": "consumer-creation",
		}).Fatal(err)
	}
	logger.Info("Successful connect to kafka")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	tgbot.HandleCommand("start", "команда для старта бота и получения ссылки на аунтификацию в системе", bot.HandleStart)
	tgbot.HandleCommand("groups", "все нужные и ненужные тебе ссылки на группы", bot.HandleLinks)

	go cons.RunListener(userIdChan, errChan, signalChan)
	go tgbot.RunBot(userIdChan, notActiveUserIdChan, errChan)
	go utils.CheckEmployee(notActiveUserIdChan, errChan, config.GetString("env.bx-api"))

	select {
	case err := <-errChan:
		logger.WithFields(logrus.Fields{
			"type": "inprogress-error",
		}).Panic(err)
	case <-signalChan:
		logger.Infoln("Shutting down...")
	}

	close(errChan)
}
