package apiserver

import (
	"github.com/lab5e/lospan/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toProtoErr(err error) error {
	if err == storage.ErrAlreadyExists {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if err == storage.ErrNotFound {
		return status.Error(codes.NotFound, err.Error())
	}
	if err == storage.ErrDeleteConstraint {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
