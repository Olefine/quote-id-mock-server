package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olefine/quote-id-mock/domain"
	"github.com/olefine/quote-id-mock/handlers"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

const schemaTable = "./store/schema.sql"
const dbName = "./store/token.db"

var port = flag.Int64("port", 8080, "Specify port to run server on.")
var env = flag.String("env", "development", "Specify environment to run on")
var db *sql.DB

func init() {
	var err error

	flag.Parse()
	db, err = sql.Open("sqlite3", dbName)
	checkErr(err)

	file, _ := os.Open(schemaTable)
	sqlContext, err := ioutil.ReadAll(file)
	checkErr(err)

	_, err = db.Exec(string(sqlContext))
	checkErr(err)

	log.Printf("Successfully connected to DB - %s\n", dbName)
}

func dbMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("tokenRepository", domain.NewTokenRepository(db))
	}
}

func loggingMiddleware(w *domain.HTTP2Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("loggingWriter", w)
	}
}

func main() {
	urlToBind := fmt.Sprintf(":%d", *port)
	broadcastTo := domain.HTTP2Writer{Out: make(chan interface{})}
	gin.DefaultWriter = domain.NewLogMultiplexor(broadcastTo)
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(dbMiddleware())

	router.POST("/token", handlers.HandleAuth)
	router.POST("/counts", handlers.HandleCounts)
	router.GET("/logs", loggingMiddleware(&broadcastTo), handlers.HandleLogs)

	router.Run(urlToBind)
}
