package records

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
