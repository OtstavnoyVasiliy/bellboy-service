package bot

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"tg-bot/pkg/db"
	"tg-bot/pkg/types"

	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

type CommandInfo struct {
	Handler    HandleCommandFunc
	Descripton string
}

type Bot struct {
	updateConfig     tgbotapi.UpdateConfig
	Bot              tgbotapi.BotAPI
	Logger           *logrus.Logger
	Database         *db.DataBase
	Host             string
	Env              string
	salt             string
	CommandsHandlers map[string]CommandInfo
}

func NewBot(config *viper.Viper, logger *logrus.Logger, database *db.DataBase) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(config.GetString("tg.token"))
	if err != nil {
		return nil, err
	}

	bot.Debug = config.GetBool("env.debugMode")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &Bot{
		updateConfig:     u,
		Bot:              *bot,
		Logger:           logger,
		Database:         database,
		Host:             "",
		Env:              config.GetString("env.type"),
		salt:             config.GetString("env.salt"),
		CommandsHandlers: make(map[string]CommandInfo),
	}, nil
}

func (b *Bot) HandleCommand(command, descripton string, handleFunc HandleCommandFunc) {
	b.CommandsHandlers[command] = CommandInfo{
		Handler:    handleFunc,
		Descripton: descripton,
	}
}

func (b *Bot) RunBot(userIdMsg <-chan int, notActiveUserIdMsg <-chan []int, errChan chan error) {
	var wg sync.WaitGroup
	leadMap := make(map[int64]types.EmployersInfo)

	updatesChan, err := b.Bot.GetUpdatesChan(b.updateConfig)
	if err != nil {
		errChan <- err
		return
	}

	for {
		select {
		case update := <-updatesChan:
			if update.Message != nil {
				switch update.Message.Chat.Type {
				case "supergroup", "group":
					if err := HandleKickBot(b, update, context.Background()); err != nil {
						errChan <- err
						return
					}

					if err := HandleAddGroup(b, update, context.Background()); err != nil {
						errChan <- err
						return
					}

					if err := HandleBotPermissions(b, update, context.Background()); err != nil {
						errChan <- err
						return
					}

				case "private":
					var msg *tgbotapi.MessageConfig

					if update.Message != nil {
						command := update.Message.Command()

						if command == "help" {
							commands := make([]string, 0, len(b.CommandsHandlers))
							for cmd := range b.CommandsHandlers {
								descripton := b.CommandsHandlers[cmd].Descripton
								commands = append(commands, fmt.Sprintf("/%s - %s", cmd, descripton))
							}

							concatCommands := strings.Join(commands, "\n")
							msgConf := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вот список всех команд и их описание:\n%s", concatCommands))
							msg = &msgConf

							b.Bot.Send(msg)
						} else {
							info, ok := b.CommandsHandlers[command]
							if !ok {
								msgConf := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда, воспользуйся командой /help, для получения всех доступных команд.")
								msg = &msgConf
							} else if msg, err = info.Handler(b, update, context.Background()); err != nil {
								errChan <- err
								return
							}

							b.Bot.Send(msg)
						}
					}
				}
			} else if update.CallbackQuery != nil {
				go HandleCallBack(b, &wg, update, &leadMap, errChan)
			}

		case msg := <-userIdMsg:
			chatIds, err := b.Database.GetActiveGroups((context.Background()))
			if err != nil {
				errChan <- err
				return
			}

			wg.Add(len(chatIds))
			for _, chatID := range chatIds {
				go HandleKickUser(b, chatID, msg, &wg, errChan)
			}
			wg.Wait()

		case msg := <-notActiveUserIdMsg:
			go HandleDismissed(b, context.Background(), msg, &leadMap, errChan)
		}
	}
}
