package main

import (
	"flag"
	"log"
	"os"
)

var pathToData string
var pathToClasses string
var pathToOutput string

func init() {
	flag.StringVar(&pathToData, "data", "", "path to csv file with dialogs")
	flag.StringVar(&pathToClasses, "classes", "", "path to json file with classes")
	flag.StringVar(&pathToOutput, "out", "", "path to output file")
	flag.Parse()
}

func ReadClasses(fileName string) []Class {
	file, err := os.ReadFile(pathToClasses) // For read access.

	if err != nil {
		log.Fatal(err)
	}

	return ParseJson(file)
}

func main() {
	classes := ReadClasses(pathToClasses)

	inData, err := os.ReadFile(pathToData)

	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create(pathToOutput)

	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	dialogs := GetDialogs(GetJsonDialogs(inData))
	_, err = outFile.Write(GetData(ProduceClasses(dialogs, classes)))

	if err != nil {
		log.Fatal(err)
	}
}
