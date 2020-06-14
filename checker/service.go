package checker

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/*
 * The service type defines a service where usernames can be checked
 * This type contains a method that determines if a username is available or not
 */
type Service struct {
	Name      string
	Available availableFunction
}

type availableFunction func(string) Status
type Status int8

const (
	StatusUnavailable = Status(iota)
	StatusAvailable
	StatusInvalid
	StatusUnknown
)

/*
 * Generates a service that checks if a username is available by sending a `method` request to `endpoint`
 * The endpoint should be formatted with a %s where the username should go
 * The response is checked using the isAvailable function
 */
func generateHTTPAvailableFunction(method string, endpoint string, status func(response *http.Response) Status) availableFunction {
	return func(username string) Status {
		client := http.Client{}

		req, err := http.NewRequest(method, fmt.Sprintf(endpoint, username), nil)
		if err != nil {
			log.Errorf("Error sending request to %s: %s", endpoint, err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Error sending request to %s: %s", endpoint, err)
		}

		return status(resp)
	}
}
