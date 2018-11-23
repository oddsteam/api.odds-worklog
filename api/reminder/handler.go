package reminder

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/slack"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

const TOKEN = "xoxb-484294901968-485201164352-IC904vZ6Bxwx2xkI2qzWgy5J"

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	userRepo := user.NewRepository(session)
	incomeRepo := income.NewRepository(session)
	incomeUsecase := income.NewUsecase(incomeRepo, userRepo)

	r = r.Group("/reminder")
	r.GET("/send", func(c echo.Context) error {
		return send(c, incomeUsecase)
	})
}

// Send Reminder godoc
// @Summary Send Reminder
// @Description Send Notification Reminder
// @Tags reminder
// @Produce  json
// @Success 200 {array} string
// @Failure 500 {object} utils.HTTPError
// @Router /reminder/send [get]
func send(c echo.Context, incomeUsecase income.Usecase) error {
	client := slack.NewClient(TOKEN)
	// incomeIndividualStatusList, err := incomeUsecase.GetIncomeStatusList("N")
	// if err != nil {
	// 	return utils.NewError(c, 500, err)
	// }
	// incomeCorpStatusList, err := incomeUsecase.GetIncomeStatusList("Y")
	// if err != nil {
	// 	return utils.NewError(c, 500, err)
	// }
	// incomeStatusList := append(incomeIndividualStatusList, incomeCorpStatusList...)
	emails := []string{}
	// for _, incomeStatus := range incomeStatusList {
	// 	if incomeStatus.Status == "N" {
	// 		emails = append(emails, incomeStatus.User.Email)
	// 	}
	// }
	emails = []string{
		"tong@odds.team",
		"work.alongkorn@gmail.com",
		"saharat@odds.team",
		"thanundorn@odds.team",
		"p.watchara@gmail.com",
		"santi@odds.team",
	}
	sendEmails := []string{}
	slackUsers, err := client.GetUsersList()
	if err != nil {
		return utils.NewError(c, 500, err)
	}
	for _, email := range emails {
		for _, member := range slackUsers.Members {
			if member.Profile.Email == email {
				im, err := client.OpenIMChannel(member.ID)
				if err != nil {
					return utils.NewError(c, 500, err)
				}
				channelID := im.Channel.ID
				client.PostMessage(channelID, "กรุณาเข้าไปกรอกเงินเดือนด้วยครับ")
				sendEmails = append(sendEmails, email)
			}
		}
	}
	return c.JSON(http.StatusOK, sendEmails)
}
