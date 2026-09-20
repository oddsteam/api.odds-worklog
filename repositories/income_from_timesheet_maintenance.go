package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

// Maintenance helpers for the income_from_timesheet collection. Rows written before the period was
// stored on the record carry no year/month field, and the collection can hold several of them for
// the same user and period — see cmd/resyncincomefromtimesheet, which rebuilds those from the
// timesheet event log and then clears them out.

// legacyIncomeFromTimesheetQuery matches the rows written before the period was stored explicitly.
func legacyIncomeFromTimesheetQuery() bson.M {
	return bson.M{"year": bson.M{"$exists": false}}
}

// EnsureIncomeFromTimesheetPeriodIndex creates the unique period index if it is not there yet. It
// fails while the collection still holds duplicate rows for the same user and period, so the
// resync has to run before it can succeed on an existing database.
func EnsureIncomeFromTimesheetPeriodIndex(session *mongo.Session) error {
	coll := session.GetCollection(incomeFromTimesheetColl)
	_, err := coll.Indexes().CreateOne(session.Ctx(), incomeFromTimesheetPeriodIndex())
	return err
}

// CountIncomeFromTimesheetRows reports how many rows the collection holds in total and how many of
// those are legacy rows without a stored period.
func CountIncomeFromTimesheetRows(session *mongo.Session) (total int64, legacy int64, err error) {
	coll := session.GetCollection(incomeFromTimesheetColl)
	ctx := session.Ctx()
	total, err = coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, 0, err
	}
	legacy, err = coll.CountDocuments(ctx, legacyIncomeFromTimesheetQuery())
	if err != nil {
		return 0, 0, err
	}
	return total, legacy, nil
}

// ListIncomeFromTimesheetRows reads the whole collection. It is only for the migration, which has
// to group every row by period to find the duplicates; nothing on the serving path reads this way.
func ListIncomeFromTimesheetRows(session *mongo.Session) ([]*models.IncomeFromTimesheet, error) {
	coll := session.GetCollection(incomeFromTimesheetColl)
	ctx := session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	rows := make([]*models.IncomeFromTimesheet, 0)
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// SetIncomeFromTimesheetPeriod stamps a recovered period onto one legacy row, leaving every other
// field — the payroll figures and any manual edit among them — untouched.
func SetIncomeFromTimesheetPeriod(session *mongo.Session, id primitive.ObjectID, year int, month time.Month, submitDate time.Time) error {
	coll := session.GetCollection(incomeFromTimesheetColl)
	_, err := coll.UpdateOne(session.Ctx(), bson.M{"_id": id}, bson.M{"$set": bson.M{
		"year":       year,
		"month":      int(month),
		"submitDate": submitDate,
	}})
	return err
}

// DeleteIncomeFromTimesheetRow removes one row by id. The migration uses it to drop the duplicates
// of a period once the row to keep has been chosen.
func DeleteIncomeFromTimesheetRow(session *mongo.Session, id primitive.ObjectID) error {
	coll := session.GetCollection(incomeFromTimesheetColl)
	_, err := coll.DeleteOne(session.Ctx(), bson.M{"_id": id})
	return err
}
