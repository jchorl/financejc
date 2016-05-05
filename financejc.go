package financejc

import (
	"net/http"

	"handlers"
)

func init() {
	http.HandleFunc("/auth", handlers.AuthUser)
}
