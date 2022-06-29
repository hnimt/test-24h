package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseSuccess struct {
	Data interface{} `json:"data"`
}

func RespondSuccess(w http.ResponseWriter, httpCode int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	response, _ := json.Marshal(data)
	w.Write(response)
}
