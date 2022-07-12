package logr

import (
	"io/ioutil"
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
	log := Get()
	log.Info("hello!")
	require.NoError(t, w.Close())
	os.Stdout = previousStdout

	loggedLine, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	assert.Contains(t, string(loggedLine), `msg="hello!"`)
	assert.Contains(t,  string(loggedLine), FieldFileName+"=logr/logr_test.go")
}
