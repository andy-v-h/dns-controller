package records

import "errors"

var (
	// ErrorRecordNoAnswers is when no answers for a record are found
	ErrorRecordNoAnswers = errors.New("no answers found for a record")
	// ErrorInvalidRecord is a generic invalid response
	ErrorInvalidRecord = errors.New("invalid record format")
	// ErrorNoRecordName is when a request / record doesn't have a primary key
	ErrorNoRecordName = errors.New("no record name")
	// ErrorNoRecordType is when a request / record doesn't have a primary key
	ErrorNoRecordType = errors.New("no record type")
	// ErrorUnsupportedType when a request for an unsupported record type occurs
	ErrorUnsupportedType = errors.New("unsupported record type")
)
