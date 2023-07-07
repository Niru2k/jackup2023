package drivers

import (
	//User-defined packages
	"blog/helper"
	"blog/repository"

	//Inbuild packages
	"fmt"

	//Third-party packages
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {
	//create a connection to ginframework database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", helper.Host, helper.Port, helper.User, helper.Password, helper.Dbname)
	Db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Established a successful connection to '%s' database!!!\n", helper.Dbname)
	repository.TableCreation(Db)
	return Db
}
