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

	err = checkRequiredEnv([]string{
		"API_KEY",
	})
	if err != nil {
		log.Fatalln(err)
	}

	src.ListenAndServe()
}

func checkRequiredEnv(required []string) error {
	for _, req := range required {
		if _, ok := os.LookupEnv(req); !ok {
			return fmt.Errorf("missing required env variable, %s", req)
		}
	}
	return nil
}
