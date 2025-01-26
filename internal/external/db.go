package external

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaDB(host, user, password, dbname, port string) {
	conexao := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(conexao))

	if err != nil {
		log.Panic("Erro na conexao com banco de dados")
	}

	DB = db
}
