package models

type ReminderSetting struct {
	Date     string `bson:"date" json:"date"`
	Message  string `bson:"message" json:"message"`
	Time     string `bson:"time" json:"time"`
	Slack    bool   `bson:"slack" json:"slack"`
	Line     bool   `bson:"line" json:"line"`
	Facebook bool   `bson:"facebook" json:"facebook"`
}

type Reminder struct {
	Name    string          `bson:"name" json:"name"`
	Setting ReminderSetting `bson:"setting" json:"setting"`
}
