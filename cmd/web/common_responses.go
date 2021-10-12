package main

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}
