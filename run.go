package rpcserver

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/designsbysm/timber/v2"
	"google.golang.org/grpc"
)

func Run(s *grpc.Server, host string, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("RPC: %s", err.Error())
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			panic(fmt.Errorf("RPC: %s", err.Error()))
		}
	}()

	timber.Info(fmt.Sprintf("RPC: listening on %s", address))

	// wait for ^c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// close
	fmt.Println("")
	s.Stop()
	listener.Close()
	timber.Info("RPC: closed")

	return nil
}
