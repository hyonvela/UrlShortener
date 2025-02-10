package usecase

import (
	"context"
	"hash/fnv"

	"example.com/m/internal/storage"
	urlshortener "example.com/m/internal/url_shortener"
	"example.com/m/pkg/logging"
)

type Usecase struct {
	s   storage.Storage
	log *logging.Logger
}

func New(storage storage.Storage, logger *logging.Logger) *Usecase {
	return &Usecase{storage, logger}
}

func (uc *Usecase) ShortenUrl(long string, ctx context.Context) (string, error) {
	id := GetID(long)

	// Так как коллизий избежать не получится из-за того что мощность
	// множества всех строк больше мощьности множества строк которые производит алгоритм сокращения,
	// необходимо предусмотреть обход коллизий

	// Проверяем, существует ли уже короткий url
	short, err := uc.s.GetShortURL(ctx, long)
	if err == nil {
		return short, err
	}

	// Если записи нет, ищем следующий доступный id
	if err == storage.ErrNotFound {
		for {
			_, err := uc.s.GetLongURLByID(ctx, id)
			if err == storage.ErrNotFound {
				// Нашли свободный id
				break
			} else if err != nil {
				return "", err
			}
			// Если id занят, увеличиваем
			id++
		}

		short = urlshortener.Shorten(id)
		return short, uc.s.SaveURL(ctx, id, long, short)
	}

	return "", err
}

func (uc *Usecase) GetLongUrl(short string, ctx context.Context) (string, error) {
	return uc.s.GetLongURL(ctx, short)
}

func GetID(url string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(url))
	return hash.Sum32()
}
