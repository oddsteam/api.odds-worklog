package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
)

type Session struct {
	client *mongo.Client
	dbName string
}

func Setup() *Session {
	c := config.Config()
	session, err := NewSession(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return session
}

type dialInfo struct {
	host     string
	database string
	username string
	password string
}

// redactedMongoDialDescription returns a mongodb:// URL shape for logs only; password is never included.
func redactedMongoDialDescription(info dialInfo) string {
	host := info.host
	if host == "" {
		host = "(no addresses)"
	}
	db := info.database
	if db == "" {
		db = "(no database)"
	}
	if info.username != "" {
		return fmt.Sprintf("mongodb://%s:***@%s/%s", info.username, host, db)
	}
	return fmt.Sprintf("mongodb://%s/%s", host, db)
}

func NewSession(config *models.Config) (*Session, error) {
	info := dialInfo{
		host:     config.MongoDBHost,
		database: config.MongoDBName,
		username: config.Username,
		password: config.Password,
	}

	clientOpts := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s", config.MongoDBHost)).
		SetAuth(options.Credential{
			AuthSource: config.MongoDBName,
			Username:   config.Username,
			Password:   config.Password,
		}).
		SetConnectTimeout(60 * time.Second).
		SetMaxPoolSize(uint64(config.MongoDBConectionPool))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		desc := redactedMongoDialDescription(info)
		return nil, fmt.Errorf("%w (attempted: %s)", err, desc)
	}

	if err := client.Ping(ctx, nil); err != nil {
		desc := redactedMongoDialDescription(info)
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("%w (attempted: %s)", err, desc)
	}

	return &Session{client: client, dbName: config.MongoDBName}, nil
}

func (s *Session) Ctx() context.Context {
	return context.Background()
}

func (s *Session) GetCollection(col string) *mongo.Collection {
	return s.client.Database(s.dbName).Collection(col)
}

func (s *Session) Close() {
	if s.client != nil {
		_ = s.client.Disconnect(context.Background())
	}
}

func (s *Session) DropDatabase(db string) error {
	if s.client != nil {
		return s.client.Database(db).Drop(s.Ctx())
	}
	return nil
}
