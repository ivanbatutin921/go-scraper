package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

var urls = []string{
	"https://jsonplaceholder.typicode.com/todos?_limit=2",
	"https://jsonplaceholder.typicode.com/posts?_limit=2",
	"https://jsonplaceholder.typicode.com/users?_limit=2",
}

func Parse(urls []string) {
	chParse := make(chan string)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			chParse <- string(body)
			wg.Done()
		}(url)
	}

	go func() {
		wg.Wait()
		close(chParse)
	}()

	writeToFile(chParse)
}

func writeToFile(ch chan string) {
	file, err := os.Create("tmp.txt")
	if err != nil {
		fmt.Println("не удалось создать файл")
	}

	defer file.Close()

	for data := range ch {
		file.WriteString(string(data))
	}
	fmt.Println("данные в файл записаны")

}

func main() {

	log.Println("начало...")
	start := time.Now()
	Parse(urls)
	log.Println("время:", time.Since(start)) // 100-200 ms
}
