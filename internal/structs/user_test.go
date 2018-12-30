package structs

import (
	"reflect"
	"testing"
)

func TestUserPassword(t *testing.T) {

	u := NewUser()

	passWrong := "wrong password"

	passList := []string{
		"Hello, it's me",
		"boringPassword",
		"84-3T[wZHcx*';k;=m",
		"指事字 zhǐshìzì",
		" الأَبْجَدِيَّة العَرَبِيَّة",
	}

	for _, pass := range passList {
		if err := u.SetPassword(pass); err != nil {
			t.Errorf("Error setting password Got: %v, instead of nil", err)
		}

		if got := u.PasswordMatch(pass); got == false {
			t.Errorf("Error matching password, Got: %v, Want: true", got)
		}

		if got := u.PasswordMatch(passWrong); got == true {
			t.Errorf("Error matching password, Got: %v, Want: false", got)
		}
	}

	if err := u.SetRandomPassword(); err != nil {
		t.Errorf("Error setting random password. Got %v, instead of nil", err)
	}

	if got := u.PasswordMatch(u.GeneratedPassword); got == false {
		t.Errorf("Error matching password, Got: %v, Want: true", got)
	}

}

func TestNewUser(t *testing.T) {

	u := NewUser()

	if got, want := reflect.TypeOf(u).String(), "*structs.User"; got != want {
		t.Errorf("Error NewUser() does not return '%s', got: %v", want, got)
	}

	if got := reflect.TypeOf(u.Data).String(); got != "map[string]interface {}" {
		t.Errorf("Error NewUser() does not return a *User, got: %v", got)
	}

}
