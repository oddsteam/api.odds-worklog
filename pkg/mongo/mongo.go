package mongo

import (
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
	mgo "gopkg.in/mgo.v2"
)

type Session struct {
	session *mgo.Session
}

func NewSession() (*Session, error) {
	session, err := mgo.Dial(config.MongoDBHost)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetPoolLimit(config.MongoDBConectionPool)
	return &Session{session}, err
}

func (s *Session) Copy() *Session {
	return &Session{s.session.Copy()}
}

func (s *Session) GetCollection(col string) *mgo.Collection {
	return s.session.DB(config.MongoDBName).C(col)
}

func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

func (s *Session) DropDatabase(db string) error {
	if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil
}
