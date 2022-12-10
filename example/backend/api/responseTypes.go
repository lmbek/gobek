package api

/*
	This file contains the type of responses
*/

type JSONDataResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type JSONMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
