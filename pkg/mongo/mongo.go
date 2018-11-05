package mongo

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	mgo "gopkg.in/mgo.v2"
)

type Session struct {
	MgoSession *mgo.Session
	DBName     string
}

func NewSession(config *models.Config) (*Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{config.MongoDBHost},
		Timeout:  60 * time.Second,
		Database: config.MongoDBName,
		Username: config.Username,
		Password: config.Password,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetPoolLimit(config.MongoDBConectionPool)

	return &Session{session, config.MongoDBName}, err
}

func (s *Session) Copy() *Session {
	return &Session{s.MgoSession.Copy(), s.DBName}
}

func (s *Session) GetCollection(col string) *mgo.Collection {
	return s.MgoSession.DB(s.DBName).C(col)
}

func (s *Session) Close() {
	if s.MgoSession != nil {
		s.MgoSession.Close()
	}
}

func (s *Session) DropDatabase(db string) error {
	if s.MgoSession != nil {
		return s.MgoSession.DB(db).DropDatabase()
	}
	return nil
}
