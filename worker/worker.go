package worker

import (
	"log"

	"github.com/elgs/cron"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/slack"
)

var ccron = cron.New()
var jobs = make(chan int, 1)

func worker(jobs chan int) {
	log.Println("Started cronjobs")
	ccron.Start()
}

func StartWorker(reminder *models.Reminder) {
	log.Println("starting worker")
	select {
	case checkJobs := <-jobs:
		ccron.RemoveFunc(checkJobs)
		j1, _ := ccron.AddFunc("0 59 23 "+reminder.Setting.Date+" * *", func() { send(reminder) })
		select {
		case jobs <- j1:
			log.Println("have jobs", j1)
		default:
			log.Println("no message sent")
		}
		log.Println("received job", checkJobs)
	default:
		j1, _ := ccron.AddFunc("0 59 23 "+reminder.Setting.Date+" * *", func() { send(reminder) })
		select {
		case jobs <- j1:
			log.Println("have jobs", j1)
		default:
			log.Println("no message sent")
		}
		go worker(jobs)
		log.Println("Worker started")
		log.Println("no job received")
	}
}

func send(reminder *models.Reminder) {
	log.Println("start send message")
	var token string
	var emails []string
	// For production
	// token = "xoxb-293071900534-486896062132-2RMbUSdX6DqoOKsVMCSXQoiM" // Odds workspace
	// user, err := ListEmailUserIncomeStatusIsNo(incomeUsecase)
	// if err != nil {
	// 	return utils.NewError(c, 500, err)
	// }
	// emails = user

	// For development
	token = "xoxb-484294901968-485201164352-IC904vZ6Bxwx2xkI2qzWgy5J" // Reminder workspace
	emails = []string{
		"tong@odds.team",
		"saharat@odds.team",
		"thanundorn@odds.team",
		"santi@odds.team",
		"work.alongkorn@gmail.com",
		"p.watchara@gmail.com",
	}
	//
	err := sendNotification(token, emails, reminder.Setting.Message)
	if err != nil {
		log.Println(err)
	}
}

func sendNotification(token string, emails []string, message string) error {
	client := slack.Client{
		Token: token,
	}
	slackUsers, err := client.GetUserList()
	if err != nil {
		return err
	}
	for _, email := range emails {
		for _, member := range slackUsers.Members {
			if member.Profile.Email == email {
				im, err := client.OpenIMChannel(member.ID)
				if err != nil {
					return err
				}
				channelID := im.Channel.ID
				client.PostMessage(channelID, message)
			}
		}
	}
	return nil
}
