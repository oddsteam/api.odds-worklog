package models

type Config struct {
	MongoDBHost          string
	MongoDBName          string
	MongoDBConectionPool int
	APIPort              string
	Username             string
	Password             string
	RabbitMQUserName     string
	RabbitMQPassword     string
	RabbitMQHost         string
}
