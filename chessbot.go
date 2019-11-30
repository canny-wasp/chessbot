package main

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

const apiToken = "insert token here"

const letterOffset = 96
const numberOffset = 48

const knightGame = "knight"

func main() {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	gameMode := ""
	var knightRound KnightGameResult

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			helpMsg := "type /knight for game or /stopgame"

			switch update.Message.Command() {
			case "start":
				msg.Text = helpMsg
			case "help":
				msg.Text = helpMsg
			case knightGame:
				gameMode = knightGame
				knightRound = playKnightGame()
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Enter all the coordinates for the knight's move from the given coordinate"))
				msg.Text = knightRound.question
			case "stopgame":
				msg.Text = helpMsg
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		} else {
			if gameMode == knightGame {
				answer := parseAnswer(update.Message.Text)

				if reflect.DeepEqual(answer, knightRound.answer) {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Right!"))
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong!"))
				}

				knightRound = playKnightGame()
				questionMsg := tgbotapi.NewMessage(update.Message.Chat.ID, knightRound.question)
				bot.Send(questionMsg)
			}
		}
	}
}

type KnightGameResult struct {
	question string
	answer   map[string]bool
}

func playKnightGame() KnightGameResult {
	q := generateQuestion()
	a := calculateAnswer(q)

	return KnightGameResult{
		question: q,
		answer:   a,
	}
}

func generateQuestion() string {
	var buffer bytes.Buffer

	letter := string(rand.Intn(8) + 97)
	number := strconv.Itoa(rand.Intn(8) + 1)

	buffer.WriteString(letter)
	buffer.WriteString(number)

	return buffer.String()
}

func calculateAnswer(q string) map[string]bool {
	runesWitoutOffset := []rune(q)
	runes := []rune{
		runesWitoutOffset[0] - letterOffset,
		runesWitoutOffset[1] - numberOffset,
	}

	all := []coordinate{
		{x: int(runes[0] - 1), y: int(runes[1] - 2)},
		{x: int(runes[0] - 1), y: int(runes[1] + 2)},
		{x: int(runes[0] + 1), y: int(runes[1] - 2)},
		{x: int(runes[0] + 1), y: int(runes[1] + 2)},
		{x: int(runes[0] - 2), y: int(runes[1] - 1)},
		{x: int(runes[0] - 2), y: int(runes[1] + 1)},
		{x: int(runes[0] + 2), y: int(runes[1] - 1)},
		{x: int(runes[0] + 2), y: int(runes[1] + 1)},
	}

	m := make(map[string]bool)

	for i := 0; i < len(all); i++ {
		if (all[i].x > 0 && all[i].x < 9) && (all[i].y > 0 && all[i].y < 9) {
			m[coordinateName(all[i])] = true
		}
	}

	return m
}

type coordinate struct {
	x int
	y int
}

func coordinateName(c coordinate) string {
	var buffer bytes.Buffer

	buffer.WriteString(string(c.x + letterOffset))
	buffer.WriteString(string(c.y + numberOffset))

	return buffer.String()
}

func parseAnswer(answer string) map[string]bool {
	answerArray := strings.Split(answer, " ")
	answerMap := make(map[string]bool)
	for i := 0; i < len(answerArray); i++ {
		answerMap[answerArray[i]] = true
	}

	return answerMap
}
