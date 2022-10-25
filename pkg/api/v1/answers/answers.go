// Package answers wraps the CRUD operations for a models.Record
package answers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"go.hollow.sh/dnscontroller/internal/models"
)

var logger *zap.SugaredLogger

// SetLogger initializes the package logger
func SetLogger(l *zap.SugaredLogger) { logger = l }

func qmAnswerTargetAndTypeAndRecord(atarget, atype string, recordUUID, ownerUUID uuid.UUID) qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = append(mods, qm.Where("record_id=?", recordUUID.String()))
	mods = append(mods, qm.Where("owner_id=?", ownerUUID.String()))
	mods = append(mods, qm.Where("target=?", atarget))
	mods = append(mods, qm.Where("type=?", atype))

	return qm.Expr(mods...)
}

// Delete removes a record from the DB
func (a *Answer) Delete(ctx context.Context, db *sqlx.DB) error {
	err := a.validate()
	if err != nil {
		return err
	}

	err = a.Find(ctx, db)
	if err != nil {
		return err
	}

	dbAnswer, err := a.ToDBModel()
	if err != nil {
		return err
	}

	_, err = dbAnswer.Delete(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

// FindOrCreate is the upsert function
func (a *Answer) FindOrCreate(ctx context.Context, db *sqlx.DB) error {
	err := a.Find(ctx, db)
	if errors.Is(err, sql.ErrNoRows) {
		return a.Create(ctx, db)
	} else if err != nil {
		return err
	}

	return nil
}

// Find looks the answer up by name,type
func (a *Answer) Find(ctx context.Context, db *sqlx.DB) error {
	if err := a.validate(); err != nil {
		return err
	}

	qm := qmAnswerTargetAndTypeAndRecord(a.Target, a.Type, a.RecordUUID, a.OwnerUUID)

	dbAnswer, err := models.Answers(qm).One(ctx, db)
	if err != nil {
		return err
	}

	return a.FromDBModel(ctx, db, dbAnswer)
}

// Create inserts a answer
func (a *Answer) Create(ctx context.Context, db *sqlx.DB) error {
	dbAnswer, err := a.ToDBModel()
	if err != nil {
		return err
	}

	if err := dbAnswer.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return a.FromDBModel(ctx, db, dbAnswer)
}

// FromDBModel converts a db type to an api type
func (a *Answer) FromDBModel(ctx context.Context, db *sqlx.DB, dbT *models.Answer) error {
	a.CreatedAt = dbT.CreatedAt
	a.UpdatedAt = dbT.UpdatedAt
	a.Target = dbT.Target
	a.Type = dbT.Type
	a.TTL = uint64(dbT.TTL)
	a.HasDetails = dbT.HasDetails

	var err error

	a.OwnerUUID, err = uuid.Parse(dbT.OwnerID)
	if err != nil {
		return err
	}

	a.RecordUUID, err = uuid.Parse(dbT.RecordID)
	if err != nil {
		return err
	}

	a.UUID, err = uuid.Parse(dbT.ID)
	if err != nil {
		return err
	}

	if a.HasDetails {
		if err = a.GetDetails(ctx, db); err != nil {
			return err
		}
	}

	logger.Debugw("db answer conterted to api model", "answer", dbT, "api-model", a)

	return a.validate()
}

// ToDBModel converts the api type to db type
func (a *Answer) ToDBModel() (*models.Answer, error) {
	dbModel := &models.Answer{
		Target:   strings.ToLower(a.Target),
		Type:     strings.ToUpper(a.Type),
		TTL:      int64(a.TTL),
		OwnerID:  a.OwnerUUID.String(),
		RecordID: a.RecordUUID.String(),
	}

	if err := a.validate(); err != nil {
		return nil, err
	}

	if a.UUID.String() != uuid.Nil.String() {
		dbModel.ID = a.UUID.String()
	}

	logger.Debugw("api answer converted to db model", "answer", a, "db-model", dbModel)

	return dbModel, nil
}

// ParseAnswers parses the answer list from the body
func ParseAnswers(c *gin.Context, answers []*Answer) error {
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &answers); err != nil {
		return err
	}

	for _, answer := range answers {
		a := answer
		if err := a.validate(); err != nil {
			return err
		}
	}

	logger.Debugw("got answers from body", "answers", answers, "body", string(bytes))

	return nil
}

func (a *Answer) validate() error {
	if a.Target == "" {
		return ErrorNoAnswerTarget
	}

	if a.Type == "" {
		return ErrorNoAnswerType
	}

	if err := isSupportedRecordType(a.Type); err != nil {
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
