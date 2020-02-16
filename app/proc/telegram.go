package proc

import (
	"time"

	log "github.com/go-pkgz/lgr"
	tb "gopkg.in/tucnak/telebot.v2"
)

// TelegramBot struct
type TelegramBot struct {
	Bot     *tb.Bot
	Timeout time.Duration
}

// Settings for Telegram bot
type Settings struct {
	Port        string
	PublicURL   string
	Token       string
	WithWebhook bool
	Timeout     time.Duration
}

// NewTelegramBot init vk client
func NewTelegramBot(settings Settings) (*TelegramBot, error) {
	var poller tb.Poller
	if settings.WithWebhook {
		poller = &tb.Webhook{
			Listen:   ":" + settings.Port,
			Endpoint: &tb.WebhookEndpoint{PublicURL: settings.PublicURL},
		}
	} else {
		poller = &tb.LongPoller{Timeout: 3 * time.Second}
	}

	pref := tb.Settings{
		Token:  settings.Token,
		Poller: poller,
	}

	log.Printf("[INFO] create bot with setting, %v", settings)
	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Printf("[ERROR] failed on create telegram bot, %s", err)
		return nil, err
	}

	result := TelegramBot{
		Bot:     bot,
		Timeout: settings.Timeout,
	}

	return &result, err
}

// Start telegram bot
func (tbot *TelegramBot) Start() {
	tbot.Bot.Start()
}

// Use set command endpoint for bot
func (tbot *TelegramBot) Use(endpoint interface{}, handler func(b *tb.Bot, m *tb.Message)) {
	tbot.Bot.Handle(endpoint, func(m *tb.Message) {
		handler(tbot.Bot, m)
	})
}
