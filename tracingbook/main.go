package main

import (
	"fmt"
	"gohome_self/tracingbook/controllers"
	"gohome_self/tracingbook/crawler"
	"gohome_self/tracingbook/infrastructure"
	"gohome_self/tracingbook/mail"
	"gohome_self/tracingbook/middlewares"
	"gohome_self/tracingbook/models"
	"gohome_self/tracingbook/seeds"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func drop(db *gorm.DB) {
	db.DropTableIfExists(
		&models.FileUpload{},
		&models.Book{}, &models.Site{}, &models.UpdateItem{},
		&models.UserRole{}, &models.Role{}, &models.User{})
}

func migrate(database *gorm.DB) {

	database.AutoMigrate(&models.Role{})
	database.AutoMigrate(&models.UserRole{})
	database.AutoMigrate(&models.User{})

	database.AutoMigrate(&models.Book{})
	database.AutoMigrate(&models.Site{})
	database.AutoMigrate(&models.UpdateItem{})

	//database.AutoMigrate(&models.FileUpload{})
}

func addDbConstraints(database *gorm.DB) {
	// TODO: it is well known GORM does not add foreign keys even after using ForeignKey in struct, but, why manually does not work neither ?

	dialect := database.Dialect().GetName() // mysql, sqlite3
	if dialect != "sqlite3" {

		database.Model(&models.UserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.UserRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

	} else if dialect == "sqlite3" {

	}

	database.Model(&models.UserRole{}).AddIndex("user_roles__idx_user_id", "user_id")
}
func create(database *gorm.DB) {
	drop(database)
	migrate(database)
	addDbConstraints(database)
}

func main() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	println("DB_DIALECT:", os.Getenv("DB_DIALECT"))
	println("DB_DRIVER:", os.Getenv("DB_DRIVER"))

	database := infrastructure.OpenDbConnection()

	defer database.Close()

	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}

		if first == "create" {
			create(database)
		} else if first == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if second == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}

	migrate(database)

	// gin.New() - new gin Instance with no middlewares
	// goGonicEngine.Use(gin.Logger())
	// goGonicEngine.Use(gin.Recovery())
	goGonicEngine := gin.Default() // gin with the Logger and Recovery Middlewares attached
	// Allow all Origins
	goGonicEngine.Use(cors.Default())

	goGonicEngine.Use(middlewares.Benchmark())

	// goGonicEngine.Use(middlewares.Cors())

	goGonicEngine.Use(middlewares.UserLoaderMiddleware())
	goGonicEngine.Static("/static", "./static")
	apiRouteGroup := goGonicEngine.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterPageRoutes(apiRouteGroup.Group("/"))

	mail.NML.Pass = os.Getenv("MIXIU_PASS")
	crawler.InitDingDian()
	crawler.FetchBooks()

	goGonicEngine.Run(":8080") // listen and serve on 0.0.0.0:8080

}
