package records

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.hollow.sh/dnscontroller/internal/models"
	"go.hollow.sh/dnscontroller/pkg/api/v1/answers"
	"go.uber.org/zap"
)

func TestNewRecord(t *testing.T) {
	var (
		u  uuid.UUID
		ts time.Time
	)

	type args struct {
		c *gin.Context
	}

	type test struct {
		name        string
		args        args
		expected    *Record
		expectedErr bool
		err         error
	}

	var tests []test

	SetLogger(zap.NewNop().Sugar())

	for _, types := range []string{"A", "SRV"} {
		tests = append(tests, test{
			name: "Happy path " + types,
			args: args{
				&[]gin.Context{
					{
						Params: gin.Params{
							{
								Key:   "name",
								Value: "example.com",
							},
							{
								Key:   "type",
								Value: types,
							},
						},
					},
				}[0],
			},
			expected: &Record{
				Name:      "example.com",
				Type:      types,
				Answers:   []*answers.Answer{},
				path:      "example.com/" + types,
				UUID:      u,
				CreatedAt: ts,
				UpdatedAt: ts,
			},
		})
	}

	tests = append(tests, test{
		name: "No record error",
		args: args{
			&[]gin.Context{
				{
					Params: gin.Params{},
				},
			}[0],
		},
		expectedErr: true,
		err:         ErrorNoRecordName,
	})

	tests = append(tests, test{
		name: "No recordtype error",
		args: args{
			&[]gin.Context{
				{
					Params: gin.Params{
						{
							Key:   "name",
							Value: "example.com",
						},
					},
				},
			}[0],
		},
		expectedErr: true,
		err:         ErrorNoRecordType,
	})

	tests = append(tests, test{
		name: "Unsupported recordtype error",
		args: args{
			&[]gin.Context{
				{
					Params: gin.Params{
						{
							Key:   "name",
							Value: "example.com",
						},
						{
							Key:   "type",
							Value: "TEAPOT",
						},
					},
				},
			}[0],
		},
		expectedErr: true,
		err:         ErrorUnsupportedType,
	})

	for _, tt := range tests {
		got, err := NewRecord(tt.args.c)
		if tt.expectedErr {
			assert.NotNil(t, err)
			assert.ErrorIs(t, tt.err, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, got)
		}
	}
}

func TestRecord_ToDBModel(t *testing.T) {
	uuids := make(map[string]uuid.UUID)

	now := time.Time{}

	uuids["happy path"], _ = uuid.Parse("52087acc-b0e9-4060-bc48-f37182b6becc")

	type fields struct {
		Name      string
		Type      string
		Answers   []*answers.Answer
		path      string
		UUID      uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		want    *models.Record
		wantErr bool
		err     error
	}{
		{
			name: "Happy Path",
			fields: fields{
				Name:      "example.com",
				Type:      "A",
				UUID:      uuids["happy path"],
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: &models.Record{
				ID:         uuids["happy path"].String(),
				Record:     "example.com",
				RecordType: "A",
				UpdatedAt:  now,
				CreatedAt:  now,
			},
		},
		{
			name: "invalid record",
			fields: fields{
				Name:      "",
				Type:      "A",
				UUID:      uuids["happy path"],
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			err:     ErrorNoRecordName,
		},
		{
			name: "invalid type",
			fields: fields{
				Name:      "example.com",
				Type:      "",
				UUID:      uuids["happy path"],
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: true,
			err:     ErrorNoRecordType,
		},
	}
	for _, tt := range tests {
		r := &Record{
			Name:      tt.fields.Name,
			Type:      tt.fields.Type,
			Answers:   tt.fields.Answers,
			path:      tt.fields.path,
			UUID:      tt.fields.UUID,
			CreatedAt: tt.fields.CreatedAt,
			UpdatedAt: tt.fields.UpdatedAt,
		}
		got, err := r.ToDBModel()

		if tt.wantErr {
			assert.NotNil(t, err)
			assert.Error(t, err, tt.err)
		} else {
			assert.Equal(t, tt.want, got)
		}
	}
}
