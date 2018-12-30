package db

import "testing"

func TestNewXID(t *testing.T) {

	xid := NewXID()

	if got, want := len(xid.String()), 20; got != want {
		t.Errorf("ID must have %d characters but got: %v", want, got)
	}

}
