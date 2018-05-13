package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nlopes/slack"
	"gopkg.in/telegram-bot-api.v4"
)

type tomlConfig struct {
	Telegram telegramConfig
	Slack    []slackAccount
}

type telegramConfig struct {
	User  int64
	Token string
}

type slackAccount struct {
	Name  string
	Token string
}

func parseConfig(filename string) (tomlConfig, error) {
	var config tomlConfig

	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return config, err
	}

	if config.Telegram.Token == "" || config.Telegram.User == 0 {
		return config, errors.New("Need to specify telegram user and token in config file")
	}

	for _, s := range config.Slack {
		if s.Token == "" || s.Name == "" {
			return config, errors.New("Need to specify slack name and token in config file")
		}
	}
	return config, nil
}

func connectTelegramBotAPI(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect telegram bot: %v", err)
	}
	return bot, nil
}

func logSlackError(workspace string, err error) {
	log.Printf("[%v]: %v", workspace, err)
}

func logSlackMessage(workspace string, message string) {
	log.Printf("[%v]: %v", workspace, message)
}

func connectSlackRTM(slackName, slackToken string, telegramUser int64, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	api := slack.New(slackToken)
	logger := log.New(os.Stdout, fmt.Sprintf("%v: ", slackName), log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	rtm := api.NewRTM()
	go rtm.ManageConnection()
	logSlackMessage(slackName, "Started listening for messages")

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			info := rtm.GetInfo()
			myUserID := info.User.ID
			myUserName := info.User.Name

			direct := strings.HasPrefix(ev.Msg.Channel, "D")
			if !direct && !strings.Contains(ev.Msg.Text, "@"+myUserID) {
				continue
			}

			presence, err := api.GetUserPresence(myUserID)
			if err != nil {
				logSlackError(slackName, fmt.Errorf("failed to get user presence %s", err))
				continue
			}
			if presence.ConnectionCount > 1 {
				logSlackError(slackName, fmt.Errorf("not sending telegram message as user is online. %+v", presence))
				continue
			}

			user, err := api.GetUserInfo(ev.User)
			if err != nil {
				logSlackError(slackName, fmt.Errorf("failed to get user '%v' info %s", ev.User, err))
				continue
			}

			message := fmt.Sprintf("[%v]:: %v sent message '%v'", slackName, user.Name, ev.Text)
			if !direct {
				channel, err := api.GetChannelInfo(ev.Msg.Channel)
				if err != nil {
					logSlackError(slackName, fmt.Errorf("failed to get channel (%v) info %s", ev.Msg.Channel, err))
					continue
				}
				updatedText := strings.Replace(ev.Text, fmt.Sprintf("<@%v>", myUserID), "@"+myUserName, 1)
				message = fmt.Sprintf("[%v]:[#%v]: %v sent message '%v'", slackName, channel.Name, user.Name, updatedText)
			}

			logSlackMessage(slackName, message)
			msg := tgbotapi.NewMessage(telegramUser, message)
			_, err = bot.Send(msg)
			if err != nil {
				logSlackError(slackName, fmt.Errorf("failed to send telegram message: %v", err))
				continue
			}
		case *slack.RTMError:
			logSlackError(slackName, fmt.Errorf("error message received: %s", ev.Error()))
		case *slack.InvalidAuthEvent:
			logSlackError(slackName, fmt.Errorf("invalid credentials - not listening to messages in this workspace"))
			wg.Done()
			return
		}
	}
}

func main() {
	configFlag := flag.String("config", "config.toml", "config file")
	flag.Parse()

	config, err := parseConfig(*configFlag)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := connectTelegramBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to telegram API as %s", bot.Self.UserName)

	var wg sync.WaitGroup
	wg.Add(len(config.Slack))
	for _, s := range config.Slack {
		go connectSlackRTM(s.Name, s.Token, config.Telegram.User, bot, &wg)
	}
	wg.Wait()
}
