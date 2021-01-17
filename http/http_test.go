package http

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	s, err := NewServer("9999")
	if err != nil {
		t.Fatalf("failed to NewServer: %+v", err)
	}

	go func() {
		err := s.Run()
		if err != nil {
			t.Errorf("failed to Run: %+v", err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	if err := s.Close(); err != nil {
		t.Fatalf("failed to Close: %+v", err)
	}
}

func TestRunFail(t *testing.T) {
	s, err := NewServer("aaa")
	if err != nil {
		t.Fatalf("failed to NewServer: %+v", err)
	}

	go func() {
		err := s.Run()
		if err == nil {
			t.Errorf("Run want to fail")
		}
	}()

	time.Sleep(100 * time.Millisecond)

	if err := s.Close(); err != nil {
		t.Fatalf("failed to Close: %+v", err)
	}
}
