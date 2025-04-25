package main

import (
	"fmt"

	"github.com/AidanRJ1/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		errorMsg := fmt.Errorf("error occured while reading config file: %v", err)
		fmt.Println(errorMsg)
	}
	fmt.Println(cfg)

	err = cfg.SetUser("Aidan")
	if err != nil {
		errorMsg := fmt.Errorf("error occured while writing to config: %v", err)
		fmt.Println(errorMsg)
	} else {
		fmt.Println("Successfuly written to file")
	}
}