package specifications

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"golang.org/x/exp/slices"
)

func AssertEndpointsEqual(t *testing.T, got, want http.Endpoints) {
	t.Helper()
	assert.Equal(t, len(got), len(want))
	for i, endpoint := range got {
		AssertEndpointEqual(t, endpoint, want[i])
	}
}

func AssertEndpointEqual(t *testing.T, got, want http.Endpoint) {
	t.Helper()
	got.ID = uuid.UUID{}
	want.ID = uuid.UUID{}
	got.Response.Body = fudgeTheWhiteSpace(got.Response.Body)
	want.Response.Body = fudgeTheWhiteSpace(want.Response.Body)
	got.LoadedAt = time.Time{}
	want.LoadedAt = time.Time{}
	assert.Equal(t, got, want)
}

func AssertResponseMatches(t *testing.T, want, got http.Response) {
	t.Helper()
	assert.Equal(t, fudgeTheWhiteSpace(want.Body), fudgeTheWhiteSpace(got.Body), "body not equal")
	assert.Equal(t, want.Status, want.Status, "status not equal")

	for key, v := range want.Headers {
		for _, value := range v {
			i := slices.Index(got.Headers[key], value)
			assert.NotEqual(t, -1, i, fmt.Sprintf("%q not found in %v", value, got.Headers[key]))
		}
	}
}

func fudgeTheWhiteSpace(in string) string {
	in = strings.Replace(in, "\t", "", -1)
	in = strings.Replace(in, "\n", "", -1)
	in = strings.Replace(in, "\r", "", -1)
	in = strings.Replace(in, " ", "", -1)
	return in
}
