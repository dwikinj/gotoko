package app

import (
	"flag"
	"fmt"
	"github.com/dwikinj/gotoko/app/database/seeders"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	//server.InitializeDB(dbConfig)
	server.InitializeRoutes()
	//seeders.DBSeed(server.DB)

}

func (server *Server) Run(addr string) {
	fmt.Println("Server running at port " + addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) InitializeDB(dbConfig DBConfig) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}
	server.DB = db
	fmt.Println("Success connect to DB")

}

func (server *Server) dbMigrate() {
	for _, model := range RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)
		if err != nil {
			log.Fatalf("failed to auto migrate. error : %v", err.Error())
		}
	}
}

func (server *Server) initCommands(config AppConfig, dbConfig DBConfig) {
	server.InitializeDB(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []*cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		}, {
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && len(value) != 0 {
		return value
	} else {
		return fallback
	}
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load env")
	}
	appConfig.AppName = getEnv("APP_NAME", "GoToko")
	appConfig.AppEnv = getEnv("APP_ENV", "Dev")
	appConfig.AppPort = getEnv("APP_PORT", "8080")

	//DB config
	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "postgres")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "1234")
	dbConfig.DBName = getEnv("DB_NAME", "gotoko")
	dbConfig.DBPort = getEnv("DB_PORT", "5432")

	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		server.initCommands(appConfig, dbConfig)
	} else {
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}
