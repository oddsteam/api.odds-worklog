package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func main() {
	sessionId := os.Getenv("SESSION")
	fmt.Println(sessionId)
	csrf := os.Getenv("CSRF")
	fmt.Println(csrf)
	body, err := getStudentLoans(sessionId, csrf)
	if err != nil {
		fmt.Println(err)
		return
	}
	var loans []models.StudentLoan
	err = json.Unmarshal(body, &loans)

	session := setUpMongo()
	defer session.Close()

	r := income.NewRepository(session)
	for i := range loans {
		loans[i].ID = bson.NewObjectId()
	}
	r.SaveStudentLoans(loans)
}

func getStudentLoans(sessionId string, csrf string) ([]byte, error) {
	url := "https://slfrd.dsl.studentloan.or.th/SLFRD/EmployeeReport/getDataByPage"
	method := "POST"

	y, m := utils.GetYearMonthInBuddistEra(time.Now())
	payload := strings.NewReader(fmt.Sprintf(`deleteFlag=&month=%d&year=%d`, m, y))

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println("making request for student loans fail")
		return nil, err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", sessionId)
	req.Header.Add("X-CSRF-TOKEN", csrf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("request for student loans fail")
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error when parsing student loans response")
		return nil, err
	}
	fmt.Println("==== Body ====")
	fmt.Println(string(body))
	fmt.Println("==============")
	return body, nil
}

func setUpMongo() *mongo.Session {
	c := config.Config()
	// Setup Mongo
	session, err := mongo.NewSession(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return session
}
