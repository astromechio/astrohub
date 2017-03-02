package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bytes"
	"fmt"
)

type jobReq struct {
	ID   string `json:"id"`
	Body []byte `json:"originalRequest"`
}

type question struct {
	First  int `json:"first"`
	Second int `json:"second"`
}

type addRes struct {
	Answer int `json:"answer"`
}

func main() {
	next := make(chan bool)

	for true {
		go handleMath(next)

		ok := <-next
		if !ok {
			break
		}
	}
}

func handleMath(next chan (bool)) {
	req, _ := http.NewRequest("GET", "http://127.0.0.1:3000/service/add/jobs", nil)
	res, _ := http.DefaultClient.Do(req)

	next <- true

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)

	var job jobReq
	_ = json.Unmarshal(resBody, &job)

	fmt.Println(string(job.Body))

	var ques question
	_ = json.Unmarshal(job.Body, &ques)

	answer := addRes{
		Answer: ques.First + ques.Second,
	}

	answerJSON, _ := json.Marshal(answer)

	answerReq, _ := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:3000/jobs/%s/response", job.ID), bytes.NewReader(answerJSON))
	answerRes, _ := http.DefaultClient.Do(answerReq)

	fmt.Println(answerRes.Status)
}
