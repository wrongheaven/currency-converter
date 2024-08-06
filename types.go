package main

type ApiResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`

	IsError      bool   `json:"error"`
	Error        string `json:"message"`
	ErrorMessage string `json:"description"`
}
