package main

import (
	"context"
	"time"

	"github.com/lab5e/lospan/pkg/pb/lospan"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createClient(addr string) (lospan.LospanClient, context.Context, context.CancelFunc, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, err
	}
	client := lospan.NewLospanClient(conn)
	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	return client, ctx, done, nil
}
