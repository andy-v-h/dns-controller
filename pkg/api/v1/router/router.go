// Package router has a router for dnscontroller
package router

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.hollow.sh/toolbox/ginjwt"
	"go.uber.org/zap"

	ax "go.hollow.sh/dnscontroller/pkg/api/v1/answers"
	rx "go.hollow.sh/dnscontroller/pkg/api/v1/records"
)

const (
	// V1URI is the path prefix for all v1 endpoints
	V1URI = "/api/v1"

	// RecordsBaseURI is the path to the regular record
	// endpoint, called by the client to .
	RecordsBaseURI = "/records"

	// RecordsNameTypeURI is the path to the endpoint used
	// for retrieving the stored record for an instance
	RecordsNameTypeURI = RecordsBaseURI + "/:name/:type"

	// RecordAnswersURI is for interactions with record's answers
	RecordAnswersURI = RecordsNameTypeURI + "/answers"

	// RecordAnswersIDURI is for interactions with record's answers
	RecordAnswersIDURI = RecordAnswersURI + "/:uuid"

	// AnswersBaseURI is the path to the regular record
	// endpoint, called by the client to .
	AnswersBaseURI = "/answers"

	// AnswersIDURI is for interactions with record's answers
	AnswersIDURI = AnswersBaseURI + "/:uuid"

	// AnswerDetailsURI gets details about the an answer by ID
	AnswerDetailsURI = AnswersIDURI + "/details"

	// AnswerDetailsIDURI

	// scopePrefix = "dnscontroller"
)

// Router provides a router for the v1 API
type Router struct {
	authMW *ginjwt.Middleware
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

// New builds a Router
func New(amw *ginjwt.Middleware, db *sqlx.DB, l *zap.SugaredLogger) *Router {
	ax.SetLogger(l)
	rx.SetLogger(l)

	return &Router{authMW: amw, db: db, logger: l}
}

// Routes will add the routes for this API version to a router group
func (r *Router) Routes(rg *gin.RouterGroup) {
	// TODO: add auth'd endpoints
	// authMw := r.AuthMW
	// rg.POST(RecordURI, authMw.AuthRequired(), authMw.RequiredScopes(upsertScopes("record")))
	rg.GET(RecordsNameTypeURI, r.getRecord)
	rg.POST(RecordsNameTypeURI, r.createRecord)
	rg.DELETE(RecordsNameTypeURI, r.deleteRecord)

	rg.POST(RecordAnswersURI, r.createRecordAnswers)
}

// GetRecordPath returns the path used by an instance to fetch Record
func GetRecordPath() string {
	return path.Join(V1URI, RecordsNameTypeURI)
}

// func upsertScopes(items ...string) []string {
// 	s := []string{"write", "create", "update"}
// 	for _, i := range items {
// 		s = append(s, fmt.Sprintf("%s:create:%s", scopePrefix, i))
// 	}

// 	for _, i := range items {
// 		s = append(s, fmt.Sprintf("%s:update:%s", scopePrefix, i))
// 	}

// 	return s
// }

// func readScopes(items ...string) []string {
// 	s := []string{"read"}
// 	for _, i := range items {
// 		s = append(s, fmt.Sprintf("%s:read:%s", scopePrefix, i))
// 	}

// 	return s
// }

// func deleteScopes(items ...string) []string {
// 	s := []string{"write", "delete"}
// 	for _, i := range items {
// 		s = append(s, fmt.Sprintf("%s:delete:%s", scopePrefix, i))
// 	}

// 	return s
// }
