package usecases

type ForUsingExportIncomeFromTimesheet interface {
	ExportIncomeFromTimesheet(role string, monthIndex string) (string, error)
}
