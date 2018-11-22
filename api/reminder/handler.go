package reminder

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	userRepo := user.NewRepository(session)
	incomeRepo := income.NewRepository(session)
	incomeUsecase := income.NewUsecase(incomeRepo, userRepo)

	r = r.Group("/reminder")
	r.GET("/send", func(c echo.Context) error {
		return send(c, incomeUsecase)
	})
}

func send(c echo.Context, incomeUsecase income.Usecase) error {
	incomeIndividualStatusList, err := incomeUsecase.GetIncomeStatusList("N")
	if err != nil {
		return utils.NewError(c, 500, err)
	}
	incomeCorpStatusList, err := incomeUsecase.GetIncomeStatusList("Y")
	if err != nil {
		return utils.NewError(c, 500, err)
	}
	incomeStatusList := append(incomeIndividualStatusList, incomeCorpStatusList...)
	emails := []string{}
	for _, incomeStatus := range incomeStatusList {
		if incomeStatus.Status == "N" {
			emails = append(emails, incomeStatus.User.Email)
		}
	}
	return c.JSON(http.StatusOK, emails)
}
