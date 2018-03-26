package main

import (
	"fmt"

	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sadlil/gologger"
	"io/ioutil"
	"strings"
	"regexp"
)

var (
	logger gologger.GoLogger
)

func main() {
	s, err := ioutil.ReadFile("buglist")
	if err != nil {
		fmt.Errorf("cannot read source material, make sure you placed it in 'buglist' file\n %s", err)
	}
	cs := string(s)

	regexpC := regexp.MustCompile(`(?:\[.+] |\[.+]: |\[.+]|\[.+]:)(.+)`)
	regexpF := regexp.MustCompile(`(?:{.+} |{.+}|{.+] |{.+])(.+)`)

	chain := NewChain()
	sl := strings.Split(cs, "\n")
	for _, line := range sl {
		if line[0] == '[' {
			r := regexpC.FindStringSubmatch(line)
			line = r[1]
		}
		if line[0] == '{' {
			r := regexpF.FindStringSubmatch(line)
			line = r[1]
		}
		if len(line) > 3 {
			if line[len(line)-1] == ' ' && line[len(line)-2] == ' ' {
				line = line[:len(line)-2]
			} else if line[len(line)-1] == ' ' {
				line = line[:len(line)-1]
			}
			chain.AddSentence(strings.Split(line, " "))
		}
	}


	logger = gologger.GetLogger(gologger.CONSOLE, gologger.ColoredLog)

	token := os.Getenv("TGBOTTOKEN")
	if token == "" || token == "token" {
		panic("you must set TGBOTTOKEN before starting")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	bot.Debug = false
	production := os.Getenv("TGPROD")
	if production != "1" {
		bot.Debug = true
	}

	logger.Info(fmt.Sprintf("Authorized as @%s", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 80

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		logger.Info(fmt.Sprintf("%d (@%s) in %d (%q) executed %q with arguments %q", update.Message.From.ID, update.Message.From.UserName, update.Message.Chat.ID, update.Message.Chat.Title, update.Message.Command(), update.Message.CommandArguments()))

		switch update.Message.Command() {
		case "bug":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, chain.RandomSentence())
			if update.Message.Chat.Type != "private" {
				msg.ReplyToMessageID = update.Message.MessageID
			}
			bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*Markov chain based bug generator.*\n\nWritten by [raccoon](t.me/Rirush), inspired by [melulza](t.me/Melulza), powered by [gophers](golang.org).\n\nType /bug to create a bug.")
			msg.ParseMode = "markdown"
			msg.DisableWebPagePreview = true
			if update.Message.Chat.Type != "private" {
				msg.ReplyToMessageID = update.Message.MessageID
			}
			bot.Send(msg)
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*Markov chain based bug generator.*\n\nWritten by [raccoon](t.me/Rirush), inspired by [melulza](t.me/Melulza), powered by [gophers](golang.org).\n\nType /bug to create a bug.")
			msg.ParseMode = "markdown"
			msg.DisableWebPagePreview = true
			if update.Message.Chat.Type != "private" {
				msg.ReplyToMessageID = update.Message.MessageID
			}
			bot.Send(msg)
		}

	}
}
