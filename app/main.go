package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"

	"github.com/Gklenskiy/vkdigest_bot/app/cmd"
	"github.com/Gklenskiy/vkdigest_bot/app/models"
	"github.com/Gklenskiy/vkdigest_bot/app/proc"
)

// Opts with all cli commands and flags
type Opts struct {
	TelegramCmd cmd.TelegramCommand `command:"bot"`
	ConsoleCmd  cmd.ConsoleCommand  `command:"cmd"`
	ServerCmd   cmd.ServerCommand   `command:"server"`

	Conf string `short:"c" long:"conf" env:"MR_CONF" default:"conf.yml" description:"config file (yml)"`
	Dbg  bool   `long:"dbg" env:"DEBUG" description:"debug mode"`

	DbPort     int    `long:"db_port" env:"DB_PORT" default:"5432" description:"port for database"`
	DbHost     string `long:"db_host" env:"DB_HOST" default:"localhost" description:"host for database"`
	DbUser     string `long:"db_user" env:"DB_USER" default:"postgres" description:"user for database"`
	DbPassword string `long:"db_password" env:"DB_PASSWORD" default:"password" description:"password for database"`
	DbName     string `long:"db_name" env:"DB_NAME" default:"postgres" description:"database name"`
}

var revision = "unknown"

func main() {
	fmt.Printf("vkdigest %s\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.CommonOptionsCommander)

		setupLog(opts.Dbg)

		conf, err := loadConfig(opts.Conf)
		if err != nil {
			log.Fatalf("[ERROR] can't load config %s, %v", opts.Conf, err)
		}

		settings := models.DbSettings{
			Port:     opts.DbPort,
			Host:     opts.DbHost,
			User:     opts.DbUser,
			Password: opts.DbPassword,
			Dbname:   opts.DbName,
		}
		error := models.InitDB(settings)
		if error != nil {
			log.Fatalf("Init Db error: %v", error)
		}

		c.SetCommon(cmd.CommonOpts{
			Revision:   revision,
			Conf:       *conf,
			DbPort:     opts.DbPort,
			DbHost:     opts.DbHost,
			DbUser:     opts.DbUser,
			DbPassword: opts.DbPassword,
			DbName:     opts.DbName,
		})

		err = c.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			log.Printf("[ERROR] failed with %+v", err)
			os.Exit(0)
		} else {
			log.Printf("[ERROR] failed with %+v", err)
			os.Exit(1)
		}
	}
}

func loadConfig(fname string) (res *proc.Conf, err error) {
	log.Printf(fname)
	res = &proc.Conf{}
	data, err := ioutil.ReadFile(fname) // nolint
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, res); err != nil {
		return nil, err
	}

	return res, nil
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
