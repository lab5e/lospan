package main

import (
	"net"
	"os"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/apiserver"
	"github.com/lab5e/lospan/pkg/keys"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":4711")
	if err != nil {
		lg.Error("Error creating listener: %v", err)
		os.Exit(1)
	}

	store := storage.NewMemoryStorage()
	ma, err := protocol.NewMA([]byte{1, 2, 3})
	if err != nil {
		lg.Error("Error creating MA: %v", err)
		os.Exit(1)
	}
	keyGen, err := keys.NewEUIKeyGenerator(ma, 0, store)
	if err != nil {
		lg.Error("Error creating EUI key generator: %v", err)
		os.Exit(1)
	}
	lospanSvc, err := apiserver.New(store, &keyGen)
	if err != nil {
		lg.Error("Error creatig lospan service: %v", err)
		os.Exit(1)
	}

	lg.Info("Listening on %s", listener.Addr().String())
	server := grpc.NewServer()
	lospan.RegisterLospanServer(server, lospanSvc)
	if err := server.Serve(listener); err != nil {
		lg.Error("Error serving gRPC: %v", err)
		os.Exit(2)
	}
}
