package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/beevik/etree"
)

func main() {
	// Create a new document
	doc := etree.NewDocument()

	// Load the XML file
	if err := doc.ReadFromFile("example.xmlBIG"); err != nil {
		fmt.Printf("Failed to read XML file: %v", err)
		return
	}

	// Pattern to match within <Explorer> tags
	pattern := "(?s)<Explorer>(.*?)" + regexp.QuoteMeta("Fred") + "(.*?)</Explorer>"
	regex := regexp.MustCompile(pattern) // Compile the regular expression once

	// Get all the //SRVD/EX elements
	explorers := doc.FindElements("//SRVD/EX")

	// Variable to store the sum of //SRVD/LG elements
	sum := 0.0

	for _, explorer := range explorers {
		// Get the text content of <Explorer> tag
		explorerText := explorer.Text()
		fmt.Printf("%v\n", explorerText)

		// Perform the match search on the text content
		matched := regex.MatchString(explorerText)

		// Check if the text matches the desired pattern
		if matched {
			// Find the associated <LG> element
			lgElements := explorer.Parent().FindElements("LG")
			for _, lgElement := range lgElements {
				// Get the text content of <LG> element
				lgText := lgElement.Text()

				// Convert the text to a float and add it to the sum
				lgValue, err := strconv.ParseFloat(lgText, 64)
				if err != nil {
					fmt.Printf("Failed to convert //SRVD/LG value to float: %v\n", err)
					continue
				}

				sum += lgValue
			}
		}
	}

	fmt.Printf("Sum of //SRVD/LG elements for matched //SRVD/EX: %f\n", sum)
}
