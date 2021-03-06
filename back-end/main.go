package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"schedule/api"
	"schedule/conf"
	"schedule/dbb"
	"schedule/middleware"
	"schedule/snowid"
	"schedule/wechatid"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func main() {
	//dbAddr := flag.String("dbaddr", "127.0.0.1", "database addr")
	//dbUser := flag.String("dbuser", "root", "database user")
	//dbPassword := flag.String("dbpassword", "root", "database password")
	//listenAddr := flag.String("listen", "127.0.0.1:8080", "web listen addr")
	//debug := flag.Bool("debug", false, "debug mode")

	//flag.Parse()

	url, user, password, listenAddr, debug := conf.ReadConf()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/schedule", user, password, url))
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%+v", err)
	}

	dbInstance := dbb.InitDB(db)
	db.SetConnMaxLifetime(time.Second * 14400)

	if !(debug) {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	if debug {
		router.Use(gin.ErrorLogger())
	}

	// init session store
	sessionAuthKey := make([]byte, 32)
	if _, err := rand.Read(sessionAuthKey); err != nil {
		log.Fatalf("%+v", xerrors.Errorf("generate session auth key failed: %w", err))
	}
	store := memstore.NewStore(sessionAuthKey)

	router.Use(sessions.Sessions("cookie", store))

	apI := api.New(dbInstance, wechatid.Wechat{}, snowid.SnowFlake{})

	router.POST("/login", apI.Login)
	router.POST("/savePerson", middleware.SessionChecker(), apI.SavePerson)

	router.POST("/newRota", middleware.SessionChecker(), apI.NewRota)
	router.GET("/rotas", middleware.SessionChecker(), apI.GetRotas)
	router.GET("/join", middleware.SessionChecker(), apI.GetJoin)

	router.GET("/rota/:rotaId", middleware.SessionChecker(), apI.GetFrees)
	router.POST("/chooseFree/:rotaId", middleware.SessionChecker(), apI.ChooseFree)

	router.GET("/generate/:rotaId", apI.Generate)
	router.GET("/download/:rotaId", apI.Download)

	router.DELETE("/delete/:rotaId", middleware.SessionChecker(), apI.DeleteRota)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) // Interrupt Signal = syscall.SIGINT

	httpServer := http.Server{Handler: router, Addr: listenAddr}

	shutdownChan := make(chan struct{})

	go func() {
		<-signalChan
		timeout, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelFunc()
		if err := httpServer.Shutdown(timeout); err != nil {
			log.Println(err)
		}
		close(shutdownChan)
	}()

	log.Println("start http server")

	err = httpServer.ListenAndServe()
	switch err {
	case http.ErrServerClosed:
		<-shutdownChan

	default:
		log.Println(err)
	}

	if err := db.Close(); err != nil {
		log.Printf("%+v", xerrors.Errorf("close db failed: %w", err))
	}
}
