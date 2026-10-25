package go_db_entities

import (
	"io"
	"log"
)

// CloseOrLog closes the given Closer and logs error, if one occurs.
func CloseOrLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("Close error: %v\n", err)
	}
}
