package backoff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yolossn/process-manager/pkg/backoff"
)

func TestNewStaticBackoff(t *testing.T) {
	bo := backoff.NewStaticBackoff(time.Second)

	sleep := bo.Duration()
	require.Equal(t, time.Second, sleep)
}
