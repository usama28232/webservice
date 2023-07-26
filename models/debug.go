package models

type HttpRequest struct {
	Method   string
	Url      string
	Data     string
	Agent    string
	Duration int64
}
