package main

import (
	"fmt"
	"trab02/service01/database"
	products_repository "trab02/service03/repositories"
)

func main() {
	db, _ := database.InitMySqlConn()
	repo := products_repository.NewProductRepository(db)
	fmt.Println(repo.GetProducts())
}
