package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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
	err = os.WriteFile("screens/.png", data, 0644)
	if err != nil {
		http.Error(w, "save error", 400)
		return
	}
	fmt.Println("Screenshot received")
	w.Write([]byte("ok"))
}
