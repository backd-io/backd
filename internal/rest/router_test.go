package rest

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {

	route := "/objects/:collection/:id/:relation/:direction"
	matches := []string{"", "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$", "^[a-zA-Z0-9-]{1,32}$", "^(in|out)$"}

	testItems := []map[string]interface{}{
		{
			"collection": "collectionOK",
			"id":         "idokidokidokidokidok",
			"rel":        "relationOK",
			"dir":        "in",
			"expect":     true,
		},
		{
			"collection": "collectionOK",
			"id":         "idnook", // invalid, no id
			"rel":        "relationOK",
			"dir":        "in",
			"expect":     false,
		},
		{
			"collection": "thisIsWayLongerThanExpectedSinceItMustBe32Chars", // invalid, longer
			"id":         "idokidokidokidokidok",
			"rel":        "relationOK",
			"dir":        "in",
			"expect":     false,
		},
		{
			"collection": "", // invalid, shorter
			"id":         "idokidokidokidokidok",
			"rel":        "relationOK",
			"dir":        "in",
			"expect":     false,
		},
		{
			"collection": "collectionOK",
			"id":         "idokidokidokidokidok",
			"rel":        "relationOK",
			"dir":        "invalid", // invalid, is not 'in' nor 'out'
			"expect":     false,
		},
	}

	for _, this := range testItems {
		req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/objects/%s/%s/%s/%s", this["collection"], this["id"], this["rel"], this["dir"]), nil)
		assert.Equal(t, this["expect"], match(route, matches, req), "they should be equal")
	}

}
