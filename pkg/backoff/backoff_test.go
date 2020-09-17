package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewStaticBackoff(t *testing.T) {
	bo := NewStaticBackoff(time.Second)

	sleep := bo.Duration()
	require.Equal(t, time.Second, sleep)
}
