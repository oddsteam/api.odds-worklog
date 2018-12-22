package invoice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/po"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

type usecase struct {
	invoiceRepo Repository
	poRepo      po.Repository
}

func NewUsecase(invoiceRepo Repository, poRepo po.Repository) Usecase {
	return &usecase{invoiceRepo, poRepo}
}

func (u *usecase) Create(i *models.Invoice) (*models.Invoice, error) {
	// _, err := u.poRepo.GetPO(id)
	// if err != nil {
	// 	return nil, errors.New("PO not found.")
	// }
	return u.invoiceRepo.Create(i)
}

func (u *usecase) Get() ([]*models.Invoice, error) {
	return u.invoiceRepo.Get()
}

func (u *usecase) GetByPO(id string) ([]*models.Invoice, error) {
	// _, err := u.poRepo.GetPO(id)
	// if err != nil {
	// 	return nil, errors.New("PO not found.")
	// }
	return u.invoiceRepo.GetByPO(id)
}

func (u *usecase) GetByID(id string) (*models.Invoice, error) {
	return u.invoiceRepo.GetByID(id)
}

func (u *usecase) NextNo(id string) (string, error) {
	// _, err := u.poRepo.GetPO(id)
	// if err != nil {
	// 	return "", errors.New("PO not found.")
	// }
	invoice, err := u.invoiceRepo.Last(id)
	if err != nil {
		println("error: " + err.Error())
		return newNo("")
	}
	println("OK")
	return newNo(invoice.InvoiceNo)
}

func newNo(last string) (string, error) {
	var no string
	limit := 999
	t := time.Now()
	s := strings.Split(last, "_")
	if s[0] == strconv.Itoa(t.Year()) {
		n, _ := strconv.Atoi(s[1])
		n++
		if n > limit {
			return "", errors.New("Over limit 999 invoices.")
		}
		no = fmt.Sprintf("%s_%03d", s[0], n)
	} else {
		no = fmt.Sprintf("%04d_001", t.Year())
	}
	return no, nil
}

func (u *usecase) Delete(id string) error {
	return u.invoiceRepo.Delete(id)
}

func (u *usecase) Update(i *models.Invoice) (*models.Invoice, error) {
	invoice, err := u.invoiceRepo.GetByID(i.ID.Hex())
	if err != nil {
		return nil, err
	}
	invoice.Amount = i.Amount
	invoice, err = u.invoiceRepo.Update(invoice)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}
