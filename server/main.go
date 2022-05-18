package main

import (
	"Blog/ent"
	"context"
	"embed"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/acme/autocert"
)

var (
	//go:embed html/*
	fs embed.FS
)

type viewEvent struct {
	page            string
	ipAddress       string
	TimestampMillis time.Time
}

var viewEventsChan chan viewEvent
var client *ent.Client

func main() {

	var err error
	// ambitiously large view event channel
	viewEventsChan = make(chan viewEvent, 64)
	go eventWorker()

	// DB create and auto migration.
	client, err = ent.Open("sqlite3", "file:blog?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// attempt to run migration
	if err := client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	var contentHandler = echo.WrapHandler(http.FileServer(http.FS(fs)))
	var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/html/$1"}) // map the html folder to
	e.GET("/*", contentHandler, contentRewrite)
	e.GET("/view/:id", metrics)

	e.AutoTLSManager.Cache = autocert.DirCache("/home/.cache")
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}

func metrics(c echo.Context) error {

	pathId := c.Param("id")
	page, err := base64.StdEncoding.DecodeString(pathId)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	userIP := net.ParseIP(ip)

	viewEventsChan <- viewEvent{
		page:            string(page),
		ipAddress:       userIP.String(),
		TimestampMillis: time.Now().UTC(),
	}

	// log view event
	return c.NoContent(http.StatusOK)
}

// Creates worker to write view events to parquet file.
func eventWorker() {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Rrecovered from panic", r)
		}
	}()

	for {
		event := <-viewEventsChan
		_, err := client.ViewEvent.
			Create().
			SetPage(event.page).
			SetIPAddress(event.ipAddress).
			SetEventTime(event.TimestampMillis).
			Save(context.Background())

		if err != nil {
			log.Error(err)
		}
	}

}
