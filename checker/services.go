package checker

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

var Services = []Service{
	{
		Name: "Fortnite",
		Available: generateHTTPAvailableFunction(
			"GET",
			"https://www.epicgames.com/id/api/account/name/state/%s",
			func(response *http.Response) Status {
				_body, _ := ioutil.ReadAll(response.Body)
				body := string(_body)

				switch body {
				case "{\"exist\":true}":
					return StatusUnavailable
				case "{\"exist\":false}":
					return StatusAvailable
				default:
					return StatusUnknown
				}
			}),
	},
	{
		Name: "Github",
		Available: generateHTTPAvailableFunction(
			"GET",
			"https://github.com/%s",
			func(response *http.Response) Status {
				if response.StatusCode == 404 {
					return StatusUnavailable
				}
				if response.StatusCode == 200 {
					return StatusAvailable
				}
				log.Debugln(response.StatusCode)
				return StatusUnknown
			},
			),
	},
	{
		Name: "UnknownCheats",
		Available: func(username string) Status {
			endpoint := "https://www.unknowncheats.me/forum/ajax.php?do=verifyusername"
			// format the payload with the username
			data := []byte(fmt.Sprintf("securitytoken=guest&do=verifyusername&username=%s", username))
			// send the request
			req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
			if err != nil {
				log.Errorf("Error sending request to %s: %s", endpoint, err)
				return StatusUnknown
			}

			// Set the headers
			req.Header.Set("authority", "www.unknowncheats.me")
			req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")
			req.Header.Set("x-requested-with", "XMLHttpRequest")
			req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
			req.Header.Set("accept", "*/*")
			req.Header.Set("origin", "https://www.unknowncheats.me")
			req.Header.Set("sec-fetch-site", "same-origin")
			req.Header.Set("sec-fetch-mode", "cors")
			req.Header.Set("sec-fetch-dest", "empty")
			req.Header.Set("referer", "https://www.unknowncheats.me/forum/register.php?do=register")
			req.Header.Set("accept-language", "en-US,en;q=0.9")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Errorf("Error sending request to %s: %s", endpoint, err)
				return StatusUnknown
			}
			defer resp.Body.Close()

			_body, _ := ioutil.ReadAll(resp.Body)
			body := string(_body)

			if strings.Contains(body, "<status>valid</status>") {
				return StatusAvailable
			}
			if strings.Contains(body, "<status>invalid</status>") {
				return StatusUnavailable
			}
			log.Warnf("Received unknown status from %s", endpoint)
			return StatusUnknown
		},
	},
}

func GetServiceNames() []string {
	var serviceNames []string
	for _, service := range Services {
		serviceNames = append(serviceNames, service.Name)
	}
	return serviceNames
}

func GetService(name string) *Service {
	for _, service := range Services {
		if strings.ToLower(service.Name) == strings.ToLower(name) {
			return &service
		}
	}
	return nil
}
