package bale

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type manager struct {
	url    string
	apiKey string
	botID  int
	client *http.Client
}

func New() *manager {

	botID, err := strconv.Atoi(os.Getenv("BALE_BOT_ID"))
	if err != nil {
		log.Println(fmt.Errorf("Bale, new bale manager: invalid id, %w", err))
		return nil
	}
	if botID <= 0 {
		log.Println(fmt.Errorf("Bale, new bale manager: invalid id.	"))
		return nil
	}

	url := os.Getenv("BALE_URL")
	if url == "" {
		log.Println(fmt.Errorf("Bale, new bale manager: invalid url."))
		return nil
	}

	apiKey := os.Getenv("BALE_API_KEY")
	if apiKey == "" {
		log.Println(fmt.Errorf("Bale, new bale manager: invalid api key."))
		return nil
	}

	client := &http.Client{Timeout: 10 * time.Second}

	return &manager{
		url:    url,
		botID:  botID,
		apiKey: apiKey,
		client: client,
	}
}
