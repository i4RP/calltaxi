package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.com/sckacr/calltaxi/app"
	"gitlab.com/sckacr/calltaxi/database"
	"gitlab.com/sckacr/calltaxi/model"
	"gitlab.com/sckacr/calltaxi/service"
)

func OperatorHandlerFunc(c *gin.Context) {
	var (
		bot = service.NewOperator(app.OperatorBot())
		dao = database.NewOperator(app.DB())
	)

	events, err := app.OperatorBot().ParseRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeFollow:
			if _, err := dao.Create(model.NewOperator(event.Source.UserID)); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}
		case linebot.EventTypeUnfollow:
			operator, err := dao.FindByUserID(event.Source.UserID)
			if err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}

			if err := dao.Delete(operator.ID); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}
		case linebot.EventTypePostback:
			if err := bot.ReactToPostback(event); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}
		}
	}
}
