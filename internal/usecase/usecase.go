package usecase

import (
	"context"
	"hash/fnv"
	"strings"

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
	var (
		short     string
		longCheck string
	)

	// Проверяем, существует ли уже short URL для данного long URL
	err := uc.s.GetShortUrl(id, long, &short, ctx)
	if err == nil {
		return short, nil
	}

	// Если записи нет, ищем следующий доступный id
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		for {
			err = uc.s.GetLongUrlByID(id, &longCheck)
			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				// Нашли свободный id
				break
			} else if err != nil {
				return "", err
			}
			// Если id занят, увеличиваем
			id++
		}

		short = urlshortener.Shorten(id)
		return short, uc.s.InsertShortUrl(id, long, short, ctx)
	}

	return "", err
}
func (uc *Usecase) GetLongUrl(short string, long *string, ctx context.Context) error {
	return uc.s.GetLongUrl(short, long, ctx)
}

func GetID(url string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(url))
	return hash.Sum32()
}
