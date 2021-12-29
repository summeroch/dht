package basic

import (
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Printf("%s\n", err)
	}
}
