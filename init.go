package abuse

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Init(path string) error {
	tempDomainMu.Lock()
	defer tempDomainMu.Unlock()

	// Close the file if its been opened previously.
	if additionsFile != nil {
		additionsFile.Close()
		additionsFile = nil
	}

	// Set additions path for future updates.
	additionsPath = path

	if len(additionsPath) > 0 {
		additionsFile, err := os.OpenFile(additionsPath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		// Load the additions into the map
		reader := bufio.NewReader(additionsFile)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					return err
				} else {
					break
				}
			}

			if s = strings.TrimSpace(s); len(s) > 0 {
				tempDomain[s] = true
			}
		}
	}

	return nil
}

func Close() error {
	if additionsFile != nil {
		if err := additionsFile.Sync(); err != nil {
			return err
		} else if err := additionsFile.Close(); err != nil {
			return err
		}
		additionsFile = nil
	}

	return nil
}
