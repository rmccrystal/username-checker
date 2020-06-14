package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"username-checker/checker"
)

func main() {
	// Init the logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})

	// Parse the flags
	threads := flag.Int("threads", 20, "The number of threads used to check the usernames")
	inputFileName := flag.String("in", "", "The input file (defaults to STDIN)")
	outputFileName := flag.String("out", "", "The output file to write valid usernames (defaults to STDOUT")

	flag.Parse()

	// We're going to use the trailing args to get the serviceName
	if len(flag.Args()) == 0 {
		log.Fatalf("No service specified. Usage: ./username-checker (flags) <service name> \n\nAvailable services: %s", strings.Join(checker.GetServiceNames(), ", "))
	}

	serviceName := flag.Args()[0]

	service := checker.GetService(serviceName)

	// If the user didn't specify a valid service
	if service == nil {
		log.Fatalf("%s is not a valid service. \n\nAvailable services: %s", serviceName, strings.Join(checker.GetServiceNames(), ", "))
	}

	// Init inputFile
	var inputFile *os.File
	var err error
	if *inputFileName != "" {
		inputFile, err = os.Open(*inputFileName)
		if err != nil {
			log.Fatalf("Error reading input file: %s", err)
		}
	} else {
		inputFile = os.Stdin
	}

	var outputFile *os.File
	if *outputFileName != "" {
		outputFile, err = os.Open(*outputFileName)
		if err != nil {
			log.Fatalf("Error reading output file: %s", err)
		}
	} else {
		outputFile = os.Stdout
	}

	// Create a new checker
	ch := checker.NewChecker(service, *threads)
	log.Println("Initialized Checker")

	// Create a new goroutine to parse the results
	go func() {
		// iterate the results channel
		for result := range ch.Results {
			switch result.Status {
			case checker.StatusAvailable:
				_, err := fmt.Fprintln(outputFile, result.Username) // Write to stdout
				if err != nil {
					log.Panic("Error writing to output file: %s", err)
				}
				log.Printf("[AVAILABLE] %s", result.Username)
			case checker.StatusUnavailable:
				log.Printf("[UNAVAILABLE]: %s", result.Username)
			case checker.StatusInvalid:
				log.Printf("[INVALID]: %s", result.Username)
			case checker.StatusUnknown:
				log.Printf("[UNKNOWN]: %s", result.Username)
			default:
				log.Panic("Received invalid status")
			}
		}
	}()

	// Grab lines from inputFile
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		ch.Usernames <- scanner.Text()
	}

	select {}
}
