package model

type ShortenURLReq struct {
	Url string `json:"url"`
}

type ShortenURLRes struct {
	Goto string `json:"result"`
}
