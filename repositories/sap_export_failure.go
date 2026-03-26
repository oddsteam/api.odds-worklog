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

func (r *sapExportFailureRepository) LogSAPExportFailure(log *models.SAPExportFailureLog) error {
	coll := r.session.GetCollection(sapExportFailureColl)
	log.ID = bson.NewObjectId()
	return coll.Insert(log)
}
