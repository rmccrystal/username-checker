package checker

type Checker struct {
	Usernames chan string
	Service   *Service
	Results   chan Result
}

type Result struct {
	Username string
	Status   Status
}

func NewChecker(service *Service, threads int) *Checker {
	usernameChan := make(chan string, 100)
	resultsChan := make(chan Result, 100) // Make a results channel with a large buffer

	checker := Checker{
		Usernames: usernameChan,
		Results:   resultsChan,
		Service:   service,
	}

	// Create workers
	for i := 0; i <= threads; i++ {
		go checker.worker()
	}

	return &checker
}

func (c Checker) worker() {
	for {
		// Get a username
		username := <-c.Usernames

		var status Status

		// Check if the username is valid
		if c.Service.Valid(username) {
			// If the username is cached
			if cachedStatus := CacheGet(c.Service.Name, username); cachedStatus != StatusUnknown {
				status = cachedStatus
			} else {
				// Else, check the username and cache it
				status = c.Service.Available(username)
				// Cache only if the status is known
				if status != StatusUnknown {
					CacheAppend(c.Service.Name, username, status)
				}
			}
		} else {
			status = StatusInvalid
		}

		c.Results <- Result{
			Username: username,
			Status:   status,
		}
	}
}
