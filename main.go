package main

import (
	"fmt"

	"github.com/jahidul39306/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	fmt.Println("DB URL:", cfg.DBURL)
	fmt.Println("Current User Name:", cfg.CurrentUserName)

	err = cfg.SetUser("Jahidul")
	if err != nil {
		fmt.Println("Error setting user:", err)
		return
	}
	fmt.Println("Updated Current User Name:", cfg.CurrentUserName)
}
