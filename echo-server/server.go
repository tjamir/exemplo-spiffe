package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

type ServerConfig struct {
	Port         int    `json:"port"`
	AuthorizedId string `json:"authorized_id"`
	SocketPath   string `json:"socket_path"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	serverConfig, err := loadConfig()

	if err != nil {
		panic(err)
	}

	err = RunServer(ctx, serverConfig)
	if err != nil {
		panic(err)
	}

}

func loadConfig() (ServerConfig, error) {
	serverConfig := ServerConfig{}
	configFile, err := os.Open("conf/echo-server.json")
	if err != nil {
		return serverConfig, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&serverConfig)
	if err != nil {
		return serverConfig, err
	}

	return serverConfig, nil
}

func RunServer(ctx context.Context, config ServerConfig) error {
	allowedId, err := spiffeid.FromString(config.AuthorizedId)
	if err != nil {
		return err
	}

	fmt.Println("Iniciando tls server")
	listen, err := spiffetls.ListenWithMode(ctx,
		"tcp", fmt.Sprintf(":%d", config.Port),
		spiffetls.MTLSServerWithSourceOptions(tlsconfig.AuthorizeID(allowedId),
			workloadapi.WithClientOptions(workloadapi.WithAddr(config.SocketPath))))

	if err != nil {
		return err
	}

	for {
		fmt.Println("Entrou no loop")
		connection, err := listen.Accept()
		fmt.Println("Aceitou")
		go func() {
			defer connection.Close()
			if err != nil {
				return
			}
			fmt.Println("Rodando a conexão")
			if err := handleConnection(connection); err != nil {
				return
			}
		}()
	}

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
