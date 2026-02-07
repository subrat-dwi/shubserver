package main

import (
	"fmt"
	"net/http"

	"github.com/subrat-dwi/shubserver/internal/app"
)

func main() {
	s := app.Setup()

	fmt.Printf("Listening to http://localhost%v\n", s.Addr)
	err := http.ListenAndServe(s.Addr, s.Router)

	if err != nil {
		fmt.Println("Server stopped:", err)
	}
}
