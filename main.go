package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"username-checker/checker"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})

	// If there are no args, return
	if len(os.Args) <= 1 {
		log.Fatalf("You must specify the service to use. \n\nAvailable services: %s", strings.Join(checker.GetServiceNames(), ", "))
	}
	// Get flags
	//var serviceName = flag.String("service", nil, )

	ch := checker.NewChecker(checker.GetService("unknowncheats"), 20)

	go func() {
		for result := range ch.Results {
			switch result.Status {
			case checker.StatusAvailable:
				fmt.Println(result.Username)	// Write to stdout
				log.Printf("[AVAILABLE] %s", result.Username)
			case checker.StatusUnavailable:
				log.Printf("[UNAVAILABLE]: %s", result.Username)
			case checker.StatusUnknown:
				log.Printf("[UNKNOWN]: %s", result.Username)
			default:
				log.Panic("Received invalid status")
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ch.Usernames <- scanner.Text()
	}

	select {}
}
