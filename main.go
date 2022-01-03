package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	vault "github.com/hashicorp/vault/api"
)

var client *vault.Client

func main() {
	setupVault()
	startServer()
}

func setupVault() {
	config := &vault.Config{
		Address:      "http://127.0.0.1:8200",
		HttpClient:   cleanhttp.DefaultPooledClient(),
		Timeout:      time.Second * 60,
		MinRetryWait: time.Millisecond * 1000,
		MaxRetryWait: time.Millisecond * 1500,
		MaxRetries:   2,
		Backoff:      retryablehttp.LinearJitterBackoff,
	}
	var err error
	client, err = vault.NewClient(config)
	if len(os.Args) < 2 {
		log.Panicln("Usage: vault_api.exe <token>")
	}
	client.SetToken(os.Args[1])
	if err != nil {
		fmt.Println("unable to initialize Vault client")
	}
	log.Println("Vault setup complete.")
}
