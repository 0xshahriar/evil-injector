package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ANSI color codes
const (
	Red   = "\033[31m"
	Reset = "\033[0m"
)

// Function to check if "evil.com" is present in the response headers
func checkForEvil(response *http.Response, domain string) bool {
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "evil.com") {
			fmt.Printf("%sHost header injection has been found on %s%s\n", Red, domain, Reset)
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <domain_list_file>")
		return
	}

	// Fancy banner
	banner := `
    ______      _ __   ____        _           __
   / ____/   __(_) /  /  _/___    (_)__  _____/ /_____  _____
  / __/ | | / / / /   / // __ \  / / _ \/ ___/ __/ __ \/ ___/
 / /___ | |/ / / /  _/ // / / / / /  __/ /__/ /_/ /_/ / /
/_____/ |___/_/_/  /___/_/ /_/_/ /\___/\___/\__/\____/_/
                            /___/	`
	fmt.Println(banner)
	fmt.Println("By 0xShahriar\n")

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Get total number of domains
	totalDomains := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalDomains++
	}

	file.Seek(0, 0) // Reset file pointer to the beginning

	// Read each domain from the file
	currentDomain := 0
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		currentDomain++
		domain := scanner.Text()
		url := "https://" + domain
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("X-Forwarded-Host", "evil.com")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		if checkForEvil(resp, domain) {
			continue
		}

		// Display progress
		fmt.Printf("\rProcessing domain %d/%d", currentDomain, totalDomains)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("\nProcessing completed.")
}
