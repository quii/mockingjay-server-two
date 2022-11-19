package specifications

import (
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

func AssertEndpointsEqual(t *testing.T, got, want mockingjay.Endpoints) {
	t.Helper()
	assert.Equal(t, len(got), len(want))
	for i, endpoint := range got {
		AssertEndpointEqual(t, endpoint, want[i])
	}
}

func AssertEndpointEqual(t *testing.T, got, want mockingjay.Endpoint) {
	t.Helper()
	got.ID = uuid.UUID{}
	want.ID = uuid.UUID{}
	got.Response.Body = fudgeTheWhiteSpace(got.Response.Body)
	want.Response.Body = fudgeTheWhiteSpace(want.Response.Body)
	got.LoadedAt = time.Time{}
	want.LoadedAt = time.Time{}
	assert.Equal(t, got, want)
}

func fudgeTheWhiteSpace(in string) string {
	in = strings.Replace(in, "\t", "", -1)
	in = strings.Replace(in, "\n", "", -1)
	in = strings.Replace(in, "\r", "", -1)
	in = strings.Replace(in, " ", "", -1)
	return in
}
