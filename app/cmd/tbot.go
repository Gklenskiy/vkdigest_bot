package cmd

import (
	"time"

	"github.com/Gklenskiy/vkdigest_bot/app/proc"
	log "github.com/go-pkgz/lgr"
)

// TelegramCommand with params
type TelegramCommand struct {
	Port        string        `long:"port" env:"PORT" default:"5000" description:"port for listen"`
	PublicURL   string        `long:"public_url" env:"PUBLIC_URL" default:"https://api.telegram.org" description:"service url"`
	Token       string        `long:"tg_token" env:"TG_TOKEN" default:"" description:"token for telegram bot"`
	Timeout     time.Duration `long:"timeout" env:"TG_TIMEOUT" default:"" description:"timeout for telegram bot"`
	WithWebhook bool          `long:"hook" env:"WEBHOOK" description:"use webhook"`
	VkAppID     string        `long:"vk_app_id" env:"VK_APP_ID" required:"true" description:"Vk Application ID"`
	AuthURL     string        `long:"auth_url" env:"AUTH_URL" required:"true" description:"Authentication URL"`

	CommonOpts
}

// Execute is the entry point for "server" command, called by flag parser
func (tcmd *TelegramCommand) Execute(args []string) error {
	log.Printf("[INFO] start bot service on port %s", tcmd.Port)

	app, err := tcmd.newTgServiceApp()
	if err != nil {
		log.Printf("[ERROR] failed to setup application, %+v", err)
		return err
	}

	log.Printf("[INFO] set commands")

	api := proc.NewBotCtrl(proc.BotCtrlSettings{
		VkBaseURL:    "https://api.vk.com/method",
		VkAPIVersion: "5.103",
		VkAppID:      tcmd.VkAppID,
		AuthURL:      tcmd.AuthURL,
	})
	app.Use("/ping", api.PingCtrl)
	app.Use("/trend", api.TrendsCtrl)
	app.Use("/start", api.StartCtrl)
	app.Use("/add", api.AddCtrl)
	app.Use("/list", api.ListCtrl)

	log.Printf("[INFO] listen commands")
	app.Start()

	log.Printf("[INFO] bot terminated")
	return nil
}

// newTgServiceApp prepares application and return it with all active parts
// doesn't start anything
func (tcmd *TelegramCommand) newTgServiceApp() (*proc.TelegramBot, error) {
	settings := proc.Settings{
		Port:        tcmd.Port,
		PublicURL:   tcmd.PublicURL,
		Token:       tcmd.Token,
		WithWebhook: tcmd.WithWebhook,
		Timeout:     tcmd.Timeout,
	}
	bot, err := proc.NewTelegramBot(settings)
	if err != nil {
		return nil, err
	}

	return bot, nil
}
