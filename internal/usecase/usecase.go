package usecase

import (
	"hash/fnv"

	"example.com/m/internal/storage"
	urlshortener "example.com/m/internal/url_shortener"
	"example.com/m/pkg/logging"
)

type Usecase struct {
	s   *storage.Storage
	log *logging.Logger
}

func New(storage *storage.Storage, logger *logging.Logger) *Usecase {
	return &Usecase{storage, logger}
}

func (uc *Usecase) ShortenUrl(url string) (string, error) {
	id := GetID(url)

	var short string
	err := uc.s.GetShortUrl(id, &short)

	uc.log.Info(short)
	if err != nil && err.Error() == "sql: no rows in result set" {
		short = urlshortener.Shorten(id)
		return short, uc.s.InsertShortUrl(id, url, short)
	}

	return short, err
}

func (uc *Usecase) GetLongUrl(short string, long *string) error {
	return uc.s.GetLongUrl(short, long)
}

func GetID(url string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(url))
	return hash.Sum32()
}
