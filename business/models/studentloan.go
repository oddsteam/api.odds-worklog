package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentLoanList struct {
	List []StudentLoan `bson:"list"`
}

type StudentLoan struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
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

func (sll *StudentLoanList) CreateIDForLoans() {
	for i := range sll.List {
		sll.List[i].ID = primitive.NewObjectID()
	}
}

func (sll *StudentLoanList) FindLoan(bankAccountName string) StudentLoan {
	for _, e := range sll.List {
		if strings.Contains(bankAccountName, e.Fullname) {
			return e
		}
	}
	return StudentLoan{}
}

func (sll *StudentLoanList) GetUpdateQuery() bson.M {
	return bson.M{"$set": bson.M{"list": sll.List}}
}

func (sll *StudentLoanList) GetFilterQuery(now time.Time) bson.M {
	return bson.M{"monthYear": GetCurrentMonthInBuddistEra(now)}
}

func (sl *StudentLoan) CSVAmount() string {
	return FormatCommas(fmt.Sprint(sl.Amount))
}
