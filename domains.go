package abuse

/*
Usage: use github.com/packetriot/abuse/gen to create the file 'domains-static.go'
that will statically initialize the tempDomain map with values originally from
domains.txt (or any new-line source).
*/

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	tempDomainMu sync.Mutex
	tempDomain   = make(map[string]bool)

	// Users may want to add domains dynamically:
	//
	additionsFile *os.File                   // Changes are saved here
	additionsPath = "/tmp/abuse-domains.txt" // Path to domain additions file
)

func IsTempEmail(email string) (bool, error) {
	tokens := strings.Split(email, "@")
	if len(tokens) > 1 {
		return IsAbusiveDomain(strings.ToLower(tokens[1])), nil
	}
	return false, fmt.Errorf("invalid email input")
}

func IsAbusiveDomain(domain string) bool {
	tempDomainMu.Lock()
	defer tempDomainMu.Unlock()

	_, exists := tempDomain[strings.ToLower(domain)]
	return exists
}

func Add(domain string) (err error) {
	tempDomainMu.Lock()
	defer tempDomainMu.Unlock()

	domain = strings.ToLower(domain)
	if tempDomain[domain] {
		return nil
	}

	// Add new domain to map
	tempDomain[domain] = true

	// When Init(string) is not called this will be nil, it will use the
	// default path and work as desired.
	if additionsFile == nil {
		if additionsFile, err = os.OpenFile(additionsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		}
	}

	if additionsFile != nil {
		if _, err = additionsFile.WriteString(fmt.Sprintf("%s\n", domain)); err == nil {
			err = additionsFile.Sync()
		}
	}

	return err
}

func AbusiveDomains() map[string]bool {
	return tempDomain
}
