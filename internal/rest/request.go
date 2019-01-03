package rest

import (
	// "encoding/json"
	"io/ioutil"
	"net/http"
)

// GetFromBody is a simple function to fill an object from the Request
func GetFromBody(r *http.Request, obj interface{}) error {

	defer r.Body.Close()
	objectByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(objectByte, &obj)

}
