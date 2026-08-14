package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
)

// Usecase is the HTTP-facing alias of the manage-users driving port.
type Usecase = usecases.ForUsingUsers
