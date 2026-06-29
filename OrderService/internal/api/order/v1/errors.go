package api

import (
	"errors"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapCreateOrderError(err error) error {
	if errors.Is(err, model.ErrPartNotFound) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return preserveGRPCStatusOrInternal(err, "failed to create order")
}

func mapGetOrderError(err error) error {
	if errors.Is(err, model.ErrOrderNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	return preserveGRPCStatusOrInternal(err, "failed to get order")
}

func mapPayOrderError(err error) error {
	switch {
	case errors.Is(err, model.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrInvalidOrderStatus), errors.Is(err, model.ErrOrderAlreadyPaid):
		return status.Error(codes.Aborted, err.Error())
	default:
		return preserveGRPCStatusOrInternal(err, "failed to pay order")
	}
}

func preserveGRPCStatusOrInternal(err error, message string) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	return status.Errorf(codes.Internal, "%s: %v", message, err)
}
