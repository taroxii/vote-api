package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/taroxii/vote-api/pkg/config"
	httpDelivery "github.com/taroxii/vote-api/pkg/interface/echo_server"
	_itemRepo "github.com/taroxii/vote-api/pkg/repository/item_repository/postgres"
	_userRepo "github.com/taroxii/vote-api/pkg/repository/user_repository/postgres"
	usecase "github.com/taroxii/vote-api/pkg/usecase"
	"github.com/taroxii/vote-api/pkg/utils/logger"
)

func init() {
	fmt.Printf("go init %s", "...")
	logger.InitializeLogger()
	config.NewConfig()
}

func main() {

	dsnPassword := ""
	sslmode := "disable"
	if config.Config.Postgres.Password != "" {
		dsnPassword = fmt.Sprintf("password=%s", config.Config.Postgres.Password)
	}

	dsn := fmt.Sprintf("host=%s user=%s %s dbname=%s port=%s sslmode=%s ", config.Config.Postgres.Host, config.Config.Postgres.Username, dsnPassword, config.Config.Postgres.Database, config.Config.Postgres.Port, sslmode)

	dbConn, err := sql.Open(`postgres`, dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	timeoutContext := time.Duration(config.Config.Timeout) * time.Second

	itemRepo := _itemRepo.NewPostgresItemRepository(dbConn)
	userRepo := _userRepo.NewPostgresUserRepository(dbConn)
	itemUsecase := usecase.NewItemUsecase(itemRepo, userRepo, timeoutContext)
	httpDelivery.NewItemHandler(e, itemUsecase)

	userUsecase := usecase.NewUserUsecase(userRepo, timeoutContext)
	httpDelivery.NewUserHandler(e, userUsecase)
	log.Fatal(e.Start(fmt.Sprintf(":%s", config.Config.ServerPort)))
}
