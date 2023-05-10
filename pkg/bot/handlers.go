package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"tg-bot/pkg/constants"
	"tg-bot/pkg/types"
	"tg-bot/pkg/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HandleCommandFunc func(b *Bot, update tgbotapi.Update, ctx context.Context) (*tgbotapi.MessageConfig, error)

func HandleKickBot(b *Bot, update tgbotapi.Update, ctx context.Context) error {
	if update.Message.LeftChatMember != nil && update.Message.LeftChatMember.ID == b.Bot.Self.ID {
		if err := b.Database.RemoveGroup(ctx, update.Message.Chat.ID); err != nil {
			return err
		}
	}
	return nil
}

func HandleBotPermissions(b *Bot, update tgbotapi.Update, ctx context.Context) error {
	var timestamp time.Time
	err := b.Database.Conn.QueryRow(ctx, "SELECT joined_at FROM active_chats WHERE tg_chat = $1", update.Message.Chat.ID).Scan(&timestamp)
	if err != nil {
		return err
	}

	elapsed := time.Since(timestamp)
	if elapsed < time.Hour {
		return nil
	}

	chatConfig := tgbotapi.ChatConfigWithUser{
		ChatID: update.Message.Chat.ID,
		UserID: b.Bot.Self.ID,
	}

	chatMember, err := b.Bot.GetChatMember(chatConfig)

	if err != nil {
		return err
	}

	if !chatMember.IsAdministrator() {
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не администратор, предоставьте права администратора для корректной работы.")
		_, err := b.Bot.Send(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func HandleAddGroup(b *Bot, update tgbotapi.Update, ctx context.Context) error {
	if update.Message.Chat != nil && update.Message.NewChatMembers != nil {
		for _, member := range *update.Message.NewChatMembers {
			if member.UserName == b.Bot.Self.UserName {
				inviterID := update.Message.From.ID
				isEmployee, err := b.Database.IsEmployee(ctx, inviterID)
				if err != nil {
					return err
				}
				if !isEmployee {
					if _, err := b.Bot.LeaveChat(tgbotapi.ChatConfig{
						ChatID: update.Message.Chat.ID,
					}); err != nil {
						return err
					}
				}
				isActive, err := b.Database.ActiveChat(ctx, update.Message.Chat.ID)
				if err != nil {
					return err
				}
				if !isActive {
					if err := b.Database.AddGroup(ctx, update.Message.Chat.ID); err != nil {
						return err
					}
				}
				break
			}
		}
	}
	return nil
}

func HandleDismissed(b *Bot, ctx context.Context, users []int, leadMap *map[int64]types.EmployersInfo, errChan chan error) {
	// fill leads users
	for _, userId := range users {
		leadTgID, name, lastName, err := b.Database.GetBxInfo(ctx, userId)
		if err != nil {
			errChan <- err
			return
		}
		b.Logger.Infof("User: %d, Leader: %d", userId, leadTgID)
		userInfo := types.UserInfo{
			UserId:   userId,
			Name:     name,
			LastName: lastName,
		}

		chats, err := b.Database.GetActiveGroups(ctx)
		if err != nil {
			errChan <- err
			return
		}
		chs := make([]string, 0)

		for _, chatId := range chats {
			info, err := b.Database.GetChatInfo(context.Background(), chatId)
			if err != nil {
				errChan <- err
				return
			}
			memConfig := tgbotapi.ChatConfigWithUser{
				ChatID: chatId,
				UserID: userId,
			}
			member, err := b.Bot.GetChatMember(memConfig)
			if err != nil {
				errChan <- err
				return
			}
			if member.IsMember() {
				chs = append(chs, info.Title)
			}
		}

		if v, ok := (*leadMap)[int64(leadTgID)]; ok {
			b.Logger.Infof("append to: %v", v)
			v.Users = append(v.Users, userInfo)
			(*leadMap)[int64(leadTgID)] = v
		} else if !ok {
			(*leadMap)[int64(leadTgID)] = types.EmployersInfo{
				Done: make(chan struct{}),
				Users: []types.UserInfo{
					{
						UserId:   userId,
						Name:     name,
						LastName: lastName,
						Chats: chs,
					},
				},
			}
		}
	}
	b.Logger.Info(leadMap)

	for leadId, info := range *leadMap {
		b.Logger.Infof("leader: %d, users: %v", leadId, info.Users)

		go HandleLeadMessage(b, ctx, int(leadId), info.Users, info.Done, errChan)
	}
}

func HandleLeadMessage(b *Bot, ctx context.Context, leadId int, users []types.UserInfo, done chan struct{}, errChan chan error) {
	b.Logger.Info("Handle lead msg")
	for _, userInfo := range users {
		

		chsJoined := strings.Join(userInfo.Chats, "\n")

		text := fmt.Sprintf("Доброе утро!\nПохоже %s %s больше не работает у нас, но до сих пор не вышел из корпоративных чатов\n%s\nНужно что-то сделать?", userInfo.Name, userInfo.LastName, chsJoined)
		tgMsg := tgbotapi.NewMessage(int64(leadId), text)
		b.Logger.Infof("Msg for lead: %d", leadId)

		kickData := types.ReplyAction{
			UserId: userInfo.UserId,
			LeadId: leadId,
			Type:   "kick",
		}
		b.Logger.Infof("kick data: %v", kickData)

		stayData := types.ReplyAction{
			UserId: userInfo.UserId,
			LeadId: leadId,
			Type:   "stay",
		}
		b.Logger.Infof("stay data: %v", stayData)

		kickDataStr, err := json.Marshal(kickData)
		if err != nil {
			errChan <- err
			return
		}

		stayDataStr, err := json.Marshal(stayData)
		if err != nil {
			errChan <- err
			return
		}

		tgMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Убрать", string(kickDataStr)),
				tgbotapi.NewInlineKeyboardButtonData("Пусть остаётся", string(stayDataStr)),
			),
		)

		// if len(chs) > 0 {
		if _, err := b.Bot.Send(tgMsg); err != nil {
			errChan <- err
			return
		}

		<-done
		// }
	}
}

func HandleKickUser(b *Bot, chatID int64, userID int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	kickConfig := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: int64(chatID),
			UserID: userID,
		},
	}

	_, err := b.Bot.KickChatMember(kickConfig)
	if err != nil {
		msg := tgbotapi.NewMessage(int64(chatID), constants.ErrorMsgText)
		b.Bot.Send(msg)
		errChan <- fmt.Errorf("Error kicking user %d: %s\n", userID, err.Error())
		return
	} else {
		b.Logger.Infof("User %d removed from the group %d\n", userID, chatID)
	}
}

func HandleCallBack(b *Bot, wg *sync.WaitGroup, update tgbotapi.Update, leadMap *map[int64]types.EmployersInfo, errChan chan error) {
	var data types.ReplyAction

	b.Logger.Infof("Incomming data: %s", update.CallbackQuery.Data)
	if err := json.Unmarshal([]byte(update.CallbackQuery.Data), &data); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, constants.ErrorMsgText)
		b.Bot.Send(msg)
		errChan <- err
		return
	}
	b.Logger.Infof("Decoded incomming data: %v", data)

	if data.LeadId == update.CallbackQuery.From.ID {
		if data.Type == "kick" {
			b.Logger.Info("Kick")
			chatIds, err := b.Database.GetActiveGroups((context.Background()))
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, constants.ErrorMsgText)
				b.Bot.Send(msg)
				errChan <- err
				return
			}

			wg.Add(len(chatIds))
			for _, chatID := range chatIds {
				go func(chatID int64) {
					defer wg.Done()
					kickConfig := tgbotapi.KickChatMemberConfig{
						ChatMemberConfig: tgbotapi.ChatMemberConfig{
							ChatID: int64(chatID),
							UserID: data.UserId,
						},
					}

					_, err := b.Bot.KickChatMember(kickConfig)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, constants.ErrorMsgText)
						b.Logger.Warnf("Error in kick member from chat: %s", err.Error())
						b.Bot.Send(msg)
						return
					} else {
						b.Logger.Infof("User %d removed from the group %d\n", data.UserId, chatID)
					}
				}(chatID)
			}
			wg.Wait()
			newTxt := fmt.Sprintf("%s\n\nУдаляем.", update.CallbackQuery.Message.Text)

			editMsg := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				newTxt,
			)
			editMsg.ReplyMarkup = nil
			b.Bot.Send(editMsg)
		} else if data.Type == "stay" {
			b.Logger.Info("Stay")
			newTxt := fmt.Sprintf("%s\n\nОставляем.", update.CallbackQuery.Message.Text)
			editMsg := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				newTxt,
			)
			editMsg.ReplyMarkup = nil
			b.Bot.Send(editMsg)
		}

		(*leadMap)[int64(data.LeadId)].Done <- struct{}{}
	}
}

// Commands
func HandleStart(b *Bot, update tgbotapi.Update, ctx context.Context) (*tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig

	// Check local DB
	exists, err := b.Database.UserExists(ctx, update.Message.From.ID)
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, constants.ErrorMsgText)
		return nil, err
	}

	signStr := utils.HashWithSalt(strconv.Itoa(update.Message.From.ID), b.salt)

	if !exists {
		link := fmt.Sprintf("%s/bellboy/start.php?chat_id=%d&chat_sign=%s&env=%s", b.Host, int64(update.Message.From.ID), signStr, b.Env)
		authLink := utils.ToLink("ссылке", link)
		txt := fmt.Sprintf("Для продолжения работы авторизуйся в Битрикс24 по %s.", authLink)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, txt)
		msg.ParseMode = "HTML"

		return &msg, nil
	}

	firstName, city, err := b.Database.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, constants.ErrorMsgText)
		return nil, err
	}

	link := fmt.Sprintf("%s/bellboy/start.php?chat_id=%d&chat_sign=%s&env=%s", b.Host, int64(update.Message.From.ID), signStr, b.Env)
	authLink := utils.ToLink("ссылке", link)

	msgText := fmt.Sprintf(constants.HelloMsgText, firstName, city, authLink)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = "HTML"

	return &msg, nil
}

func HandleLinks(b *Bot, update tgbotapi.Update, ctx context.Context) (*tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig

	isAuth, err := b.Database.IsAuth(ctx, update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	if !isAuth {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ты не авторизовался, чтобы пройти авторизацию запусти команду /start и перейди по сгенерированной ссылке!")
		return &msg, nil
	}

	chats, err := b.Database.GetChatLinks(ctx, update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	allLinks := make([]string, len(chats.All))
	for _, record := range chats.All {
		link, err := b.Database.GetInviteLink(ctx, int64(record.ID))
		if err != nil {
			return nil, err
		}

		htmlLink := utils.ToLink(record.Title, link)

		str := fmt.Sprintf("%s - %s", htmlLink, record.Description)
		allLinks = append(allLinks, str)
	}

	cityLinks := make([]string, len(chats.City))
	for _, record := range chats.City {
		link, err := b.Database.GetInviteLink(ctx, int64(record.ID))
		if err != nil {
			return nil, err
		}

		htmlLink := utils.ToLink(record.Title, link)

		str := fmt.Sprintf("%s - %s", htmlLink, record.Description)
		cityLinks = append(cityLinks, str)
	}

	devLinks := make([]string, len(chats.Dev))
	if len(chats.Dev) > 0 {
		for _, record := range chats.Dev {
			link, err := b.Database.GetInviteLink(ctx, int64(record.ID))
			if err != nil {
				return nil, err
			}

			htmlLink := utils.ToLink(record.Title, link)

			str := fmt.Sprintf("%s - %s", htmlLink, record.Description)
			devLinks = append(devLinks, str)
		}
	}

	recommendedLinks := make([]string, len(chats.Recommended))
	for _, record := range chats.Recommended {
		link, err := b.Database.GetInviteLink(ctx, int64(record.ID))
		if err != nil {
			return nil, err
		}

		htmlLink := utils.ToLink(record.Title, link)

		str := fmt.Sprintf("%s - %s", htmlLink, record.Description)
		recommendedLinks = append(recommendedLinks, str)
	}

	joinedAll := strings.Join(allLinks, "\n")
	allResult := fmt.Sprintf("Обязятельные ссылки для вступления:%s\n", joinedAll)

	joinedCity := strings.Join(cityLinks, "\n")
	cityResult := fmt.Sprintf("Ссылки по городу в котором ты живешь:%s\n", joinedCity)

	joinedDev := strings.Join(devLinks, "\n")
	devResult := fmt.Sprintf("Ссылки по твоему отделу:%s\n", joinedDev)

	joinedRec := strings.Join(recommendedLinks, "\n")
	recResult := fmt.Sprintf("Ссылки по интересам:%s\n", joinedRec)

	msgText := fmt.Sprintf("%s\n%s\n%s\n%s", allResult, cityResult, devResult, recResult)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	msg.ParseMode = "HTML"

	return &msg, nil
}
