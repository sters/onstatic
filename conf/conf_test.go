package conf

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap/zapcore"
)

func Test_createLogFileIfEmpty(t *testing.T) {
	p := getCurrentGoFilePath()
	p = filepath.Join(filepath.Dir(p), "/tmp.log")

	if err := createLogFileIfEmpty(p); err != nil {
		t.Errorf("failed: %+v", err)
	}

	if err := os.Remove(p); err != nil {
		t.Errorf("failed to cleanup: %+v", err)
	}
}

func Test_logger(t *testing.T) {
	dir := filepath.Dir(getCurrentGoFilePath())

	fileOut := filepath.Join(dir, "/tmp_out.log")
	fileErr := filepath.Join(dir, "/tmp_err.log")

	_ = logger(zapcore.DebugLevel, fileOut, fileErr)

	if err := os.Remove(fileOut); err != nil {
		t.Errorf("failed to cleanup: %+v", err)
	}
	if err := os.Remove(fileErr); err != nil {
		t.Errorf("failed to cleanup: %+v", err)
	}
}

func Test_getenv(t *testing.T) {
	var want string

	want = "default"
	if got := getenv("key", "default"); want != got {
		t.Errorf("want = %v, got = %v", want, got)
	}

	want = "defined"
	os.Setenv(prefix+"key", want)
	if got := getenv("key", "default"); want != got {
		t.Errorf("want = %v, got = %v", want, got)
	}
}

func Test_getenvInt(t *testing.T) {
	var want int

	want = 1
	if got := getenvInt("key", 1); want != got {
		t.Errorf("want = %v, got = %v", want, got)
	}

	want = 2
	os.Setenv(prefix+"key", "2")
	if got := getenvInt("key", 1); want != got {
		t.Errorf("want = %v, got = %v", want, got)
	}
}

func Test_getenvBool(t *testing.T) {
	var want bool

	want = false
	if got := getenvBool("key", false); want != got {
		t.Errorf("want = %v, got = %v", want, got)
	}

	want = true
	os.Setenv(prefix+"key", "true")
	if got := getenvBool("key", false); want != got {
		t.Errorf("want = %v, got = %v", want, got)
	}
}
