package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx := r.session.Ctx()
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(limit))
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*models.SAPExportFailureLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, err
}

func (r *sapExportFailureRepository) LogSAPExportFailure(log *models.SAPExportFailureLog) error {
	coll := r.session.GetCollection(sapExportFailureColl)
	log.ID = primitive.NewObjectID()
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, log)
	return err
}
