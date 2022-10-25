// Package answers wraps the CRUD operations for a models.Record
package answers

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.hollow.sh/dnscontroller/internal/models"
	"go.uber.org/zap"
)

func TestParseAnswers(t *testing.T) {
	SetLogger(zap.NewNop().Sugar())

	happyPathBody := io.NopCloser(strings.NewReader(`[{"target":"1.1.2.1","type":"a","has_details":false,"owner_id":"bf9455c6-6987-4fdb-96ac-7f3f9dfabbe4","record_id":"52087acc-b0e9-4060-bc48-f37182b6becc"}]`))
	badBody := io.NopCloser(strings.NewReader(`bad and boujie`))
	typeBody := io.NopCloser(strings.NewReader(`[{"target":"1.1.2.1","has_details":false,"owner_id":"bf9455c6-6987-4fdb-96ac-7f3f9dfabbe4","record_id":"52087acc-b0e9-4060-bc48-f37182b6becc"}]`))
	targetBody := io.NopCloser(strings.NewReader(`[{"has_details":false,"owner_id":"bf9455c6-6987-4fdb-96ac-7f3f9dfabbe4","record_id":"52087acc-b0e9-4060-bc48-f37182b6becc"}]`))

	answers := []*Answer{}

	type args struct {
		c   *gin.Context
		ans []*Answer
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Path",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Body: happyPathBody,
					},
				},
				ans: answers,
			},
			wantErr: false,
		},
		{
			name: "Bad payload",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Body: badBody,
					},
				},
				ans: answers,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Body: typeBody,
					},
				},
				ans: answers,
			},
			wantErr: true,
		},
		{
			name: "invalid target",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Body: targetBody,
					},
				},
				ans: answers,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		err := ParseAnswers(tt.args.c, tt.args.ans)

		if tt.wantErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestAnswer_ToDBModel(t *testing.T) {
	uuids := make(map[string]uuid.UUID)

	now := time.Time{}

	uuids["happy path"], _ = uuid.Parse("52087acc-b0e9-4060-bc48-f37182b6becc")

	type fields struct {
		UUID       uuid.UUID
		Target     string
		Type       string
		TTL        uint64
		HasDetails bool
		Details    []*Detail
		OwnerUUID  uuid.UUID
		RecordUUID uuid.UUID
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		want    *models.Answer
		wantErr bool
	}{
		{
			name: "Happy path",
			fields: fields{
				UUID:       uuids["happy path"],
				Target:     "example.COM",
				Type:       "SrV",
				TTL:        10,
				HasDetails: false,
				OwnerUUID:  uuids["happy path"],
				RecordUUID: uuids["happy path"],
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			want: &models.Answer{
				ID:         uuids["happy path"].String(),
				Target:     "example.com",
				Type:       "SRV",
				TTL:        10,
				HasDetails: false,
				OwnerID:    uuids["happy path"].String(),
				RecordID:   uuids["happy path"].String(),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Answer{
				UUID:       tt.fields.UUID,
				Target:     tt.fields.Target,
				Type:       tt.fields.Type,
				TTL:        tt.fields.TTL,
				HasDetails: tt.fields.HasDetails,
				Details:    tt.fields.Details,
				OwnerUUID:  tt.fields.OwnerUUID,
				RecordUUID: tt.fields.RecordUUID,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
			}
			got, err := a.ToDBModel()
			if (err != nil) != tt.wantErr {
				t.Errorf("Answer.ToDBModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Answer.ToDBModel()\n\t got:\t %v\n\t want:\t %v", got, tt.want)
			}
		})
	}
}
