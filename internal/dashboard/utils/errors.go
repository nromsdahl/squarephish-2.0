package utils

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// RespondWithError is a utility function to respond with an error message
// and log the error message. It is used for serving internal server errors.
// Parameters:
//   - w: The http.ResponseWriter
//   - message: The message to respond with
//   - err: The error to log
//
// It returns nothing.
func RespondWithError(w http.ResponseWriter, message string, err error) {
	log.Printf("%s: %v", message, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// RespondWithErrorMessage is a utility function to respond with an error message
// and log the error message. It is used for error messages.
// Parameters:
//   - w: The http.ResponseWriter
//   - message: The message to respond with
//   - err: The error to log
//
// It returns nothing.
func RespondWithErrorMessage(w http.ResponseWriter, message string, err error) {
	log.Printf("%s: %v", message, err)
	fmt.Fprintf(w, "%s: %v", message, err)
}

// RespondWithMessage is a utility function to respond with a message
// and log the message. It is used for success messages.
// Parameters:
//   - w: The http.ResponseWriter
//   - message: The message to respond with
//
// It returns nothing.
func RespondWithMessage(w http.ResponseWriter, message string) {
	log.Print(message)
	fmt.Fprint(w, message)
}
