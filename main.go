package main

import (
	"fmt"

	"github.com/AidanRJ1/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		errorMsg := fmt.Errorf("error occured while reading config file: %v", err)
		fmt.Println(errorMsg)
	}

	fmt.Println(config)
}