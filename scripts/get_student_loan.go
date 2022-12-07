package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	sessionId := os.Getenv("SESSION")
	fmt.Println(sessionId)
	csrf := os.Getenv("CSRF")
	fmt.Println(csrf)
	err := getStudentLoans(sessionId, csrf)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getStudentLoans(sessionId string, csrf string) error {
	url := "https://slfrd.dsl.studentloan.or.th/SLFRD/EmployeeReport/getDataByPage"
	method := "POST"

	payload := strings.NewReader(`deleteFlag=&month=12&year=2565`)

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println("making request for student loans fail")
		return err
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", sessionId)
	req.Header.Add("X-CSRF-TOKEN", csrf)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("request for student loans fail")
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error when parsing student loans response")
		return err
	}
	fmt.Println(string(body))
	return nil
}
