package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if strings.HasPrefix(update.Message.Command(), "start") {
			text := fmt.Sprintf("Hi %v, I am XurBot and I'll show you Xur's location. Just use /xur_location command.",
				update.Message.From.FirstName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
		if strings.HasPrefix(update.Message.Command(), "item_") {
			hash := hashRegexp(update.Message.Command())
			item, errHash := ParseHashDataOneItem(hash)
			if errHash != nil {
				log.Println(errHash)
			}
			if reflect.ValueOf(item).IsZero() {
				imp := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				imp.Text = "Incorrect command"
				if _, err := bot.Send(imp); err != nil {
					log.Println(err)
				}
			} else {
				imp := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, nil)
				text := fmt.Sprintf("<b>%v</b>\n"+
					"<i>%v</i>\n\n"+
					"%v\n\n<a href=\"https://www.light.gg/db/items/%v\">More info</a>",
					item.Name,
					item.ItemTypeAndTierDisplayName,
					item.FlavorText,
					item.Hash,
				)
				imp.FileID = "https://www.bungie.net" + item.Screenshot
				imp.Caption = text
				imp.ParseMode = "HTML"
				imp.UseExisting = true
				if _, err := bot.Send(imp); err != nil {
					log.Println(err)
				}
			}

		}

		if strings.HasPrefix(update.Message.Command(), "xur_location") {
			xur := getXurData()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "HTML"
			if xur == nil {
				text_msg := fmt.Sprintln("Xur isn't here")
				msg.Text = text_msg
			} else {
				text_msg := []string{fmt.Sprintf("I am on <b>%v</b> in <b>%v</b> and I have something for you:\n",
					xur.LocationName,
					xur.PlaceName),
				}
				xurItems := ParseHashesData()
				text_msg = append(text_msg, xurItems...)
				msg.Text = strings.Join(text_msg[:], "\n")
			}
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func hashRegexp(s string) string {
	re := regexp.MustCompile(`^item_(\d+)$`)
	regexpResult := re.FindString(s)
	if regexpResult == "" {
		return ""
	} else {
		return regexpResult[5:]
	}
}
