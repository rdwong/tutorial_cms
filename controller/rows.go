package controller

import (
	"net/http"

	"github.com/lightwell/cms_utilities_go_v2/cms"
)

var rowToModelLookup = map[string]func(columnNameToRowIndex map[string]int, rowValues [][]interface{}) ([]interface{}, error){}

// Add endpoints specific to this CMS here
var customEndpoints = map[string]func(w http.ResponseWriter, r *http.Request){
	"ping": func(w http.ResponseWriter, r *http.Request) {},
	// NOTE: Add endpoints specific to this CMS here
}

func InitialiseLookupsEndpoints() {
	cms.SetCustomEndpoints(customEndpoints)
	cms.SetRowToModelLookup(rowToModelLookup)
}
