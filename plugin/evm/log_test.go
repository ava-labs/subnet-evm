package evm

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	require := require.New(t)
	_, err := InitLogger("alias", "info", true, os.Stderr)
	require.NoError(err)
	log.Info("test")
}
