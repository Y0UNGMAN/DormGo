package main

import (
	"fmt"

	"github.com/Y0UNGMAN/DormGo/backend/utils"
)

func main() {
	token, err := utils.GenToken(9, "tester")
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
}
