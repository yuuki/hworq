package db

import "testing"

func TestNew(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if db == nil {
		t.Errorf("failed to connect postgresql")
	}
}
