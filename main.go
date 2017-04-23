package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var status200 uint64 = 0
	var summ uint64 = 0
	var wg sync.WaitGroup
	messages := make(chan bool)
	signals := make(chan bool)

	go func() {
		var counter uint64
		counter = 0
		for {
			select {
			case msg := <-messages:
				if msg {
					counter++
				}
			case <-signals:
				atomic.AddUint64(&status200, counter)
				break
			default:

			}
		}
	}()

	file, err := os.Open("topmillion.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(file)
	for s.Scan() {
		site := strings.Fields(s.Text())[1]
		wg.Add(1)
		summ++
		go func(site string) {
			defer wg.Done()
			resp, err := http.Get("http://" + site)
			if (err == nil) && (resp.StatusCode == 200) {
				messages <- true
			}
		}(site)
	}
	wg.Wait()
	signals <- true
	time.Sleep(time.Second)
	fmt.Println("summ: ", summ, "Status 200:", atomic.LoadUint64(&status200))
}
