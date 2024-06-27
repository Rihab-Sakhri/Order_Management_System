package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	guuid "github.com/google/uuid"
	"github.com/sofiukl/oms/oms-core/models"
)

// PrintErrorf  prints the error
func PrintErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

// RespondWithError respond with error
func RespondWithError(w http.ResponseWriter, code int, message string, details string) {
	errStr := fmt.Sprintf("code: %b, error: %s", code, message)
	fmt.Println(errStr)
	RespondWithJSON(w, code, message, details, nil)
}

// RespondWithJSON respond with json
func RespondWithJSON(w http.ResponseWriter, code int, message string, details string, payload map[string]interface{}) {
	isError := false
	if code != 200 {
		isError = true
	}
	gResponse := models.GenericResponse{
		Error:   isError,
		Message: message,
		Result:  payload,
		Details: "",
	}
	response, _ := json.Marshal(gResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GenUUID generate uuid4
func GenUUID() string {
	id := guuid.New().String()
	return id
}

// ConvertToMap converts interface to map
func ConvertToMap(v interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(v)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}
