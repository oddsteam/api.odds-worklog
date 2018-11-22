package reminder

import (
	"encoding/json"
	"io/ioutil"
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

	// Step 1 query user
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

	// Step 2 get user slack
	slackUser, _ := getAllUserSlack()
	slackUserID := []string{}
	for _, email := range emails {
		for _, member := range slackUser.Members {
			if member.Profile.Email == email {
				slackUserID = append(slackUserID, member.ID)
			}
		}
	}

	// Step 3 get channel id
	// Step 4 call post message
	return c.JSON(http.StatusOK, slackUserID)
}

type SlackUsersResponse struct {
	Ok      bool `json:"ok"`
	Members []struct {
		ID       string `json:"id"`
		TeamID   string `json:"team_id"`
		Name     string `json:"name"`
		Deleted  bool   `json:"deleted"`
		Color    string `json:"color"`
		RealName string `json:"real_name"`
		Tz       string `json:"tz"`
		TzLabel  string `json:"tz_label"`
		TzOffset int    `json:"tz_offset"`
		Profile  struct {
			Title                 string `json:"title"`
			Phone                 string `json:"phone"`
			Skype                 string `json:"skype"`
			RealName              string `json:"real_name"`
			RealNameNormalized    string `json:"real_name_normalized"`
			DisplayName           string `json:"display_name"`
			DisplayNameNormalized string `json:"display_name_normalized"`
			StatusText            string `json:"status_text"`
			StatusEmoji           string `json:"status_emoji"`
			StatusExpiration      int    `json:"status_expiration"`
			AvatarHash            string `json:"avatar_hash"`
			Email                 string `json:"email"`
			Image24               string `json:"image_24"`
			Image32               string `json:"image_32"`
			Image48               string `json:"image_48"`
			Image72               string `json:"image_72"`
			Image192              string `json:"image_192"`
			Image512              string `json:"image_512"`
			StatusTextCanonical   string `json:"status_text_canonical"`
			Team                  string `json:"team"`
		} `json:"profile"`
		IsAdmin           bool `json:"is_admin"`
		IsOwner           bool `json:"is_owner"`
		IsPrimaryOwner    bool `json:"is_primary_owner"`
		IsRestricted      bool `json:"is_restricted"`
		IsUltraRestricted bool `json:"is_ultra_restricted"`
		IsBot             bool `json:"is_bot"`
		IsAppUser         bool `json:"is_app_user"`
		Updated           int  `json:"updated"`
		Has2Fa            bool `json:"has_2fa,omitempty"`
	} `json:"members"`
	CacheTs int `json:"cache_ts"`
}

func getAllUserSlack() (SlackUsersResponse, error) {
	token := ""
	url := "https://slack.com/api/users.list?token=" + token
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var data SlackUsersResponse
	_ = json.Unmarshal(body, &data)
	return data, nil
}
