package main

import (
	"io"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

typeOfOperation:= map[string]bool{
	"add":      true,
	"list":     true,
	"findById": true,
	"remove":   true,
}

func Perform(args Arguments, writer io.Writer) error {
	if len(typeOfOperation) == 0 {
		return errors.New("-operation flag has to be specified")
	}
}

func add(filename string, item string) ([]byte, error) {
}

func findById(fileContents []byte, writer io.Writer, args Arguments) error {
}

func list(fileContents []byte, writer io.Writer) error {
}

func remove(fileContents []byte, writer io.Writer, args Arguments) error {
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
