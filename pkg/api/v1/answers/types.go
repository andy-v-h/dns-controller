package answers

import (
	"time"

	"github.com/google/uuid"
)

// Answer is the API model for an Answer
//
// Data relationships:
//   - An Answer may be many to one for Answers:Record
//   - Am Answer is one to one for Answer:Owner
type Answer struct {
	UUID       uuid.UUID `json:"uuid"`
	Target     string    `json:"target"`
	Type       string    `json:"type"`
	TTL        uint64    `json:"ttl"`
	HasDetails bool      `json:"has_details"`
	Details    []*Detail `json:"details,omitempty"`
	OwnerUUID  uuid.UUID `json:"owner_id"`
	RecordUUID uuid.UUID `json:"record_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// Detail is the API model for an Detail
//
// Data Relationships:
//   - An Detail is 1:1 with an Answer, however they are optional
type Detail struct {
	UUID       uuid.UUID `json:"uuid"`
	AnswerUUID uuid.UUID `json:"answer_id"`
	Port       int64     `json:"port"`
	Priority   int64     `json:"priority"`
	Protocol   string    `json:"protocol"`
	Weight     int64     `json:"weight"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
