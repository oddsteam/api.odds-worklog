package income

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetStudentLoansQueryWithCurrentMonth(t *testing.T) {
	n := time.Date(2022, time.Month(11), 1, 13, 30, 29, 0, time.UTC)
	query := loanQuery(n)
	assert.Equal(t, "11/2565", query["list.monthYear"])
}
