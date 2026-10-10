package handler

import (
	data "PlayerVaultAPI/internal/Data"
	error_project "PlayerVaultAPI/internal/Error"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

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

func saveHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println ("Получен запрос на регестрацию нового игрока")
	if r.Method != http.MethodPost{
		error_project.ErrorWriteLog(error_project.ErrorMethod)
		return
	}

	userType := r.Header.Get("Content-type")
	if !strings.HasPrefix(userType, "application/json"){
		http.Error(
			w,
			error_project.ErrorNetWriteLog(nil, error_project.ErrorRequstType), 
			http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var data data.Player
	unmarsh := json.NewDecoder(r.Body)
	unmarsh.DisallowUnknownFields() //ошибка при лишних полях

	if err := unmarsh.Decode(&data); err != nil{
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		switch {
			case errors.As(err, &syntaxError):
				http.Error(
					w,
					error_project.ErrorNetWriteLog(err,error_project.ErrorJsonSyntax),
				http.StatusBadRequest)
			case errors.As(err, &unmarshallTypeError):
				http.Error(
					w,
					error_project.ErrorNetWriteLog(err, error_project.ErrorJsonTypeDate),
					http.StatusBadRequest)
			case errors.Is(err, io.EOF): //пустое тело
				http.Error(
					w,
					error_project.ErrorNetWriteLog(err, error_project.ErrorJsonEmptyBody),
					http.StatusBadRequest)
			case strings.HasPrefix(err.Error(), "json: unknown field"):
				//error при лишних неизвестных полях (при наличий DisallowUnknownFields)
			default:
				http.Error(w, 
				error_project.ErrorNetWriteLog(err, error_project.ErrorJsonIncorrectRequst),
				http.StatusBadRequest)
			}
		return
	}
	if data.Nickname == ""{
		http.Error(
			w,
			error_project.ErrorNetWriteLog(nil,error_project.ErrorNicknameEmpty),
			http.StatusBadRequest)
		return
	}
	if data.Score < 0 {
		http.Error(
			w,
			error_project.ErrorNetWriteLog(nil,error_project.ErrorScoreInput),
			http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type","applictaion/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))


}

func StartServer(){
	var mux = http.NewServeMux()
	var server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server start on 8080 port, для остановки сервера используйте комбинацию ctrl+c")

	mux.HandleFunc("/NewPlayer", saveHandler)

	go func(){
		if err := server.ListenAndServe(); err != nil{
			if err.Error() == "http: Server closed"{}else{
			error_project.ErrorWriteLog(error_project.DualError(err, error_project.ErrorStartServer))
			}
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