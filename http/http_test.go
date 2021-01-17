package http

import (
	"testing"
)

func TestRunFail(t *testing.T) {
	s, err := NewServer("aaa")
	if err != nil {
		t.Fatalf("failed to NewServer: %+v", err)
	}

	if err := s.Run(); err == nil {
		t.Errorf("Run want to fail")
	}

	if err := s.Close(); err != nil {
		t.Fatalf("failed to Close: %+v", err)
	}
}
