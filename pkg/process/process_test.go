package process_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yolossn/process-manager/pkg/config"
	"github.com/yolossn/process-manager/pkg/process"
)

func TestNew(t *testing.T) {
	t.Parallel()

	command := config.Command{
		Command:    "ls",
		Args:       []string{"/"},
		MaxRetries: 4,
	}

	proc := process.New(command)

	require.Equal(t, command.MaxRetries, proc.MaxRetries())
}

func TestRun(t *testing.T) {
	t.Parallel()

	// Success case
	command := config.Command{
		Command:    "curl",
		Args:       []string{"https://postman-echo.com/get?foo=bar"},
		MaxRetries: 1,
	}

	proc := process.New(command)
	completeChan := make(chan *process.Process, 1)
	proc.Run(completeChan)
	<-completeChan

	require.Contains(t, proc.Output(), "foo")
	require.Contains(t, proc.Output(), "bar")
	require.Equal(t, 0, proc.Retries())
	require.Equal(t, true, proc.IsSuccessful())

	// Fail case
	command = config.Command{
		Command:    "sleep",
		Args:       []string{"10", "2"},
		MaxRetries: 3,
	}

	proc = process.New(command)
	proc.Run(completeChan)
	<-completeChan

	require.Equal(t, 2, proc.Retries())
	require.Equal(t, false, proc.IsSuccessful())
	require.Contains(t, proc.Error(), "usage: sleep seconds")

	// Stop
	command = config.Command{
		Command:    "sleep",
		Args:       []string{"10"},
		MaxRetries: 3,
	}
	proc = process.New(command)
	time.Sleep(time.Second)
	proc.Stop()
	proc.Run(completeChan)
	<-completeChan

	require.Equal(t, false, proc.IsSuccessful())
}
