package api

import (
	"encoding/json"
)

func ReturnError(errorMessage string) []byte {

	errmsg := make(map[string]string)
	errmsg["ERROR"] = errorMessage

	envelope := make(map[string]map[string]string)
	envelope["ErrorMessage"] = errmsg

	output, _ := json.Marshal(envelope)

	return output

}
