package records

import (
	"time"

	"github.com/google/uuid"

	"go.hollow.sh/dnscontroller/pkg/api/v1/answers"
)

// Record is the API model for a record
type Record struct {
	Name      string            `json:"record"`
	Type      string            `json:"record_type"`
	Answers   []*answers.Answer `json:"answers,omitempty"`
	path      string
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
