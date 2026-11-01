package usecases

type ForUsingDeleteIncome interface {
	DeleteCurrentMonthIncome(userID string) error
}
