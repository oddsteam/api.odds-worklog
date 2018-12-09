package file

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
)

func TestUsecase_getTranscriptFilename(t *testing.T) {
	u := mock.MockUser

	filename := getTranscriptFilename(&u)
	assert.NotEmpty(t, filename)

	path := "files/transcripts"
	filenameExp := fmt.Sprintf("%s/%s_%s_", path, strings.ToLower(u.FirstName), strings.ToLower(u.LastName))
	assert.Contains(t, filename, filenameExp)
	assert.Contains(t, filename, ".pdf")
	assert.Equal(t, len(filenameExp)+16, len(filename))
}
