package health

import (
	"net/http"

	"github.com/subrat-dwi/shubserver/internal/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, 200, map[string]string{
		"status": "ok",
	})
}
