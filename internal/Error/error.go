package error_project

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	ErrorGetUser = errors.New("Нету игрока с таким айди")
	ErrorStartServer = errors.New("ошибка работы сервера по причине: ")
	ErrorMethod = errors.New("выбран не верный метод")
	ErrorJsonSyntax = errors.New("ошибка в синтаксисе запроса:")
	ErrorJsonTypeDate = errors.New("в запросе есть не правильные данные:")
	ErrorJsonEmptyBody = errors.New("пустое тело запроса:")
	ErrorJsonIncorrectRequst = errors.New("Неккоректный запрос:")
	ErrorRequstType = errors.New("Неправильный Content-Type: ")
	ErrorNicknameEmpty = errors.New(":Пустое имя пользователя")
	ErrorScoreInput = errors.New("Счет не может быть отрицательным")
	ErrorIdValid = errors.New("Игрок с таким айди уже есть")
	ErrorIOread = errors.New("Неудачный формат данных:")
	ErrorEmptyMap = errors.New("В БД нету игроков")
)

func DualError(err error, staterr error)error{
	return errors.New(staterr.Error() + err.Error())
}

func ErrorWriteLog(errinput error){
	file, err := os.OpenFile("error.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{}
	defer file.Close()
	logWrite := fmt.Sprintf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"),"Ошибка", errinput.Error())
	if _, err := file.WriteString(logWrite); err != nil{}
}

func ErrorNetWriteLog(httpErr error, statErr error)string{
	if httpErr == nil{
		file, err := os.OpenFile("error.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil{}
		defer file.Close()
		logWrite := fmt.Sprintf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"),"Ошибка", statErr.Error())
		if _, err := file.WriteString(logWrite); err != nil{}
		return statErr.Error()
	}
	FullError := statErr.Error() + httpErr.Error()
	file, err := os.OpenFile("error.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{}
	defer file.Close()
	logWrite := fmt.Sprintf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"),"Ошибка", FullError)
	if _, err := file.WriteString(logWrite); err != nil{}
	return FullError
}