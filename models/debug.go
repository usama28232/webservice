package models

type Param struct {
	Debug bool
}

type HttpRequest struct {
	Method   string
	Url      string
	Data     string
	Agent    string
	Duration int64
}
