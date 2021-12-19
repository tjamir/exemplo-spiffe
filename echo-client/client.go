package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

type ClientConfig struct {
	ServerPort    int    `json:"server_port"`
	ServerAddress string `json:"server_address"`
	AuthorizedId  string `json:"authorized_id"`
	SocketPath    string `json:"socket_path"`
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

func loadConfig() (ClientConfig, error) {
	clientConfig := ClientConfig{}
	configFile, err := os.Open("conf/echo-client.json")
	if err != nil {
		return clientConfig, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&clientConfig)
	if err != nil {
		return clientConfig, err
	}

	return clientConfig, nil
}

func RunServer(ctx context.Context, config ClientConfig) error {
	allowedId, err := spiffeid.FromString(config.AuthorizedId)
	if err != nil {
		return err
	}
	address := fmt.Sprintf("%s:%d", config.ServerAddress, config.ServerPort)
	log.Println("Establecendo conexão tls com o servidor", address)

	dial, err := spiffetls.DialWithMode(ctx,
		"tcp", address,
		spiffetls.MTLSClientWithSourceOptions(tlsconfig.AuthorizeID(allowedId),
			workloadapi.WithClientOptions(workloadapi.WithAddr(config.SocketPath))))
	if err != nil {
		log.Println("Erro ao estabelecer conexão tls com o servidor", address)
		return err
	}
	defer dial.Close()

	log.Println("Conexão estabelecida, enviando mensagens", address)

	var message string
	fmt.Println("Digite uma mensagem:")
	_, err = fmt.Scanln(&message)
	if err != nil {
		log.Default().Println("Erro lendo mensagem", err)
	}
	err = sendMessage(dial, ctx, config, message, allowedId)
	if err != nil {
		log.Default().Println("Erro enviando mensagem", err)
	}
	return nil

}

func sendMessage(dial net.Conn, ctx context.Context, config ClientConfig, message string, allowedId spiffeid.ID) error {

	buffWriter := bufio.NewWriter(dial)
	_, err := buffWriter.WriteString(message + "\n")
	if err != nil {
		return err
	}
	err = buffWriter.Flush()
	if err != nil {
		return err
	}
	log.Println("Mensagem enviada, aguardando resposta")
	buffReader := bufio.NewReader(dial)
	message, err = buffReader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Println("Resposta recebida:", message)
	return nil
}
