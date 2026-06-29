package main

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestForwardHTTPStatus(t *testing.T) {
	ctx := runtime.NewServerMetadataContext(context.Background(), runtime.ServerMetadata{
		HeaderMD: metadata.Pairs("x-http-code", "204"),
	})
	recorder := httptest.NewRecorder()

	err := forwardHTTPStatus(ctx, recorder, nil)

	require.NoError(t, err)
	require.Equal(t, 204, recorder.Code)
	require.Empty(t, recorder.Header().Get("Grpc-Metadata-X-Http-Code"))
}
