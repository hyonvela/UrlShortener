package tests

import (
	"context"
	"net"
	"testing"

	"example.com/m/internal/handlers"
	"example.com/m/internal/storage"
	grpcV1 "example.com/m/pkg/grpc.v1"
	"example.com/m/pkg/logging"
	"example.com/m/tests/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestGRPCHandlers(t *testing.T) {
	mockUsecase := new(mocks.MockUsecase)
	mockLogger := logging.GetLogger("text", "test")
	grpcHandler := handlers.NewGRPCHandler(mockUsecase, mockLogger)

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	grpcV1.RegisterUrlShortenerServiceServer(srv, grpcHandler)

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := grpcV1.NewUrlShortenerServiceClient(conn)

	t.Run("ShortenUrl success", func(t *testing.T) {
		mockUsecase.ShortURL = "abc123"
		mockUsecase.Err = nil

		resp, err := client.ShortenUrl(context.Background(), &grpcV1.UrlRequest{Url: "https://example.com"})
		assert.NoError(t, err)
		assert.Equal(t, "abc123", resp.Url)
	})

	t.Run("ShortenUrl error", func(t *testing.T) {
		mockUsecase.Err = storage.ErrNotFound

		_, err := client.ShortenUrl(context.Background(), &grpcV1.UrlRequest{Url: "https://example.com"})
		assert.Error(t, err)
	})

	t.Run("GetLongUrl success", func(t *testing.T) {
		mockUsecase.LongURL = "https://example.com"
		mockUsecase.Err = nil

		resp, err := client.GetLongUrl(context.Background(), &grpcV1.UrlRequest{Url: "abc123"})
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", resp.Url)
	})

	t.Run("GetLongUrl error", func(t *testing.T) {
		mockUsecase.Err = storage.ErrNotFound

		_, err := client.GetLongUrl(context.Background(), &grpcV1.UrlRequest{Url: "invalid"})
		assert.Error(t, err)
	})
}
