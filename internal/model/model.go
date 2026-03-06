package model

type ShortenURLReq struct {
	URL string `json:"url"`
}

type ShortenURLRes struct {
	Goto string `json:"result"`
}

type ShortenURL struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
