// Package records wraps the CRUD operations for a models.Record
package records

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.hollow.sh/dnscontroller/internal/models"
	"go.hollow.sh/dnscontroller/pkg/api/v1/answers"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// SetLogger initializes the package logger
func SetLogger(l *zap.SugaredLogger) { logger = l }

func qmRecordNameAndType(rname, rtype string) (qm.QueryMod, error) {
	if err := isSupportedRecordType(rtype); err != nil {
		return nil, err
	}

	mods := []qm.QueryMod{
		qm.Where("record=?", rname),
		qm.Where("record_type=?", rtype),
	}

	return qm.Expr(mods...), nil
}

func qmAnswerRecordID(id uuid.UUID) qm.QueryMod { return qm.Where("record_id=?", id.String()) }

// GetAnswers fetches the ansers for a record
func (r *Record) GetAnswers(ctx context.Context, db *sqlx.DB) error {
	qm := qmAnswerRecordID(r.UUID)

	dbAnswers, err := models.Answers(qm).All(ctx, db)
	if err != nil {
		return err
	}

	// If there are any answers for a record, fetch them
	numAnswers, err := models.Answers(qm).Count(ctx, db)
	if err != nil {
		return err
	}

	if numAnswers > 0 {
		logger.Debugw("answers found for record, fetching", "record", r, "count", numAnswers, "db-answers", dbAnswers)

		for i := range dbAnswers {
			a := &answers.Answer{}
			if err := a.FromDBModel(ctx, db, dbAnswers[i]); err != nil {
				return err
			}

			r.Answers = append(r.Answers, a)
		}
	} else {
		r.Answers = []*answers.Answer{}
	}

	return nil
}

// Delete removes a record from the DB
func (r *Record) Delete(ctx context.Context, db *sqlx.DB) error {
	err := r.validate()
	if err != nil {
		return err
	}

	err = r.Find(ctx, db)
	if err != nil {
		return err
	}

	dbRecord, err := r.ToDBModel()
	if err != nil {
		return err
	}

	_, err = dbRecord.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

// CreateOrFind is the upsert function
func (r *Record) CreateOrFind(ctx context.Context, db *sqlx.DB) error {
	err := r.Create(ctx, db)
	if err != nil {
		return r.Find(ctx, db)
	}

	return nil
}

// Find looks the record up by name,type
func (r *Record) Find(ctx context.Context, db *sqlx.DB) error {
	if err := r.validate(); err != nil {
		return err
	}

	qm, err := qmRecordNameAndType(r.Name, r.Type)
	if err != nil {
		return err
	}

	dbRecord, err := models.Records(qm).One(ctx, db)
	if err != nil {
		return err
	}

	return r.FromDBModel(ctx, db, dbRecord)
}

// Create inserts a record
func (r *Record) Create(ctx context.Context, db *sqlx.DB) error {
	dbRecord, err := r.ToDBModel()
	if err != nil {
		return err
	}

	if err := dbRecord.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return r.FromDBModel(ctx, db, dbRecord)
}

// FromDBModel converts a db type to an api type
func (r *Record) FromDBModel(ctx context.Context, db *sqlx.DB, dbT *models.Record) error {
	r.CreatedAt = dbT.CreatedAt
	r.UpdatedAt = dbT.UpdatedAt
	r.Name = dbT.Record
	r.Type = dbT.RecordType

	var err error

	r.UUID, err = uuid.Parse(dbT.ID)
	if err != nil {
		return err
	}

	if err := r.validate(); err != nil {
		return err
	}

	logger.Debugw("db record conterted to api model", "record", dbT, "api-model", r)

	return r.GetAnswers(ctx, db)
}

// ToDBModel converts the api type to db type
func (r *Record) ToDBModel() (*models.Record, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	dbModel := &models.Record{
		Record:     strings.ToLower(r.Name),
		RecordType: strings.ToUpper(r.Type),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}

	if r.UUID.String() != uuid.Nil.String() {
		dbModel.ID = r.UUID.String()
	}

	logger.Debugw("api record converted to db model", "record", r, "db-model", dbModel)

	return dbModel, nil
}

// NewRecord creates a record from the URL params and validates it
func NewRecord(c *gin.Context) (*Record, error) {
	// Try to get record info from URL params
	record := &Record{
		Name: strings.ToLower(c.Param("name")),
		Type: strings.ToUpper(c.Param("type")),
	}
	record.path = record.Name + "/" + record.Type

	record.Answers = []*answers.Answer{}

	// Try to validate
	if err := record.validate(); err != nil {
		return nil, err
	}

	return record, nil
}

func (r *Record) validate() error {
	if r.Name == "" {
		return ErrorNoRecordName
	}

	if r.Type == "" {
		return ErrorNoRecordType
	}

	if err := isSupportedRecordType(r.Type); err != nil {
		return err
	}

	return nil
}

func isSupportedRecordType(rtype string) error {
	var supportedTypes = map[string]bool{
		"A":   true,
		"SRV": true,
	}

	if _, ok := supportedTypes[rtype]; ok {
		return nil
	}

	return ErrorUnsupportedType
}
