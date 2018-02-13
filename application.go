package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olefine/quote-id-mock/domain"
	"github.com/olefine/quote-id-mock/handlers"
	"golang.org/x/net/websocket"
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
var broadcastTo domain.HTTP2Writer
var broadcaster = domain.NewBroadcaster()

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

// WSHandler handles logging collection
func WSHandler(ws *websocket.Conn) {
	for {
		select {
		case msg := <-broadcastTo.Out:
			broadcaster.Join(ws)
			broadcaster.Broadcast(msg.([]byte))
		}
	}
}

func main() {
	urlToBind := fmt.Sprintf(":%d", *port)

	broadcastTo = domain.HTTP2Writer{Out: make(chan interface{})}
	logMultiplexor := domain.NewLogMultiplexor(broadcastTo)

	gin.DefaultWriter = logMultiplexor
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Use(domain.RequestLogger())
	router.Use(dbMiddleware())
	router.Use(gin.Recovery())

	router.POST("/token", handlers.HandleAuth)
	router.POST("/counts", handlers.HandleCounts)

	router.GET("/ws", func(c *gin.Context) {
		handler := websocket.Handler(WSHandler)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/logs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	router.Run(urlToBind)
}
