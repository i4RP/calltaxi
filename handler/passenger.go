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

func PassengerHandlerFunc(c *gin.Context) {
	var (
		bot = service.NewPassenger(app.PassengerBot())
		dao = database.NewPassenger(app.DB())
	)

	events, err := app.PassengerBot().ParseRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeFollow:
			if _, err := dao.Create(model.NewPassenger(event.Source.UserID)); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}

			if err := bot.ReactToFollow(event); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}
		case linebot.EventTypeUnfollow:
			passenger, err := dao.FindByUserID(event.Source.UserID)
			if err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}

			if err := dao.Delete(passenger.ID); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}
		case linebot.EventTypeMessage:
			if err := bot.ReactToMessage(event); err != nil {
				log.Println("["+event.Type+"]", err)
				c.AbortWithError(http.StatusBadRequest, err)
			}
		}
	}
}
