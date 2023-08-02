package models

type Param struct {
	Debug bool
}

type HttpRequest struct {
	Trxid     string
	SessionId string
	Method    string
	Url       string
	Data      string
	Agent     string
	Duration  int64
	Status    int
}
