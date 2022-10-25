// Package answers wraps the CRUD operations for a models.Record
package answers

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
