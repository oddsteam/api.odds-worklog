package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type StudentLoanList struct {
	List []StudentLoan `bson:"list"`
}

type StudentLoan struct {
	ID        bson.ObjectId `bson:"_id" json:"id,omitempty"`
	Fullname  string        `bson:"customerName" json:"customerName"`
	Amount    int           `bson:"paidAmount" json:"paidAmount"`
	MonthYear string        `bson:"monthYear" json:"monthYear"`
}

func CreateStudentLoanList(studentLoanResponse []byte) (StudentLoanList, error) {
	var loans []StudentLoan
	err := json.Unmarshal(studentLoanResponse, &loans)
	loanlist := StudentLoanList{List: loans}
	return loanlist, err
}

func (sll *StudentLoanList) FindLoan(u User) StudentLoan {
	for _, e := range sll.List {
		if e.Fullname == u.BankAccountName {
			return e
		}
	}
	return StudentLoan{}
}

func (sll *StudentLoanList) GetUpdateQuery() bson.M {
	return bson.M{"$set": bson.M{"list": sll.List}}
}

func (sll *StudentLoanList) GetFilterQuery(now time.Time) bson.M {
	return bson.M{"monthYear": utils.GetCurrentMonthInBuddistEra(now)}
}

func (sl *StudentLoan) CSVAmount() string {
	return utils.SetValueCSV(utils.FormatCommas(fmt.Sprint(sl.Amount)))
}
