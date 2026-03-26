package repositories

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const sapExportFailureColl = "sap_export_failure"

type sapExportFailureRepository struct {
	session *mongo.Session
}

func NewSAPExportFailureRepository(session *mongo.Session) usecases.ForLoggingSAPExportFailure {
	return &sapExportFailureRepository{session: session}
}

// NewSAPExportFailureLister returns the same repository type for read-only listing.
func NewSAPExportFailureLister(session *mongo.Session) usecases.ForListingSAPExportFailures {
	return &sapExportFailureRepository{session: session}
}

func (r *sapExportFailureRepository) List(limit int) ([]*models.SAPExportFailureLog, error) {
	coll := r.session.GetCollection(sapExportFailureColl)
	var logs []*models.SAPExportFailureLog
	err := coll.Find(nil).Sort("-createdAt").Limit(limit).All(&logs)
	return logs, err
}

func (r *sapExportFailureRepository) LogSAPExportFailure(log *models.SAPExportFailureLog) error {
	coll := r.session.GetCollection(sapExportFailureColl)
	log.ID = bson.NewObjectId()
	return coll.Insert(log)
}
