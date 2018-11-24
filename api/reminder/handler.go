package reminder

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/setting"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/slack"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

// const TOKEN = "xoxb-484294901968-485201164352-IC904vZ6Bxwx2xkI2qzWgy5J" // Reminder workspace
const TOKEN = "xoxb-293071900534-486896062132-2RMbUSdX6DqoOKsVMCSXQoiM" // Odds workspace

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	userRepo := user.NewRepository(session)
	incomeRepo := income.NewRepository(session)
	settingRepo := setting.NewRepository(session)
	incomeUsecase := income.NewUsecase(incomeRepo, userRepo)

	r = r.Group("/reminder")
	r.GET("/send", func(c echo.Context) error {
		return send(c, incomeUsecase, settingRepo)
	})
}

// Send Reminder godoc
// @Summary Send Reminder
// @Description Send Notification Reminder
// @Tags reminder
// @Produce  json
// @Success 200 {object} string
// @Failure 500 {object} utils.HTTPError
// @Router /reminder/send [get]
func send(c echo.Context, incomeUsecase income.Usecase, setting setting.Repository) error {
	isDev := true
	var emails []string
	if isDev {
		emails = []string{
			"tong@odds.team",
			"saharat@odds.team",
			"thanundorn@odds.team",
			"santi@odds.team",
		}
	} else {
		user, err := listEmailUserIncomeStatusIsNo(incomeUsecase)
		if err != nil {
			return utils.NewError(c, 500, err)
		}
		emails = user
	}

	s, err := setting.GetReminder()
	if err != nil {
		return utils.NewError(c, 500, err)
	}

	err = sendNotification(emails, s.Setting.Message)
	if err != nil {
		return utils.NewError(c, 500, err)
	}

	return c.JSON(http.StatusOK, true)
}

func listEmailUserIncomeStatusIsNo(incomeUsecase income.Usecase) ([]string, error) {
	emails := []string{}
	incomeIndividualStatusList, err := incomeUsecase.GetIncomeStatusList("N")
	if err != nil {
		return nil, err
	}
	incomeCorpStatusList, err := incomeUsecase.GetIncomeStatusList("Y")
	if err != nil {
		return nil, err
	}
	incomeStatusList := append(incomeIndividualStatusList, incomeCorpStatusList...)
	for _, incomeStatus := range incomeStatusList {
		if incomeStatus.Status == "N" {
			emails = append(emails, incomeStatus.User.Email)
		}
	}
	return emails, nil
}

func sendNotification(emails []string, message string) error {
	client, err := slack.NewClient(TOKEN)
	slackUsers, err := client.GetUsersList()
	if err != nil {
		return err
	}
	for _, email := range emails {
		for _, member := range slackUsers.Members {
			if member.Profile.Email == email {
				im, err := client.OpenIMChannel(member.ID)
				if err != nil {
					return err
				}
				channelID := im.Channel.ID
				// fmt.Println(channelID)
				client.PostMessage(channelID, message)
			}
		}
	}
	return nil
}
