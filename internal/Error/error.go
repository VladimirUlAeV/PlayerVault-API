package error_project

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrorGetUser = errors.New("Ошибка в /data: Блок get не отработал")
	ErrorStartServer = errors.New("Ошибка в /handler: ошибка запуска по причине: ")

)

func ErrorWriteLog(errinput error){
	data := []byte(errinput.Error())
	if err := os.WriteFile("error.log",data, 0644); err != nil{
		fmt.Print(err)
	}
}