package main

import (
	"go/format"
	"io"
	"log"
	"os"
)

func main() {
	filename := os.Getenv("GOFILE")

	if len(filename) == 0 {
		if len(os.Args) < 2 {
			log.Fatal("programm syntax: go-validate <golang file>")
		}

		filename = os.Args[1]
	}

	file, err := ParseFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	content, err := PrepareTemplate(file)
	if err != nil {
		log.Fatal(err)
	}

	formatted, err := format.Source(content.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	genFilename := filename[:len(filename)-3] + "_validation_generated.go"
	if err := writeFile(genFilename, formatted); err != nil {
		log.Fatal(err)
	}
}

// Writing generated validation file.
func writeFile(filename string, content []byte) error {
	genFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer genFile.Close()

	if _, err := io.WriteString(genFile, string(content)); err != nil {
		return err
	}

	return nil
}
