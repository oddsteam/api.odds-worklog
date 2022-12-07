package income

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStudentLoansQueryWithCurrentMonth(t *testing.T) {
	query := loanQuery()
	assert.Equal(t, "12/2565", query["list.monthYear"])
}
