package handlers

import (
	"fmt"
	"io"
	"monitoring/model"
	"net/http"
)

func RecvStream(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")
	file, _, err := r.FormFile("frame")
	if err != nil {
		fmt.Println("frame get error")
		return
	}
	defer file.Close()

	data, _ := io.ReadAll(file)
	model.Mu.Lock()
	model.LatestFrame[ip] = data
	model.Mu.Unlock()
	// fmt.Println("stream start")
	w.Write([]byte("ok"))
}
