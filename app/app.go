package app

import (
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	"gitlab.com/sckacr/calltaxi/config"
	"gitlab.com/sckacr/calltaxi/database"
)

type App struct {
	db           *gorm.DB
	passengerBot *linebot.Client
	operatorBot  *linebot.Client
}

var _app = new(App)

func Init(cfg *config.Config) error {
	var err error

	_app.db, err = database.New(cfg.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	if err = database.Migrate(_app.db); err != nil {
		return errors.Wrap(err, "failed to migrate")
	}

	_app.passengerBot, err = linebot.New(
		cfg.PassengerChannelSecret,
		cfg.PassengerChannelToken,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create passenger bot")
	}

	_app.operatorBot, err = linebot.New(
		cfg.OperatorChannelSecret,
		cfg.OperatorChannelToken,
	)

	if err != nil {
		return errors.Wrap(err, "failed to create operator bot")
	}

	return nil
}

func DB() *gorm.DB {
	return _app.db
}

func PassengerBot() *linebot.Client {
	return _app.passengerBot
}

func OperatorBot() *linebot.Client {
	return _app.operatorBot
}
