package main

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func main() {
	var w sync.WaitGroup
	for i := 0; i < 10; i++ {
		w.Add(1)
		go SendEvent(&w)
	}
	w.Wait()
}

func SendEvent(w *sync.WaitGroup) {
	defer w.Done()
	url := "http://localhost:8080/publish"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{
    "userid": "%s",
	"payload": "%s"
}`, uuid.NewString(), uuid.NewString()))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
