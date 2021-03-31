package service

import (
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	"gitlab.com/sckacr/calltaxi/app"
	"gitlab.com/sckacr/calltaxi/database"
	"gitlab.com/sckacr/calltaxi/model"
)

type Passenger interface {
	ParseRequest(*gin.Context) ([]*linebot.Event, error)
	ReactToFollow(*linebot.Event) error
	ReactToMessage(*linebot.Event) error
}

type PassengerImpl struct {
	*linebot.Client
}

func NewPassenger(bot *linebot.Client) Passenger {
	return &PassengerImpl{bot}
}

func (bot *PassengerImpl) ParseRequest(c *gin.Context) ([]*linebot.Event, error) {
	return bot.ParseRequest(c)
}

func (bot *PassengerImpl) ReactToFollow(event *linebot.Event) error {
	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTemplateMessage(
			"登録ありがとうございます",
			linebot.NewButtonsTemplate(
				"",
				"",
				"タクシーくんの友だち登録ありがとうございます\nまずは名義の登録をお願いします",
				linebot.NewMessageAction("名前を登録する", "名義変更"),
			),
		),
	).Do(); err != nil {
		return errors.Wrap(err, "failed to react to follow")
	}

	return nil
}

func (bot *PassengerImpl) ReactToMessage(event *linebot.Event) error {
	switch event.Message.(type) {
	case *linebot.TextMessage:
		return errors.Wrap(
			bot.reactToTextMessage(event),
			"failed to react to text message",
		)
	case *linebot.LocationMessage:
		return errors.Wrap(
			bot.reactToLocationMessage(event),
			"failed to react to location message",
		)
	}

	return nil
}

func (bot *PassengerImpl) reactToTextMessage(event *linebot.Event) error {
	var (
		userID       = event.Source.UserID
		passengerDao = database.NewPassenger(app.DB())
		requestDao   = database.NewRequest(app.DB())

		passenger *model.Passenger
		err       error
	)

	passenger, err = passengerDao.FindByUserID(userID)
	if err != nil {
		return errors.Wrap(err, "failed to find passenger")
	}

	message := event.Message.(*linebot.TextMessage).Text
	switch message {
	case "キャンセル":
		request, err := requestDao.FindLatest(userID)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return errors.Wrap(err, "failed to find request")
		}

		if !request.Finished {
			request.Finished = true

			if request, err = requestDao.Update(request); err != nil {
				return errors.Wrap(err, "failed to update request")
			}

			return errors.Wrap(
				bot.replyTextMessage(event.ReplyToken, "またのご利用をお待ちしております."),
				"failed to react to text message",
			)
		}
	case "名義変更":
		passenger.ChangingName = true

		if passenger, err = passengerDao.Update(passenger); err != nil {
			return errors.Wrap(err, "failed to update passenger")
		}

		return errors.Wrap(
			bot.replyTextMessage(event.ReplyToken, "タクシーを呼ぶ際のお名前をカタカナで教えてください"),
			"failed to react to text message",
		)
	default:
		if passenger.ChangingName {
			for _, c := range message {
				if !unicode.Is(unicode.Katakana, c) {
					return errors.Wrap(
						bot.replyTextMessage(event.ReplyToken, "名前はカタカナでお願いします"),
						"failed to react to text message",
					)
				}
			}

			passenger.Name = message
			passenger.ChangingName = false

			passenger, err = passengerDao.Update(passenger)
			if err != nil {
				return errors.Wrap(err, "failed to update passenger")
			}

			return errors.Wrap(
				bot.replyTextMessage(event.ReplyToken, "\""+message+"\"様で名前を登録しました"),
				"failed to react to text message",
			)
		}
	}

	return errors.Wrap(
		bot.replySelectActionMessage(event.ReplyToken),
		"failed to react to text message",
	)
}

func (bot *PassengerImpl) reactToLocationMessage(event *linebot.Event) error {
	var (
		userID     = event.Source.UserID
		requestDao = database.NewRequest(app.DB())

		err error
	)

	request := model.NewRequest()
	request.PassengerID = userID
	request.Address = event.Message.(*linebot.LocationMessage).Address

	if request, err = requestDao.Create(request); err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTextMessage(
			"近くの配車可能なタクシーを探しています, 1 ~ 3分程度で配車の手配が完了します. (キャンセルの場合は20秒以内にお願いします.)",
		).WithQuickReplies(linebot.NewQuickReplyItems(
			linebot.NewQuickReplyButton("", linebot.NewMessageAction("キャンセル", "キャンセル")),
		)),
	).Do(); err != nil {
		return errors.Wrap(err, "failed to react to location message")
	}

	time.Sleep(20 * time.Second)

	request, err = requestDao.FindByID(request.ID)
	if err != nil {
		return errors.Wrap(err, "failed to find request")
	}

	if !request.Finished {
		if err := bot.pushTextMessage(userID, "配車の手続きに入ります. 少々お待ちください..."); err != nil {
			return errors.Wrap(err, "failed to push text message to passenger")
		}

		if err := NewOperator(app.OperatorBot()).PushRequestToAll(request); err != nil {
			return errors.Wrap(err, "failed to push request to operators")
		}
	}

	return nil
}

func (bot *PassengerImpl) replyTextMessage(token, message string) error {
	_, err := bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do()
	return err
}

func (bot *PassengerImpl) pushTextMessage(to, message string) error {
	_, err := bot.PushMessage(to, linebot.NewTextMessage(message)).Do()
	return err
}

func (bot *PassengerImpl) replySelectActionMessage(token string) error {
	_, err := bot.ReplyMessage(token,
		linebot.NewTemplateMessage(
			"アクション選択",
			linebot.NewButtonsTemplate("", "アクション選択", "何をご所望ですか? お選びください",
				linebot.NewMessageAction("名前を変更する", "名義変更"),
			),
		),
	).Do()

	return err
}
