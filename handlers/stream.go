package handlers

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"monitoring/model"
	"net/http"
	"time"

	"github.com/kbinani/screenshot"
)

func capture() []byte {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("stream cp error", 400)
		return nil
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 60})
	return buf.Bytes()
}

func Streaming(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", 500)
		return
	}

	for {
		frame := model.LatestFrame[ip]
		if frame == nil {
			continue
		}

		fmt.Fprintf(w, "--frame\r\n")
		fmt.Fprintf(w, "Content-Type: image/jpeg\r\n\r\n")
		w.Write(frame)
		fmt.Fprintf(w, "\r\n")

		flusher.Flush()

		time.Sleep(200 * time.Millisecond) // ~5 FPS
	}
}
