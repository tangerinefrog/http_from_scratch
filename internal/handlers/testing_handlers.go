package handlers

import (
	"strconv"

	"github.com/tangerinefrog/http_from_scratch/internal/request"
	"github.com/tangerinefrog/http_from_scratch/internal/response"
)

func TestYourProblemHandler(w *response.Writer, r *request.Request) {
	body := []byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Problem on your side.</p></body></html>")

	w.AddHeader("Content-Type", "text/html")
	w.AddHeader("Content-Length", strconv.Itoa(len(body)))

	w.WriteHeaders(response.StatusBadRequest)
	w.Write(body)
}

func TestMyProblemHandler(w *response.Writer, r *request.Request) {
	body := []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>My bad.</p></body></html>")

	w.AddHeader("Content-Type", "text/html")
	w.AddHeader("Content-Length", strconv.Itoa(len(body)))

	w.WriteHeaders(response.StatusInternalServerError)
	w.Write(body)
}

func TestGoodHandler(w *response.Writer, r *request.Request) {
	body := []byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was OK.</p></body></html>")

	w.AddHeader("Content-Type", "text/html")
	w.AddHeader("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeaders(response.StatusOK)
	w.Write(body)
}
