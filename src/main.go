package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

type Explorer struct {
	Name   string
	Length float64
}

type ByLength []Explorer

func (a ByLength) Len() int           { return len(a) }
func (a ByLength) Less(i, j int) bool { return a[i].Length > a[j].Length }
func (a ByLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	// Check if TMLU file paths are provided as arguments
	if len(os.Args) < 2 {
		fmt.Println("Please provide one or more TMLU file paths as arguments.")
		return
	}

	// Maps to store the total length of exploration for each explorer and surveyor
	explorerLengths := make(map[string]float64)
	surveyorLengths := make(map[string]float64)

	// Iterate through the XML file paths provided as arguments
	for _, filePath := range os.Args[1:] {
		// Create a new document
		doc := etree.NewDocument()

		// Load the XML file
		if err := doc.ReadFromFile(filePath); err != nil {
			fmt.Printf("Failed to read XML file %s: %v\n", filePath, err)
			continue
		}

		// Find all the //SRVD/EX elements
		explorers := doc.FindElements("//SRVD/EX")

		for _, explorer := range explorers {
			// Get the text content of the <EX> element
			exContent := explorer.Text()

			// Extract explorer and surveyor names using HTML notation
			explorerNames := extractNames(exContent, "<Explorer>", "</Explorer>")
			surveyorNames := extractNames(exContent, "<Surveyor>", "</Surveyor>")

			// Find the associated length of exploration (//SRVD/LG)
			lengthElement := explorer.FindElement("../LG")
			if lengthElement != nil {
				lgText := lengthElement.Text()

				// Convert the length to a float
				lgValue, err := strconv.ParseFloat(lgText, 64)
				if err != nil {
					fmt.Printf("Failed to convert //SRVD/LG value to float in file %s: %v\n", filePath, err)
					continue
				}

				// Add the length to each explorer's total length
				for _, explorerName := range explorerNames {
					explorerLengths[explorerName] += lgValue
				}

				// Add the length to each surveyor's total length
				for _, surveyorName := range surveyorNames {
					surveyorLengths[surveyorName] += lgValue
				}
			}
		}
	}

	// Sort explorers by their total lengths of exploration in descending order
	var sortedExplorers []Explorer
	for explorer, length := range explorerLengths {
		sortedExplorers = append(sortedExplorers, Explorer{Name: explorer, Length: length})
	}
	sort.Sort(ByLength(sortedExplorers))

	// Print the explorers and their total lengths of exploration
	fmt.Println("Explorers:")
	for _, ex := range sortedExplorers {
		fmt.Printf("%-20s Length: %f\n", ex.Name, ex.Length)
	}

	fmt.Println()

	// Sort surveyors by their total lengths of exploration in descending order
	var sortedSurveyors []Explorer
	for surveyor, length := range surveyorLengths {
		sortedSurveyors = append(sortedSurveyors, Explorer{Name: surveyor, Length: length})
	}
	sort.Sort(ByLength(sortedSurveyors))

	// Print the surveyors and their total lengths of exploration
	fmt.Println("Surveyors:")
	for _, surveyor := range sortedSurveyors {
		fmt.Printf("%-28s Length: %f\n", surveyor.Name, surveyor.Length)
	}
}

// extractNames extracts the names from the content between startTag and endTag.
func extractNames(content, startTag, endTag string) []string {
	var names []string

	// Find the start and end positions of the tags
	startIndex := strings.Index(content, startTag)
	endIndex := strings.Index(content, endTag)

	// Extract the substring between the tags
	if startIndex != -1 && endIndex != -1 {
		substring := content[startIndex+len(startTag) : endIndex]

		// Split the substring by commas and trim whitespace
		nameList := strings.Split(substring, ",")
		for _, name := range nameList {
			names = append(names, strings.TrimSpace(name))
		}
	}

	return names
}
