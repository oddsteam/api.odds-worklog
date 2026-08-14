package repositories

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const userColl = "user"
const archivedColl = "archived_user"

type userRepository struct {
	session *mongo.Session
}

func NewUserRepository(session *mongo.Session) usecases.ForManagingUsers {
	return &userRepository{session}
}

func NewTimesheetUserRepository(session *mongo.Session) usecases.ForGettingTimesheetUser {
	return &timesheetUserRepository{&userRepository{session}}
}

// timesheetUserRepository adapts repository's GetByEmail to usecases.ForGettingTimesheetUser's
// contract: callers expect usecases.ErrTimesheetUserNotFound for "no matching user," not the
// raw mongodriver.ErrNoDocuments the shared method returns.
type timesheetUserRepository struct {
	*userRepository
}

func (r *timesheetUserRepository) GetByEmail(email string) (*models.User, error) {
	user, err := r.userRepository.GetByEmail(email)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, usecases.ErrTimesheetUserNotFound
	}
	return user, err
}

func (r *userRepository) Create(u *models.User) (*models.User, error) {
	t := time.Now()
	u.ID = primitive.NewObjectID()
	u.Create = t
	u.LastUpdate = t

	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) Get() ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByRole(role string) ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{"role": role})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, bson.M{"_id": bsonutil.MustObjectIDFromHex(id)}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetBySiteID(id string) ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{"siteId": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByPeakCode(peakCode string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, bson.M{"peakCode": peakCode}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	user.LastUpdate = time.Now()
	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	_, err := coll.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id string) error {
	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	_, err := coll.DeleteOne(ctx, bson.M{"_id": bsonutil.MustObjectIDFromHex(id)})
	return err
}

func (r *userRepository) CreateArchivedUser(user models.User) (*models.ArchivedUser, error) {
	t := time.Now()
	a := models.ArchivedUser{
		ArchivedDate: t,
		User:         user,
	}
	a.ID = primitive.NewObjectID()

	coll := r.session.GetCollection(archivedColl)
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
