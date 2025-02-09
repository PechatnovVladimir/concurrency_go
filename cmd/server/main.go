package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"kvdatabase/internal/config"
	"kvdatabase/internal/initialization"
	"log"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	ConfigFileName := os.Getenv("CONFIG_FILE_NAME")

	cfg := &config.Config{}
	if ConfigFileName != "" {
		data, err := os.ReadFile(ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err = config.Load(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	i, err := initialization.NewInit(cfg)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = i.StartApp(ctx)
	if err != nil {
		fmt.Println("11")
		log.Fatal(err)
	}

}
