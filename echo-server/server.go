package main

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

type ServerConfig struct {
	Port         int
	AuthorizedId string
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10)

	defer cancel()

	serverConfig := ServerConfig{
		Port:         8080,
		AuthorizedId: "",
	}

	RunServer(ctx, serverConfig)

}

func RunServer(ctx context.Context, config ServerConfig) error {

	allowedId := spiffeid.MustSpiffeID("spiffe://example.org/spire/agent/server")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	for {
		connection, err := listen.Accept()
		go func() {
			defer connection.Close()
			if err != nil {
				return
			}
			if err := handleConnection(connection); err != nil {
				return
			}
		}()
	}
	return nil

}

func handleConnection(connection net.Conn) error {
	buffReader := bufio.NewReader(connection)
	request, err := buffReader.ReadString('\n')
	buffWriter := bufio.NewWriter(connection)
	if err != nil {
		return err
	}
	buffWriter.WriteString(fmt.Sprintf("%s\n", request))
	buffWriter.Flush()
	return nil
}
