package main

import (
	"github.com/Y0UNGMAN/DormGo/backend/model"
	"github.com/Y0UNGMAN/DormGo/backend/router"
)

func main() {
	model.Database()

	r := router.App()
	r.Run(":8080")
}
