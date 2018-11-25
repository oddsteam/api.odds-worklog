package reminder

import "gitlab.odds.team/worklog/api.odds-worklog/models"

type Repository interface {
	SaveReminder(reminder *models.Reminder) (*models.Reminder, error)
	GetReminder() (*models.Reminder, error)
}
