package checker

type Checker struct {
	Usernames chan string
	Service   *Service
	Results   chan Result
}

type Result struct {
	Username string
	Status Status
}

func NewChecker(service *Service, threads int) *Checker {
	usernameChan := make(chan string, 100)
	resultsChan := make(chan Result, 100) // Make a results channel with a large buffer

	checker := Checker{
		Usernames: usernameChan,
		Results: resultsChan,
		Service: service,
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
		status := c.Service.Available(username)
		c.Results <- Result{
			Username: username,
			Status:   status,
		}
	}
}
