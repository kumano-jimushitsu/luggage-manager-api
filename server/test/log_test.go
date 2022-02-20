package test

import (
	"log"
	"os"
	"testing"
)

func TestSaveLogToFile(t *testing.T) {
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
					log.Println(err)
	}
	defer f.Close()
	
	logger := log.New(f, "", log.LstdFlags)
	logger.Println("text to append")
	logger.Println("more text to append")
}
