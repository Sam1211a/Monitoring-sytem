package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Device struct {
	IP       string `json:"ip"`
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

func CheckHost(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":80", 200*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func Hostname(ip string) string {
	name, err := net.LookupAddr(ip)
	if err != nil || len(name) == 0 {
		return "Unknown"
	}
	return name[0]
}

// var tmp = template.Must(template.ParseFiles("templates/dashboard.html"))

func ScanIp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	base := "192.168.1."
	var wg sync.WaitGroup
	results := make(chan Device)
	for i := 1; i <= 255; i++ {
		wg.Add(1)
		ip := fmt.Sprintf("%s%d", base, i)

		status := "DOWN"
		os := "unknown"
		go func(ip string) {
			defer wg.Done()
			host := Hostname(ip)
			if CheckHost(ip) {
				status = "UP"
				os = "Windows"
			}

			results <- Device{
				IP:       ip,
				Status:   status,
				Hostname: host,
				OS:       os,
			}
		}(ip)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		jsonData, _ := json.Marshal(res)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()

	}

	fmt.Fprintf(w, "data: {\"done\": true}\n\n")
	flusher.Flush()
}
