package handler

import (
	error_project "PlayerVaultAPI/internal/Error"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var mux = http.NewServeMux()
var server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
}

func latencyStop(){
	fmt.Print("Остановка сервера")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print("..")
	time.Sleep(1 * time.Second)
	fmt.Print(".100%")
	fmt.Println("")
}

func testHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println ("Получен и обработан запрос")
}

func StartServer(){
	fmt.Println("Server start on 8080 port")

	mux.HandleFunc("/", testHandler)

	go func(){
		if err := server.ListenAndServe(); err != nil{
			totalerr := errors.New(error_project.ErrorStartServer.Error() + err.Error())
			error_project.ErrorWriteLog(totalerr)
	}
	}()

	//Улавливать килл 
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	latencyStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil{
		error_project.ErrorWriteLog(err)
	}
}