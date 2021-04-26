package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aep/parted/src"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("could not load env: ", err)
	}

	err = checkRequiredEnv(map[string]struct{}{
		"API_KEY": {},
	})
	if err != nil {
		log.Fatalln(err)
	}

	src.Main()
}

func checkRequiredEnv(req map[string]struct{}) error {
	for key := range req {
		if _, ok := os.LookupEnv(key); !ok {
			return fmt.Errorf("missing required env variable, %s", key)
		}
	}
	return nil
}
