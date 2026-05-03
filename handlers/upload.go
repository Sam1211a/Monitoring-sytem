package handlers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func Getip() string {
	addre, _ := net.InterfaceAddrs()
	for _, addr := range addre {
		if ipnet, ok := addr.(*net.IPAddr); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}
func UploadImg(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "invalid file", 400)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read error", 500)
		return
	}
	filename := fmt.Sprintf("screens/%s.png", r.FormValue("ip"))
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		http.Error(w, "save error", 400)
		return
	}
	// fmt.Println("Screenshot received")
	w.Write([]byte("ok"))
}
