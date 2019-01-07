package structs

import (
	"fmt"
	"reflect"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/scrypt"
)

// User is the struct that sets the `user` entity on the domain
//   for domains type `backd`:
//     - passwords are required, but if not set on creation it will return a random one
//     - activate and validated must be take in account
//   for domains type `active directory`:
//     - passwords, active and validated are meaningless
type User struct {
	ID                string                 `json:"_id" bson:"_id"`                        // (required - autogenerated)
	Username          string                 `json:"username" bson:"un"`                    // (required) Username is the entity that will be used for logon. If email will be used as username then both must match
	Name              string                 `json:"name" bson:"n"`                         // (required) Name of the user (optional, gets filled from the authorizators that returns it)
	Email             string                 `json:"email" bson:"e"`                        // (required) Email of the user (the one used to notify by mail)
	Description       string                 `json:"desc,omitempty" bson:"d"`               // (optional) Description
	Password          string                 `json:"password,omitempty" bson:"-"`           // (optional) Password is only used to get the initial password on user POST
	GeneratedPassword string                 `json:"generated_password,omitempty" bson:"-"` // GeneratedPassword will be filled only when password creation was ramdom
	PasswordKey       []byte                 `json:"-" bson:"pk"`                           // PasswordKey  can not be retrieved by using an API
	PasswordSalt      string                 `json:"-" bson:"ps"`                           // PasswordSalt can not be retrieved by using an API
	Active            bool                   `json:"active,omitempty" bson:"ac"`            // (required) Active defines when the user can interact with the APIs (some authorizations can leave it as active if the authentication system will allow or restrict the user)
	Validated         bool                   `json:"validated,omitempty" bson:"va"`         // (required) Validated shows if the user needs to make any action to active its email (and probably its account too)
	Data              map[string]interface{} `json:"data,omitempty" bson:"da,omitempty"`    // (optional) Data is the arbitrary information that can be stored for the user
	Groups            []string               `json:"groups,omitempty" bson:"-"`             // Groups is a commodity to include all the groups on the session
	Metadata          `json:"_meta" bson:"_meta"`
}

// UserValidator is the JSON schema validation for the applications collection
func UserValidator() map[string]interface{} {
	return BuildValidator(
		map[string]interface{}{
			"_id": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{20}$",
			},
			"un": map[string]interface{}{
				"bsonType":  "string",
				"maxLength": 254, // be coherent, email can have a maximun of 254
				"minLength": 2,
			},
			"n": map[string]interface{}{
				"bsonType":  "string",
				"maxLength": 254,
			},
			"e": map[string]interface{}{
				"bsonType":  "string",
				"maxLength": 254,
				"pattern":   `(?:[a-z0-9!#$%&'*+\/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+\/=?^_{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`,
			},
			"d": map[string]interface{}{
				"bsonType":  "string",
				"maxLength": 1000, // always is good to have a limit :)
			},
			"pk": map[string]interface{}{
				"bsonType": "binData",
			},
			"ps": map[string]interface{}{
				"bsonType": "string",
			},
			"ac": map[string]interface{}{
				"bsonType": "bool",
			},
			"va": map[string]interface{}{
				"bsonType": "bool",
			},
			"da": map[string]interface{}{
				"bsonType": "object",
			},
		},
		[]string{"_id", "un", "n", "e", "ac", "va"},
	)

}

// Indexes
var (
	UserIndexes = []Index{
		{
			Fields: []string{"_id"},
			Unique: true,
		},
	}
)

// AES 256 like
const (
	scryptN      = 16384
	scryptR      = 8
	scryptP      = 1
	scriptKeyLen = 32
)

// NewUser returns an initialized user struct
func NewUser() *User {
	return &User{
		Data: make(map[string]interface{}),
	}
}

// SetPassword sets the passwordy on the user struct
func (u *User) SetPassword(passwd string) error {

	var (
		err          error
		bytePassword []byte
	)

	// if there is no password then set one randomly
	if passwd == "" {
		return u.SetRandomPassword()
	}

	if u.PasswordSalt == "" {
		// func(length, numDigits, numSymbols int, noUpper, allowRepeat bool)
		u.PasswordSalt, err = password.Generate(16, 2, 2, false, true)
		if err != nil {
			return err
		}
	}

	bytePassword, err = u.password(passwd)

	if err == nil {
		u.PasswordKey = bytePassword[:]
	}

	return err

}

// SetRandomPassword creates a random password for the user
//   (16 alphanumeric characters)
func (u *User) SetRandomPassword() error {

	var (
		passwd string
		err    error
	)

	passwd, err = password.Generate(16, 2, 2, false, true)
	if err != nil {
		return err
	}

	u.GeneratedPassword = passwd
	return u.SetPassword(passwd)

}

// PasswordMatch verifies if password match with the stored one
func (u *User) PasswordMatch(passwd string) bool {

	var (
		err          error
		bytePassword []byte
	)

	bytePassword, err = u.password(passwd)
	if err != nil {
		return false
	}

	fmt.Println(string(bytePassword), string(u.PasswordKey))

	if reflect.DeepEqual(bytePassword, u.PasswordKey) {
		return true
	}

	return false

}

func (u *User) password(passwd string) ([]byte, error) {
	return scrypt.Key(
		[]byte(passwd),
		[]byte(u.PasswordSalt),
		scryptN,
		scryptR,
		scryptP,
		scriptKeyLen)
}
