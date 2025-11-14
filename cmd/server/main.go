package main

import (
	"fmt"
	"net/http"

	"github.com/subrat-dwi/shubserver/internal/app"
)

func main() {
	s := app.Setup()

	fmt.Printf("Listening to http://localhost%v\n", s.Addr)
	http.ListenAndServe(s.Addr, s.Router)
}
