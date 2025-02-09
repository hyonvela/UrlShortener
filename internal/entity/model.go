package entity

type ShortUrl struct {
	ShortUrl string `json:"short_url" db:"short_url"`
}

type LongUrl struct {
	LongSUrl string `json:"long_url" db:"long_url"`
}
