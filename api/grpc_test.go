package api_test

import (
	"context"
	"testing"
	"url-shortener-golang/api"
	pb "url-shortener-golang/proto"
	"url-shortener-golang/storage"
	"url-shortener-golang/testutils"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCShorten(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupMock   func(*testutils.MockStorage)
		input       string
		expected    string
		expectedErr codes.Code
	}{
		{
			name: "Success",
			setupMock: func(ms *testutils.MockStorage) {
				ms.ExpectPost("https://valid.com", "valid12345", nil)
			},
			input:    "https://valid.com",
			expected: "valid12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := testutils.NewMockStorage()
			tt.setupMock(mockStorage)

			server := api.NewGRPCServer(mockStorage)

			resp, err := server.Shorten(ctx, &pb.ShortenRequest{
				OriginalUrl: tt.input,
			})

			if tt.expectedErr != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				if !ok {
					t.Fatal("Expected gRPC status error")
				}
				assert.Equal(t, tt.expectedErr, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp.ShortUrl)
			}
		})
	}
}

func TestGRPCResolve(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupMock   func(*testutils.MockStorage)
		input       string
		expected    string
		expectedErr codes.Code
	}{
		{
			name: "Success",
			setupMock: func(ms *testutils.MockStorage) {
				ms.ExpectGet("valid54321", "https://original.com", nil)
			},
			input:    "valid54321",
			expected: "https://original.com",
		},
		{
			name: "Not Found",
			setupMock: func(ms *testutils.MockStorage) {
				ms.ExpectGet("nothing", "", storage.ErrNotFound)
			},
			input:       "nothing",
			expectedErr: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := testutils.NewMockStorage()
			tt.setupMock(mockStorage)

			server := api.NewGRPCServer(mockStorage)

			resp, err := server.Resolve(ctx, &pb.ResolveRequest{
				ShortUrl: tt.input,
			})

			if tt.expectedErr != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				if !ok {
					t.Fatal("Expected gRPC status error")
				}
				assert.Equal(t, tt.expectedErr, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp.OriginalUrl)
			}
		})
	}
}
