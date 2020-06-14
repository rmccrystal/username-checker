/*
 * The Cache caches usernames statuses
 */
package checker

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Days that a cache entry is valid
const CacheExpireDays = 7

type cacheEntry struct {
	ServiceName string
	Username    string
	Status      Status
}

func (c cacheEntry) Serialize() string {
	return fmt.Sprintf("%s|%s|%d", c.ServiceName, c.Username, c.Status)
}

/*
 * The cache file contains lined formatted as such:
 * (serviceName)|(username)|(status)
 */
const cacheFileName = "cache.txt"

var cacheMutex = sync.Mutex{}

// Reads a list of cacheEntries from the cache file
func readCacheFile(fileName string) ([]cacheEntry, error) {
	// Lock the mutex
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Open our cache file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []cacheEntry
	// Iterate through lines in file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line with |
		split := strings.Split(line, "|")

		// If the line if not formatted correctly, return an error
		if len(split) != 3 {
			return nil, fmt.Errorf("line in cache file %s was not formatted correctly: %s", fileName, line)
		}

		// Parse the status number
		statusNum, err := strconv.Atoi(split[2])
		if err != nil {
			return nil, fmt.Errorf("status must be a number in cache file %s: %s", fileName, line)
		}

		// Make sure we have a valid status number
		if statusNum > int(StatusUnknown) {
			return nil, fmt.Errorf("status must be less than %d in cache file %s: %d", int(StatusUnknown), fileName, statusNum)
		}

		// Create our new cache entry from this line
		var newCacheEntry cacheEntry
		newCacheEntry.ServiceName = split[0]
		newCacheEntry.Username = split[1]
		newCacheEntry.Status = Status(statusNum)

		// Append our new entry
		entries = append(entries, newCacheEntry)
	}

	return entries, nil
}

func appendCacheFile(fileName string, entries []cacheEntry) error {
	// Lock the mutex
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Open our cache file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Loop over the entries we want to append and append them
	for _, entry := range entries {
		text := entry.Serialize()
		if _, err := file.WriteString(text + "\n"); err != nil {
			// return the error if there is an error writing
			return fmt.Errorf("error writing to cache file %s: %s", fileName, err)
		}
	}

	return nil
}

/*
 * Gets an item from the cache
 * If there is none, returns StatusUnknown
 */
func CacheGet(serviceName string, username string) Status {
	entries, err := readCacheFile(cacheFileName)
	if err != nil {
		log.Errorf("error reading cache file: %s", err)
	}

	// Iterate through the entries until we find an entry with the specified serviceName and username
	for _, entry := range entries {
		if entry.Username == username && entry.ServiceName == serviceName {
			return entry.Status
		}
	}

	// We didn't find an entry, return StatusUnknown
	return StatusUnknown
}

func CacheAppend(serviceName string, username string, status Status) {
	err := appendCacheFile(cacheFileName, []cacheEntry{{
		ServiceName: serviceName,
		Username:    username,
		Status:      status,
	}})

	if err != nil {
		log.Errorf("error appending to cache file: %s", err)
		return
	}
}
