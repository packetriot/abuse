package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

var (
	inputPath string

	t *template.Template
)

func init() {
	flag.StringVar(&inputPath, "input", "", "path to include file, newline separate domain names")
}

func updateStaticGoFile() error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Load the additions into the map
	var domains []string
	domainMap := make(map[string]bool)

	reader := bufio.NewReader(f)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}

		s = strings.ToLower(strings.TrimSpace(s))
		if len(s) > 0 {
			// Do no insert duplicates, will break compile.
			if !domainMap[s] {
				domains = append(domains, s)
				domainMap[s] = true
			}
		}
	}

	t.ExecuteTemplate(os.Stdout, "default", domains)

	return nil
}

var staticGoTemplate = `package abuse

func init() {
	tempDomain = map[string]bool{
		{{ range $domain := . }}"{{ $domain }}": true,
		{{ end }}
	}
}
`

func main() {
	flag.Parse()
	if len(inputPath) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var err error
	if t, err = template.New("default").Parse(staticGoTemplate); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := updateStaticGoFile(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
