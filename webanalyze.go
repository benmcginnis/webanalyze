package webanalyze

import (
	"bufio"
	"io"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
	// AppDefs provides access to the unmarshalled apps.json file
	AppDefs *AppsDefinition
)

// Result type encapsulates the result information from a given host
type Result struct {
	Host     string        `json:"host"`
	Matches  []Match       `json:"matches"`
	Duration time.Duration `json:"duration"`
	Error    error         `json:"error"`
}

// Init sets up all the workders, reads in the host data and returns the results channel or an error
func Init(workers int, hosts io.Reader, appsFile string) (chan Result, error) {
	results := make(chan Result)
	c := make(chan string)

	if err := loadApps(appsFile); err != nil {
		return results, err
	}

	// start worker
	initWorker(workers, c, results)

	// send hosts line by line to worker channel
	go func() {
		scanner := bufio.NewScanner(hosts)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		close(c)

		// wait for workers to finish, the close result channel to signal finish of scan
		wg.Wait()
		close(results)
	}()

	return results, nil
}
