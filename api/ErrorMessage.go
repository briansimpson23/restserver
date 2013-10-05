package api

import (
	"encoding/json"
	"fmt"
)

func BuildErrorMessage(statusCode int, message string, errorCode int) []byte {

	type ErrorMessage struct {
		StatusCode int
		Message    string
		ErrorCode  int
		MoreInfo   string
	}

	//---------------------------------------
	// populate the error message struct
	//---------------------------------------
	e := ErrorMessage{statusCode, message, errorCode, fmt.Sprintf("http://developer.mybucket.com/docs/errors/%d", errorCode)}

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = statusCode
	resp.Status = "error"
	resp.Data = e

	output, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return output
}
