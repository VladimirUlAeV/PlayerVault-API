package handler

import (
	dataa "PlayerVaultAPI/internal/Data"
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
	"strconv"
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
type Handler struct{
	storage dataa.PlayerStorage
}

func NewHandler(i dataa.PlayerStorage) *Handler{
	return &Handler{
		storage: i,
	}
}

func (h *Handler)saveHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println ("Получен запрос на регестрацию нового игрока")
	if r.Method != http.MethodPost{
		http.Error(w,error_project.ErrorNetWriteLog(nil, error_project.ErrorMethod),http.StatusBadRequest)
		return
	}

	userType := r.Header.Get("Content-type") //проверка контента
	if !strings.HasPrefix(userType, "application/json"){
		http.Error(
			w,
			error_project.ErrorNetWriteLog(nil, error_project.ErrorRequstType), 
			http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var data dataa.Player
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
	}

	if err := h.storage.Save(data); err != nil{
		http.Error(
			w,
			error_project.ErrorNetWriteLog(nil,err),
			http.StatusBadRequest)
		return
	}

	
	w.Header().Set("Content-Type","applictaion/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))


}

func (h *Handler)getHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Получен запрос на выдачу данных по айди")
	if r.Method != http.MethodGet{
		http.Error(w,error_project.ErrorNetWriteLog(nil, error_project.ErrorMethod),http.StatusBadRequest)
		return
	}

	BodyBytes, err := io.ReadAll(r.Body)
	if err != nil{
		http.Error(w, error_project.ErrorNetWriteLog(err, error_project.ErrorIOread),http.StatusBadRequest)
	}
	defer r.Body.Close()

	rawText := strings.TrimSpace(string(BodyBytes))
	if rawText == ""{
		http.Error(w, error_project.ErrorNetWriteLog(nil, error_project.ErrorJsonEmptyBody), http.StatusBadRequest)
		return
	}

	IdPlayer, err := strconv.Atoi(rawText)
	if err != nil{
		http.Error(w, error_project.ErrorNetWriteLog(err, error_project.ErrorIOread),http.StatusBadRequest)
		return
	}

	Playere, err := h.storage.Get(IdPlayer)
	if err != nil{
		http.Error(w, error_project.ErrorNetWriteLog(nil, err), http.StatusBadRequest)
		return
	}

	//отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Playere); err != nil{
		http.Error(w, error_project.ErrorNetWriteLog(nil, err), http.StatusBadRequest)
		return
	}
}

func (h *Handler)getAllHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Получен запрос на выдачу данных всех игроков")
	if r.Method != http.MethodGet{
		http.Error(w,error_project.ErrorNetWriteLog(nil, error_project.ErrorMethod),http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	AllPlayer, err := h.storage.GetAll()
	if err != nil{
		http.Error(w, error_project.ErrorNetWriteLog(nil, err), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(AllPlayer); err != nil{
		http.Error(w, error_project.ErrorNetWriteLog(nil, err), http.StatusBadRequest)
		return
	}

}

func StartServer(h *Handler){

	var mux = http.NewServeMux()
	var server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server start on 8080 port, для остановки сервера используйте комбинацию ctrl+c")

	mux.HandleFunc("/NewPlayer", h.saveHandler)
	mux.HandleFunc("/Player", h.getHandler)
	mux.HandleFunc("/AllPlayer", h.getAllHandler)

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