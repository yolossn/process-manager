package manager_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yolossn/process-manager/pkg/config"
	"github.com/yolossn/process-manager/pkg/manager"
)

func TestNewManager(t *testing.T) {
	t.Parallel()

	commands := readConfig(t, "../../config.yaml")
	man := manager.New(commands)
	require.NotNil(t, man)
}

func TestRun(t *testing.T) {
	t.Parallel()

	// Complete run
	commands := readConfig(t, "../../test/data/test_config.yaml")
	man := manager.New(commands)
	require.NotNil(t, man)

	done := man.Run()

	<-done
	t.Log(man.FailCount())
	t.Log(man.SuccessCount())

	require.Equal(t, 1, man.FailCount())
	require.Equal(t, 3, man.SuccessCount())

	// Interrupt run
	man = manager.New(commands)
	require.NotNil(t, man)
	done = man.Run()

	time.Sleep(time.Second)
	man.Stop()
	<-done
	require.NotEqual(t, 1, man.FailCount())
	require.NotEqual(t, 3, man.SuccessCount())
}

func readConfig(t *testing.T, file string) []config.Command {
	t.Helper()

	conf, err := config.FromYaml(file)
	require.NoError(t, err)
	return conf.Commands
}
