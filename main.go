package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime"
)

type Request struct {
	Input []rune `json:"input"`
}

type Response struct {
	Output string `json:"output"`
}

func solve(input []rune) string {
	// len(string) returns the number of bytes.
	// This is, of course, a memory waste, but I don't care
	output := make([]rune, len(input))
	currentPos := 0
	for _, ch := range input {
		switch ch {
		case '#':
			if currentPos > 0 {
				currentPos--
			}
		default:
			output[currentPos] = ch
			currentPos++
		}
	}
	return string(output[0:currentPos])
}

func handle(writer http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	if r.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body := bytes.Buffer{}
	_, err := io.Copy(&body, r.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	var request Request
	err = json.Unmarshal(body.Bytes(), &request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	solution := solve(request.Input)
	marshal, err := json.Marshal(&Response{Output: solution})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write(marshal)
}

func main() {
	//println(solve("abc#de#語"))
	runtime.GOMAXPROCS(1)
	http.HandleFunc("/solve", handle)
	http.ListenAndServe(":8080", nil)
}
