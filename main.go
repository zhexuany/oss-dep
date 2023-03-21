package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type BOMData struct {
	BomFormat    string         `json:"bomFormat"`
	SpecVersion  string         `json:"specVersion"`
	SerialNumber string         `json:"serialNumber"`
	Version      int            `json:"version"`
	Metadata     Metadata       `json:"metadata"`
	Components   []Components   `json:"components"`
	Dependencies []Dependencies `json:"dependencies"`
}
type Hashes struct {
	Alg     string `json:"alg"`
	Content string `json:"content"`
}
type Tools struct {
	Vendor  string   `json:"vendor"`
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Hashes  []Hashes `json:"hashes"`
}
type License struct {
	ID string `json:"id"`
}
type Licenses struct {
	License License `json:"license"`
}
type ExternalReferences struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
type Component struct {
	Publisher          string               `json:"publisher"`
	Group              string               `json:"group"`
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	Description        string               `json:"description"`
	Licenses           []Licenses           `json:"licenses"`
	Purl               string               `json:"purl"`
	ExternalReferences []ExternalReferences `json:"externalReferences"`
	Type               string               `json:"type"`
	BomRef             string               `json:"bom-ref"`
}
type Metadata struct {
	Timestamp time.Time `json:"timestamp"`
	Tools     []Tools   `json:"tools"`
	Component Component `json:"component"`
}
type Components struct {
	Publisher          string               `json:"publisher,omitempty"`
	Group              string               `json:"group"`
	Name               string               `json:"name"`
	Version            string               `json:"version"`
	Description        string               `json:"description,omitempty"`
	Licenses           []Licenses           `json:"licenses,omitempty"`
	Purl               string               `json:"purl"`
	ExternalReferences []ExternalReferences `json:"externalReferences,omitempty"`
	Type               string               `json:"type"`
	BomRef             string               `json:"bom-ref"`
	Scope              string               `json:"scope,omitempty"`
	Hashes             []Hashes             `json:"hashes,omitempty"`
}
type Dependencies struct {
	Ref       string   `json:"ref"`
	DependsOn []string `json:"dependsOn"`
}

func main() {
	// define a flag for file name
	fileName := flag.String("in", "", "JSON file name")

	outputFile := flag.String("out", "", "csv file name")

	// parse flags
	flag.Parse()

	// check if file name is provided
	if *fileName == "" {
		fmt.Println("Please provide a file name")
		return
	}

	if *outputFile == "" {
		fmt.Println("Please provide a output file name")
		return
	}

	// Read the JSON data file
	jsonData, err := ioutil.ReadFile(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading JSON data file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the JSON data into a Data struct
	var data BOMData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON data: %v\n", err)
		os.Exit(1)
	}

	convertJson2Csv(data.Components, *outputFile)
}

func convertJson2Csv(components []Components, outputFile string) {
	// create a CSV file for writing
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create a CSV writer from the file
	writer := csv.NewWriter(file)
	// write header row with keys
	header := []string{"group", "name", "version", "license", "scope", "description", "hash-md5"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print component details
	for _, c := range components {
		row := make([]string, len(header))
		// assign values by matching keys

		row[0] = c.Group

		row[1] = c.Name

		row[2] = c.Version
		// check if hobbies is a list and keep only the first one

		if len(c.Licenses) > 0 {
			row[3] = c.Licenses[0].License.ID
		}

		row[4] = c.Scope

		row[5] = c.Description

		if len(c.Hashes) > 0 {
			row[6] = c.Hashes[0].Content
		}

		err = writer.Write(row)
	}

}
