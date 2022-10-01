package telegram

import (
	"cardWithWords/internal/pkg/storage"
	"fmt"
	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

const textButton = "give me a new card"

type telegramBot struct {
	api           *tba.BotAPI
	wordsQuantity int
	words         storage.Words
}

type Telegram interface {
	ListenAndServeForWords(wg *sync.WaitGroup, done chan struct{}) error
}

func GetAccessToTelegramApi(token string, storage storage.Words, wordsQuantity int) (Telegram, error) {
	var (
		bot = new(telegramBot)
		err error
	)

	bot.wordsQuantity = wordsQuantity
	bot.words = storage

	bot.api, err = tba.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("[GetAccessToTelegramApi]: couldn't create new bot api: %w", err)
	}

	return bot, nil
}

// ListenAndServeForWords generate words by request
func (tb *telegramBot) ListenAndServeForWords(wg *sync.WaitGroup, done chan struct{}) error {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	uc := tba.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	uc.Timeout = 15

	// Start polling Telegram for updates
	updates := tb.api.GetUpdatesChan(uc)

	for {
		select {
		case <-done:
			wg.Done()
			return nil
		case u := <-updates:
			if u.Message == nil {
				continue
			}

			card, err := tb.words.GetCard(tb.wordsQuantity)
			if err != nil {
				log.Printf("[ListenAndServeForWords] couldn't get card for %s: %v\n", u.Message.From.UserName, err)
				card = "try to get card later"
			}

			msg := tba.NewMessage(u.Message.Chat.ID, card)
			msg.ReplyMarkup = tba.NewReplyKeyboard(
				[]tba.KeyboardButton{
					{Text: textButton},
				},
			)

			_, err = tb.api.Send(msg)
			if err != nil {
				log.Printf("[ListenAndServeForWords] couldn't send msg to %s: %v\n", u.Message.From.UserName, err)
			}
		}
	}
}
