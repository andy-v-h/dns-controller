package answers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go.hollow.sh/dnscontroller/internal/models"
)

func qmAnswerDetailAnswerID(id uuid.UUID) qm.QueryMod {
	return qm.Where("answer_id=?", id.String())
}

// GetDetails returns the Detail for an Answer if it's set
func (a *Answer) GetDetails(ctx context.Context, db *sqlx.DB) error {
	if !a.HasDetails {
		return ErrorNoAnswerDetail
	}

	qm := qmAnswerDetailAnswerID(a.UUID)

	dbModels, err := models.AnswerDetails(qm).All(ctx, db)
	if err != nil {
		return err
	}

	for _, dbModel := range dbModels {
		ad := &Detail{}
		if err := ad.FromDBModel(dbModel); err != nil {
			return err
		}

		a.Details = append(a.Details, ad)
	}

	return nil
}

// Delete remobes a detail from the DB
func (d *Detail) Delete(ctx context.Context, db *sqlx.DB) error {
	err := d.validate()
	if err != nil {
		return err
	}

	err = d.Find(ctx, db)
	if err != nil {
		return err
	}

	dbDetail, err := d.ToDBModel()
	if err != nil {
		return err
	}

	_, err = dbDetail.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

// FindOrCreate is the upsert function
func (d *Detail) FindOrCreate(ctx context.Context, db *sqlx.DB) error {
	err := d.Find(ctx, db)
	if errors.Is(err, sql.ErrNoRows) {
		return d.Create(ctx, db)
	} else if err != nil {
		return err
	}

	return nil
}

// Find looks the detail up by answer_id
func (d *Detail) Find(ctx context.Context, db *sqlx.DB) error {
	if err := d.validate(); err != nil {
		return err
	}

	qm := qmAnswerDetailAnswerID(d.AnswerUUID)

	dbDetail, err := models.AnswerDetails(qm).One(ctx, db)
	if err != nil {
		return err
	}

	return d.FromDBModel(dbDetail)
}

// Create a answer detail
func (d *Detail) Create(ctx context.Context, db *sqlx.DB) error {
	dbDetail, err := d.ToDBModel()
	if err != nil {
		return err
	}

	if err := dbDetail.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	// Set the values back
	return d.FromDBModel(dbDetail)
}

// FromDBModel converts a db type to an api type
func (d *Detail) FromDBModel(dbT *models.AnswerDetail) error {
	// Not null values
	d.CreatedAt = dbT.CreatedAt
	d.UpdatedAt = dbT.UpdatedAt

	var err error

	d.UUID, err = uuid.Parse(dbT.ID)
	if err != nil {
		return err
	}

	d.AnswerUUID, err = uuid.Parse(dbT.AnswerID)
	if err != nil {
		return err
	}

	// Null values

	d.Port = dbT.Port.Int64
	d.Priority = dbT.Priority.Int64
	d.Protocol = dbT.Protocol.String
	d.Weight = dbT.Weight.Int64

	return d.validate()
}

// ToDBModel converts the api type to db type
func (d *Detail) ToDBModel() (*models.AnswerDetail, error) {
	if err := d.validate(); err != nil {
		return nil, err
	}

	dbModel := &models.AnswerDetail{
		AnswerID:  d.UUID.String(),
		Port:      null.NewInt64(int64(d.Port), true),
		Priority:  null.NewInt64(int64(d.Priority), true),
		Weight:    null.NewInt64(int64(d.Weight), true),
		Protocol:  null.StringFrom(d.Protocol),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}

	if d.UUID.String() != uuid.Nil.String() {
		dbModel.ID = d.UUID.String()
	}

	return dbModel, nil
}

func (d *Detail) validate() error {
	if d.UUID.String() == uuid.Nil.String() {
		return ErrorNoAnswerDetailID
	}

	if d.AnswerUUID.String() == uuid.Nil.String() {
		return ErrorNoAnswerDetailAnswerID
	}

	return nil
}
