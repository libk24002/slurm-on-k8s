package utils

import (
	"log"
)

func HandleException(err error) {
	if err != nil {
		log.Printf("Exception: %v", err)
	}
}
