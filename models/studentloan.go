package models

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type StudentLoanList struct {
	List []StudentLoan `bson:"list"`
}

type StudentLoan struct {
	ID       bson.ObjectId `bson:"_id"`
	Fullname string        `bson:"customerName"`
	Amount   int           `bson:"paidAmount"`
}

func (sll *StudentLoanList) FindLoan(u User) StudentLoan {
	for _, e := range sll.List {
		if e.Fullname == u.BankAccountName {
			return e
		}
	}
	return StudentLoan{}
}

func (sl *StudentLoan) CSVAmount() string {
	return utils.SetValueCSV(utils.FormatCommas(fmt.Sprint(sl.Amount)))
}
