package login

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user/mocks"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestVerifyAudience(t *testing.T) {
	t.Run("956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com is currect Audience", func(t *testing.T) {
		assert.Equal(t, clientID, "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com")
	})

	t.Run("956316396976-cnrmemp4r4coc62oqmn9uin7iq3o3eev.apps.googleusercontent.com2 is wrong Audience", func(t *testing.T) {
		assert.NotEqual(t, clientID, "956316396976-cnrmemp4r4coc62oqmn9uin7iq3o3eev.apps.googleusercontent.com2")
	})
}

func TestHandleToken(t *testing.T) {
	t.Run("when user is first login got FirstLogin = 'Y'", func(t *testing.T) {
		user := new(models.User)
		token, err := handleToken(user)

		assert.NoError(t, err)
		assert.Equal(t, "Y", token.FirstLogin)
	})

	t.Run("when user is't first login got FirstLogin = 'N'", func(t *testing.T) {
		token, err := handleToken(&mocks.MockUser)

		assert.NoError(t, err)
		assert.Equal(t, "N", token.FirstLogin)
	})
}

func TestGenToken(t *testing.T) {
	t.Run("generate token success", func(t *testing.T) {
		token, err := genToken(&mocks.MockUser)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
