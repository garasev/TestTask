package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	n := 0
	i := 0
	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		n, _ = strconv.Atoi(scanner.Text())

		// Проверьте ошибку сканирования (например, конец файла)
		if scanner.Err() != nil {
			fmt.Println("Ошибка чтения ввода:", scanner.Err())
			break // Выход из цикла при ошибке
		}

		if n == 0 {
			t := time.NewTicker(time.Millisecond * 100)
			for {
				select {
				case <-t.C:
					i += 1
					req(client, i)
				}
			}
		}

		for i := 0; i < n; i++ {
			req(client, i)
		}
	}
}

func req(client *http.Client, i int) {
	req, err := http.NewRequest("GET", "http://localhost:8090/api/", nil)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
		return
	}

	fmt.Printf("Запрос #%d: Статус ответа: %s\n", i+1, resp.Status)

	resp.Body.Close()
}
