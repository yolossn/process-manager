package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yolossn/process-manager/pkg/config"
)

func TestFromFile(t *testing.T) {
	t.Parallel()

	// File exists
	conf, err := config.FromYaml("../../config.yaml")
	require.NoError(t, err)

	require.Equal(t, 5, len(conf.Commands))
	require.Equal(t, []string{"PWD=/Users/santhoshnagarajs/git/"}, conf.Commands[0].EnvStrings())
	// File not exists
	conf, err = config.FromYaml("../../not_exists.yaml")
	require.Error(t, err)
	require.Nil(t, conf)

}
