package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Function to get api credentials from config file
func getCreds() (string, string) {
	// Read the config file stored in the same dir
	readCreds, err := os.Open("hooThat.config")
	if err != nil {
		log.Fatal(err)
	}
	credScanner := bufio.NewScanner(readCreds)
	credScanner.Split(bufio.ScanLines)
	var rtn []string
	var i = 0
	for credScanner.Scan() {
		rtn = append(rtn, credScanner.Text())
		i++
	}
	// Return first two lines from config file
	readCreds.Close()
	return rtn[0], rtn[1]
}

// Function to take in an IP, parse the resulting json data and return the country of origin.
func getCountry(ipIn string) string {
	// Generated struct from https://mholt.github.io/json-to-go/
	type AutoGenerated struct {
		Country struct {
			Names struct {
				En string `json:"en"`
			}
		}
	}
	var country AutoGenerated

	if net.ParseIP(ipIn) != nil {
		// Read geolite username/ password from config file for basic authentication.
		usr, pwd := getCreds()
		// Run the IP through geolite service to ID the country
		geoLocate, err := http.Get("https://" + usr + ":" + pwd + "@geolite.info/geoip/v2.1/country/" + ipIn + "?pretty")
		if err != nil {
			log.Fatal(err)
		} else {
			// Assign the geolite return to a variable and close the connection
			geoBody, err := ioutil.ReadAll(geoLocate.Body)
			geoLocate.Body.Close()
			if err != nil {
				log.Fatal(err)
			} else {
				// Parse the returned json and pass the desired info into a struct and return to caller.
				json.Unmarshal(geoBody, &country)
				return (country.Country.Names.En)
			}
		}
	}
	// Return details if there's an issue with IP passed in
	return string("error parsing IP: " + ipIn)
}

// Function to take an array of strings, count unique values and return the items and count in right html
func countryInstances(countries []string) {
	// Create a dictionary of values for each element passed and count unique entities.
	countryDict := make(map[string]int)
	for _, country := range countries {
		countryDict[country] = countryDict[country] + 1
	}
	// Print the items and count withing a html structure
	fmt.Println("<html><head><title>Traffic Stats</title></head><body>")
	for key, value := range countryDict {
		fmt.Println("<br><strong>"+key+":</strong>", value)
	}
	fmt.Println("</body></html>")
}

// Main function that parses IPs out of access.logs and passes data to the other functions.
func main() {
	// Read out access.log, final version will iterate through last X amount of logs.
	accessFile, err := os.Open("access.log.1")
	if err != nil {
		fmt.Print(err)
		accessFile.Close()
	} else {
		// Declare some required arrays for use later. Read in our log file
		var unique_ips []string
		var countryCount []string
		var i int
		accessScan := bufio.NewScanner(accessFile)
		accessScan.Split(bufio.ScanLines)
		// Splitting the log file line using space as the delimiter.
		for accessScan.Scan() {
			accessLine := strings.Split(accessScan.Text(), " ")
			i = 0
			//Compare the IP in the current line with the list of unique IPs already logged
			for _, ip := range unique_ips {
				//If the current IP is not unique the increase out iterator.
				if ip == accessLine[0] {
					i = i + 1
				}
			}
			// If iterator is zero we have a unique IP so we should log it. Run the final list of logged IPs against our geolite IP identification function
			if i == 0 {
				unique_ips = append(unique_ips, accessLine[0])
				countryCount = append(countryCount, getCountry(accessLine[0]))
			}
		}
		// Close log file once finished
		accessFile.Close()
		countryInstances(countryCount)
	}
}
