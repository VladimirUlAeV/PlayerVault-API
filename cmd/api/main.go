package main

import (
	"PlayerVaultAPI/internal/handler"
	"bufio"
	"fmt"
	"os"
	"time"
)

func menu(){
	fmt.Println("Выберите команду:")
	fmt.Println("1 - запустить сервер")
	fmt.Println("2 - открыть логи")
	fmt.Println("0 - выйти из программы")
}

func latencyStart(){
	fmt.Print("CLI start.")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print("..")
	time.Sleep(1 * time.Second)
	fmt.Print(".100%")
	fmt.Println("")
}

func main(){
	latencyStart()
	scan := bufio.NewScanner(os.Stdin)
	for{
		fmt.Println("")
		menu()
		scan.Scan()
		inputUser := scan.Text()
		switch inputUser{
			case "1":
				handler.StartServer()
				
			case "2":
				//Открытие лог файла	
			case "0":
				return
			default:
				fmt.Println("Команда не найдена")
		}
	}
}