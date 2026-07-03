package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

func ForwardHTTPStatus(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	serverMetadata, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	values := serverMetadata.HeaderMD.Get("x-http-code")
	if len(values) == 0 {
		return nil
	}

	statusCode, err := strconv.Atoi(values[0])
	if err != nil {
		return fmt.Errorf("parse HTTP status code: %w", err)
	}

	delete(serverMetadata.HeaderMD, "x-http-code")
	w.Header().Del("Grpc-Metadata-X-Http-Code")
	w.WriteHeader(statusCode)
	return nil
}
