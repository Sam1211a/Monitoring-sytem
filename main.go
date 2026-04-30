package main

import (
	"fmt"
	"monitoring/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.Dashboard)
	http.HandleFunc("/scan", handlers.ScanIp)
	http.HandleFunc("/upload", handlers.UploadImg)
	http.Handle("/screens/", http.StripPrefix("/screens/", http.FileServer(http.Dir("screens"))))
	fmt.Println(`Server running at http://localhost:8080`)
	http.ListenAndServe(":8080", nil)
}
