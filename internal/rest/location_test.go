package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {

	// assert equality
	assert.Equal(t, "users/123/test", Location("users", "123", "test"), "they should be equal")

}
