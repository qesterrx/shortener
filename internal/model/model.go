package model

type ShortenURLReq struct {
	URL string `json:"url"`
}

type ShortenURLRes struct {
	Goto string `json:"result"`
}
