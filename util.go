package go_db_entities

import (
	"io"
	"log"
)

func CloseOrLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("Close error: %v\n", err)
	}
}
