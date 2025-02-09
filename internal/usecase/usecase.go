package usecase

import (
	"context"
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

func (uc *Usecase) ShortenUrl(long string, ctx context.Context) (string, error) {
	id := GetID(long)

	var short string
	err := uc.s.GetShortUrl(id, long, &short, ctx)

	if err != nil && err.Error() == "sql: no rows in result set" {
		short = urlshortener.Shorten(id)
		return short, uc.s.InsertShortUrl(id, long, short, ctx)
	}

	return short, err
}

func (uc *Usecase) GetLongUrl(short string, long *string, ctx context.Context) error {
	return uc.s.GetLongUrl(short, long, ctx)
}

func GetID(url string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(url))
	return hash.Sum32()
}
