package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	"gitlab.com/sckacr/calltaxi/app"
	"gitlab.com/sckacr/calltaxi/database"
	"gitlab.com/sckacr/calltaxi/model"
)

type Operator interface {
	ParseRequest(*gin.Context) ([]*linebot.Event, error)
	PushRequestToAll(*model.Request) error
	ReactToPostback(*linebot.Event) error
}

type OperatorImpl struct {
	*linebot.Client
}

func NewOperator(bot *linebot.Client) Operator {
	return &OperatorImpl{bot}
}

func (bot *OperatorImpl) ParseRequest(c *gin.Context) ([]*linebot.Event, error) {
	return bot.ParseRequest(c)
}

func (bot *OperatorImpl) PushRequestToAll(request *model.Request) error {
	dao := database.NewOperator(app.DB())

	operators, err := dao.FindAll()
	if err != nil {
		return errors.Wrap(err, "failed to find all operators")
	}

	to := make([]string, 0)
	for _, operator := range operators {
		to = append(to, operator.UserID)
	}

	if _, err := bot.Multicast(
		to,
		linebot.NewTextMessage(
			request.Address+"にタクシーを要求しているお客様がいます",
		).WithQuickReplies(linebot.NewQuickReplyItems(
			linebot.NewQuickReplyButton(
				"",
				linebot.NewPostbackAction(
					"タクシーを呼ぶ",
					"call:"+request.ID,
					"タクシーを呼ぶ",
					"",
				),
			),
		)),
	).Do(); err != nil {
		errors.Wrap(err, "failed to multicast")
	}

	return nil
}

func (bot *OperatorImpl) ReactToPostback(event *linebot.Event) error {
	var (
		data         = strings.Split(event.Postback.Data, ":")
		requestDao   = database.NewRequest(app.DB())
		operatorDao  = database.NewOperator(app.DB())
		passengerDao = database.NewPassenger(app.DB())
	)

	switch data[0] {
	case "call":
		request, err := requestDao.FindByID(data[1])
		if err != nil {
			return errors.Wrap(err, "failed to find request")
		}

		if request.Finished || request.OperatorID != "" {
			_, err := bot.PushMessage(
				event.Source.UserID,
				linebot.NewTextMessage(
					"申し訳ありません, 他の方が先に申し込みました",
				),
			).Do()

			return err
		}

		request.OperatorID = event.Source.UserID

		request, err = requestDao.Update(request)
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}

		operators, err := operatorDao.FindAll()
		if err != nil {
			return errors.Wrap(err, "failed to find all operators")
		}

		to := make([]string, 0)
		for _, operator := range operators {
			if operator.UserID != request.OperatorID {
				to = append(to, operator.UserID)
			}
		}

		if _, err := bot.Multicast(
			to,
			linebot.NewTextMessage(
				"締め切りました",
			),
		).Do(); err != nil {
			errors.Wrap(err, "failed to multicast")
		}

		passenger, err := passengerDao.FindByUserID(request.PassengerID)
		if err != nil {
			return errors.Wrap(err, "failed to find passenger")
		}

		if _, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(
				"ご協力ありがとうございます\n"+
					passenger.Name+"様名義で"+
					request.Address+"にタクシーを呼んでください\n"+
					"電話番号は03-5755-2151です\n"+
					"タクシーが到着するまでの時間も聞いてください",
			),
		).Do(); err != nil {
			return errors.Wrap(err, "failed to reply message to operator")
		}

		if _, err := bot.PushMessage(
			request.OperatorID,
			linebot.NewTemplateMessage(
				"時間の確認",
				linebot.NewButtonsTemplate(
					"",
					"時間の確認",
					"タクシーは何分で到着しますか",
					linebot.NewPostbackAction(
						"5分以内",
						"finish:a:"+request.ID,
						"5分以内",
						"",
					),
					linebot.NewPostbackAction(
						"10分以内",
						"finish:b:"+request.ID,
						"10分以内",
						"",
					),
					linebot.NewPostbackAction(
						"15分以内",
						"finish:c:"+request.ID,
						"15分以内",
						"",
					),
					linebot.NewPostbackAction(
						"15分以上",
						"finish:d:"+request.ID,
						"15分以上",
						"",
					),
				),
			).WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewPostbackAction(
						"配車に失敗",
						"error:"+request.ID,
						"配車に失敗",
						"",
					),
				),
			)),
		).Do(); err != nil {
			return errors.Wrap(err, "failed to push message to operator")
		}
	case "finish":
		request, err := requestDao.FindByID(data[2])
		if err != nil {
			return errors.Wrap(err, "failed to find request")
		}

		if request.Finished {
			return errors.Wrap(err, "already finished")
		}

		if _, err := app.PassengerBot().PushMessage(
			request.PassengerID,
			linebot.NewTextMessage("配車の手続きが完了しました. "+when(data[1])),
		).Do(); err != nil {
			return errors.Wrap(err, "failed to push message to passenger")
		}

		request.Finished = true

		if _, err := requestDao.Update(request); err != nil {
			return errors.Wrap(err, "failed to update request")
		}

	case "error":
		request, err := requestDao.FindByID(data[1])
		if err != nil {
			return errors.Wrap(err, "failed to find request")
		}

		if request.Finished {
			return errors.Wrap(err, "already finished")
		}

		if _, err := bot.PushMessage(
			request.OperatorID,
			linebot.NewTextMessage("残念です"),
		).Do(); err != nil {
			return errors.Wrap(err, "failed to push message to operator")
		}

		if _, err := app.PassengerBot().PushMessage(
			request.PassengerID,
			linebot.NewTextMessage("申し訳ありません, 配車に失敗しました"),
		).Do(); err != nil {
			return errors.Wrap(err, "failed to push message to passenger")
		}

		request.Finished = true

		if _, err := requestDao.Update(request); err != nil {
			return errors.Wrap(err, "failed to update request")
		}
	}

	return nil
}

func when(c string) string {
	switch c {
	case "a":
		return "5分以内にタクシーが到着します, お待ちください"
	case "b":
		return "10分以内にタクシーが到着します, お待ちください"
	case "c":
		return "15分以内にタクシーが到着します, お待ちください"
	case "d":
		return "タクシーの到着には15分以上かかります, しばらくお待ちください"
	}

	return ""
}
