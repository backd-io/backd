package structs

import (
	"testing"
	"time"
)

func TestMetadata(t *testing.T) {

	var (
		t1, t2, t3 int64
		m          Metadata
		owner      string
		updater    string
		domain     string
	)

	t1 = time.Now().Unix()
	domain = "domain"
	owner = "oo"
	updater = "uu"

	m.SetCreate(domain, owner)
	time.Sleep(1 * time.Second)
	if got, want := m.CreatedAt, m.UpdatedAt; got != want {
		t.Errorf("Error Metadata: After creation created_at and updated at must be the same but got: '%v' and '%v'", got, want)
	}
	if got, want := t1, m.CreatedAt; got > want {
		t.Errorf("Error Metadata: Create time ('%v') must be older than test time ('%v')", got, want)
	}
	if got, want := t1, m.UpdatedAt; got > want {
		t.Errorf("Error Metadata: Update time ('%v') must be older than test time ('%v')", got, want)
	}
	if got, want := m.CreatedBy, FullUsername(domain, owner); got != want {
		t.Errorf("Error Metadata: Owner error: got: '%v', want: '%v'", got, want)
	}
	if got, want := m.UpdatedBy, FullUsername(domain, owner); got > want {
		t.Errorf("Error Metadata: Updater error: got: '%v', want: '%v'", got, want)
	}

	time.Sleep(1 * time.Second)
	t2 = time.Now().Unix()
	if got, want := t2, m.CreatedAt; got < want {
		t.Errorf("Error Metadata: Create ('%v') must not be older than 2nd test time ('%v')", want, got)
	}
	if got, want := t2, m.UpdatedAt; got < want {
		t.Errorf("Error Metadata: Update ('%v') must not be older than 2nd test time ('%v')", want, got)
	}

	m.SetUpdate(domain, updater)
	time.Sleep(1 * time.Second)
	t3 = time.Now().Unix()
	if got, want := m.CreatedAt, m.UpdatedAt; got == want {
		t.Errorf("Error Metadata: After update created_at and updated at must different but got: '%v' and '%v'", got, want)
	}
	if got, want := m.CreatedAt, m.UpdatedAt; got > want {
		t.Errorf("Error Metadata: After update created_at ('%v') must be older than updated_at ('%v')", got, want)
	}
	if got, want := t1, m.CreatedAt; got > want {
		t.Errorf("Error Metadata: Create time ('%v') must be older than test time ('%v')", got, want)
	}
	if got, want := t2, m.CreatedAt; got < want {
		t.Errorf("Error Metadata: Create ('%v') must not be older than 2nd test time ('%v')", want, got)
	}
	if got, want := t2, m.UpdatedAt; got > want {
		t.Errorf("Error Metadata: Update ('%v') must be older than 2nd test time ('%v')", want, got)
	}
	if got, want := t3, m.CreatedAt; got < want {
		t.Errorf("Error Metadata: Create time ('%v') must not be older than 3nd test time ('%v')", want, got)
	}
	if got, want := t3, m.UpdatedAt; got < want {
		t.Errorf("Error Metadata: 3rd test time ('%v') must be older than Update time ('%v')", got, want)
	}
	if got, want := m.CreatedBy, FullUsername(domain, owner); got != want {
		t.Errorf("Error Metadata: Owner error: got: '%v', want: '%v'", got, want)
	}
	if got, want := m.UpdatedBy, FullUsername(domain, updater); got > want {
		t.Errorf("Error Metadata: Updater error: got: '%v', want: '%v'", got, want)
	}
	m.SetUpdate(domain, owner)
	if got, want := m.CreatedBy, FullUsername(domain, owner); got != want {
		t.Errorf("Error Metadata: Owner error: got: '%v', want: '%v'", got, want)
	}
	if got, want := m.UpdatedBy, FullUsername(domain, owner); got > want {
		t.Errorf("Error Metadata: Updater error: got: '%v', want: '%v'", got, want)
	}
}
