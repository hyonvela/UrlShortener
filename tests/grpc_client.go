package tests

import (
	"context"
	"log"

	grpcV1 "example.com/m/internal/grpc.v1"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := grpcV1.NewUrlShortenerServiceClient(conn)

	shortenReq := &grpcV1.UrlRequest{Url: "https://example.com"}

	shortenResp, err := client.ShortenUrl(context.Background(), shortenReq)
	if err != nil {
		log.Fatalf("Ошибка при вызове ShortenUrl: %v", err)
	}

	log.Printf("Сокращенный URL: %s", shortenResp.Url)

	getLongReq := &grpcV1.UrlRequest{Url: shortenResp.Url}

	getLongResp, err := client.GetLongUrl(context.Background(), getLongReq)
	if err != nil {
		log.Fatalf("Ошибка при вызове GetLongUrl: %v", err)
	}

	log.Printf("Длинный URL: %s", getLongResp.Url)
}
