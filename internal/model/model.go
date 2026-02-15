package model

type ShortenUrlReq struct {
	Url string `json:"url"`
}

type ShortenUrlRes struct {
	Goto string `json:"result"`
}
