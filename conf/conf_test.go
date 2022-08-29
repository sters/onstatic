package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func Test_createLogFileIfEmpty(t *testing.T) {
	p := getCurrentGoFilePath()
	p = filepath.Join(filepath.Dir(p), "/tmp.log")

	assert.NoError(t, createLogFileIfEmpty(p))
	assert.NoError(t, os.Remove(p))
}

func Test_logger(t *testing.T) {
	dir := filepath.Dir(getCurrentGoFilePath())

	fileOut := filepath.Join(dir, "/tmp_out.log")
	fileErr := filepath.Join(dir, "/tmp_err.log")

	_ = logger(zapcore.DebugLevel, fileOut, fileErr)

	assert.NoError(t, os.Remove(fileOut))
	assert.NoError(t, os.Remove(fileErr))
}

func Test_getenv(t *testing.T) {
	var want string

	want = "default"
	assert.Equal(t, want, getenv("key", "default"))
	assert.Equal(t, want, getenv("key", "default"))

	want = "defined"
	t.Setenv(prefix+"key", want)
	assert.Equal(t, want, getenv("key", "default"))
}

func Test_getenvInt(t *testing.T) {
	var want int

	want = 1
	assert.Equal(t, want, getenvInt("key", 1))

	want = 2
	t.Setenv(prefix+"key", "2")
	assert.Equal(t, want, getenvInt("key", 1))
}

func Test_getenvBool(t *testing.T) {
	var want bool

	want = false
	assert.Equal(t, want, getenvBool("key", false))

	want = true
	t.Setenv(prefix+"key", "true")
	assert.Equal(t, want, getenvBool("key", false))
}
