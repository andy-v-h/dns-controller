package answers

import "errors"

var (
	// ErrorInvalidAnswers is a generic invalid response
	ErrorInvalidAnswers = errors.New("invalid answers format")
	// ErrorNoAnswerTarget is when a request / answer doesn't have a primary key
	ErrorNoAnswerTarget = errors.New("no answer name")
	// ErrorNoAnswerType is when a request / answer doesn't have a primary key
	ErrorNoAnswerType = errors.New("no answer type")
	// ErrorNoAnswerDetail is when a request / answer doesn't have a primary key
	ErrorNoAnswerDetail = errors.New("no answer_details for answer provied")
	// ErrorNoAnswerDetailID there is no uuid set
	ErrorNoAnswerDetailID = errors.New("no uuid set")
	// ErrorNoAnswerDetailAnswerID there is forgien key
	ErrorNoAnswerDetailAnswerID = errors.New("no answer_id set")
	// ErrorUnsupportedType when a request for an unsupported record type occurs
	ErrorUnsupportedType = errors.New("unsupported record type")
)
