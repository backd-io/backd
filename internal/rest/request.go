package rest

import (
	// "encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/backd-io/backd/internal/constants"
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

// QueryStrings returns query and sort formatted for db request
func QueryStrings(r *http.Request) (query, sort map[string]interface{}, skip, limit int64, err error) {

	var (
		queryString   string
		sortString    string
		pageString    string
		perPageString string
		pageNumber    int64
		perPageI      int64
	)

	queryString = r.URL.Query().Get("q")
	sortString = r.URL.Query().Get("sort")
	pageString = r.URL.Query().Get("page")
	perPageString = r.URL.Query().Get("per_page")

	// if queryString is empty return a nil interface (do nothing)
	if queryString != "" {
		err = json.Unmarshal([]byte(queryString), &query)
		if err != nil {
			return
		}
	}

	// if sortString is empty it will return a sort array already empty
	if sortString != "" {
		err = json.Unmarshal([]byte(sortString), &sort)
		if err != nil {
			return
		}
	}

	// pagination
	if pageString == "" {
		pageNumber = 1
	} else {

		pageNumber, err = strconv.ParseInt(pageString, 10, 64)
		if err != nil {
			skip = 0
			limit = constants.DefaultPerPage
			return
		}
	}

	if perPageString == "" {
		perPageI = constants.DefaultPerPage
	} else {
		perPageI, err = strconv.ParseInt(pageString, 10, 64)
		if err != nil {
			skip = 0
			limit = constants.DefaultPerPage
			return
		}
	}

	skip = (pageNumber - 1) * perPageI
	limit = pageNumber * perPageI

	return
}
