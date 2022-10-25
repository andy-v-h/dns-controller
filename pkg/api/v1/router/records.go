package router

import (
	"github.com/gin-gonic/gin"

	ax "go.hollow.sh/dnscontroller/pkg/api/v1/answers"
	rx "go.hollow.sh/dnscontroller/pkg/api/v1/records"
)

func (r *Router) deleteRecord(c *gin.Context) {
	record, err := rx.NewRecord(c)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	if err := record.Delete(c.Request.Context(), r.db); err != nil {
		badRequestResponse(c, "failed to delete record", err)
		return
	}

	deletedResponse(c)
}

func (r *Router) createRecord(c *gin.Context) {
	record, err := rx.NewRecord(c)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	err = record.Create(c.Request.Context(), r.db)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	createdResponse(c)
}

func (r *Router) getRecord(c *gin.Context) {
	record, err := rx.NewRecord(c)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	err = record.Find(c.Request.Context(), r.db)
	if err != nil {
		dbErrorResponse(c, err)
		return
	}

	successResponse(c, record)
}

func (r *Router) createRecordAnswers(c *gin.Context) {
	record, err := rx.NewRecord(c)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	err = record.Find(c.Request.Context(), r.db)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	answers := []*ax.Answer{}

	err = ax.ParseAnswers(c, answers)
	if err != nil {
		badRequestResponse(c, ax.ErrorInvalidAnswers.Error(), err)
		return
	}

	r.logger.Debugw("pasred answers for record", "answers", answers, "record", record)

	var preparedAnswers []*ax.Answer

	for i := range answers {
		tempA := answers[i]
		tempA.RecordUUID = record.UUID

		preparedAnswers = append(preparedAnswers, tempA)

		r.logger.Debugw("prepped answer", "answer", answers[i], preparedAnswers[i])
	}

	for i := range preparedAnswers {
		if err := preparedAnswers[i].Create(c.Request.Context(), r.db); err != nil {
			badRequestResponse(c, "failed to write answers to datastore", err)
			return
		}
	}

	// Clear the cached value and go fetch the record with new answers
	record, err = rx.NewRecord(c)
	if err != nil {
		badRequestResponse(c, rx.ErrorInvalidRecord.Error(), err)
		return
	}

	if err := record.Find(c.Request.Context(), r.db); err != nil {
		dbErrorResponse(c, err)
		return
	}

	r.logger.Debugw("answers created for record", "answers", answers, "record", record)

	createdResponse(c)
}
