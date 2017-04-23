package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func worker(id int, jobs <-chan string, signals <-chan bool, result chan<- bool) {
	for {
		select {
		case site := <-jobs:
			if len(site) > 1 {
				fmt.Println("Worker: ", id, "Site: ", site)
				resp, err := http.Get("http://" + site)
				if (err == nil) && (resp.StatusCode == 200) {
					fmt.Println("Worker: ", id, "Site: ", site, resp.StatusCode)
					result <- true
				} else {
					log.Println(err)
				}
			}
		case sig := <-signals:
			if sig {
				break
			}
		default:
		}
	}
}

func main() {
	var summ, status200 uint64 = 0, 0
	signals := make(chan bool)
	jobs := make(chan string, 10)
	result := make(chan bool)

	for w := 1; w <= 10; w++ {
		go worker(w, jobs, signals, result)
	}

	file, err := os.Open("topmillion.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(file)
	for s.Scan() {

		site := strings.Fields(s.Text())[1]
		summ++
		jobs <- site

	}
	close(jobs)
	for range result {
		atomic.AddUint64(&status200, 1)
	}

	signals <- true
	close(result)
	time.Sleep(time.Second)
	fmt.Println("summ: ", summ, "Status 200:", atomic.LoadUint64(&status200))
}
