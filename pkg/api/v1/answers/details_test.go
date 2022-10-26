package answers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"go.hollow.sh/dnscontroller/internal/models"
)

func TestDetail_ToDBModel(t *testing.T) {
	uuids := make(map[string]uuid.UUID)

	now := time.Time{}

	uuids["happy path"], _ = uuid.Parse("52087acc-b0e9-4060-bc48-f37182b6becc")

	type fields struct {
		UUID       uuid.UUID
		AnswerUUID uuid.UUID
		Port       int64
		Priority   int64
		Protocol   string
		Weight     int64
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		want    *models.AnswerDetail
		wantErr bool
		err     error
	}{
		{
			name: "happy path",
			fields: fields{
				UUID:       uuids["happy path"],
				AnswerUUID: uuids["happy path"],
				Port:       1337,
				Priority:   10,
				Protocol:   "tcp",
				Weight:     10,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			want: &models.AnswerDetail{
				ID:        uuids["happy path"].String(),
				AnswerID:  uuids["happy path"].String(),
				Priority:  null.Int64From(10),
				Port:      null.Int64From(1337),
				Protocol:  null.StringFrom("tcp"),
				Weight:    null.Int64From(10),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "null answer id",
			fields: fields{
				UUID:       uuids["happy path"],
				AnswerUUID: uuid.Nil,
				Port:       1337,
				Priority:   10,
				Protocol:   "tcp",
				Weight:     10,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			err:     ErrorNoAnswerDetailAnswerID,
			wantErr: true,
		},
		{
			name: "null details id",
			fields: fields{
				AnswerUUID: uuids["happy path"],
				UUID:       uuid.Nil,
				Port:       1337,
				Priority:   10,
				Protocol:   "tcp",
				Weight:     10,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			err:     ErrorNoAnswerDetailAnswerID,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		d := &Detail{
			ID:        tt.fields.UUID,
			AnswerID:  tt.fields.AnswerUUID,
			Port:      tt.fields.Port,
			Priority:  tt.fields.Priority,
			Protocol:  tt.fields.Protocol,
			Weight:    tt.fields.Weight,
			CreatedAt: tt.fields.CreatedAt,
			UpdatedAt: tt.fields.UpdatedAt,
		}
		got, err := d.ToDBModel()

		if tt.wantErr {
			assert.NotNil(t, err)
			assert.Error(t, err, tt.err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		}
	}
}
