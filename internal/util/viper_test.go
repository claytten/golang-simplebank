package util_test

import (
	"testing"

	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	//testing true if app.env is in the root folder
	config, err := util.LoadConfig("app", "../..")
	require.NoError(t, err)
	require.NotEmpty(t, config)
	require.NotEmpty(t, config.Environment)
	require.NotEmpty(t, config.DBDriver)
	require.NotEmpty(t, config.DBSource)

	//testing false if app.env is not in the root folder
	config2, err2 := util.LoadConfig("example", "../..")
	require.Error(t, err2)
	require.Empty(t, config2)
}
