package logr

import (
	"io/ioutil"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	previousStdout := os.Stdout
	defer func() {
		os.Stdout = previousStdout
	}()
	os.Stdout = w

	// Init must be called after redirecting os.Stdout, so the slog handler
	// captures the pipe's file descriptor instead of the original stdout.
	Init(slog.LevelDebug)
	log := Get()
	log.Info("hello!")
	require.NoError(t, w.Close())
	os.Stdout = previousStdout

	loggedLine, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	assert.Contains(t, string(loggedLine), `msg=hello!`)
	assert.Contains(t, string(loggedLine), "source=")
	assert.Contains(t, string(loggedLine), "logr/logr_test.go")
}
