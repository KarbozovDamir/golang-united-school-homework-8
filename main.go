package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Arguments map[string]string

func parseArgs() Arguments {
	var id = flag.String("id", "", "takes an id")
	var operation = flag.String("operation", "", "takes operations (add, list, findById, remove)")
	var item = flag.String("item", "", "takes user info")
	var fileName = flag.String("fileName", "", "takes file name")
	flag.Parse()

	return Arguments{
		"id":        *id,
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
}

func Perform(args Arguments, writer io.Writer) error {
	var operation func(Arguments, io.Writer) error

	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "add":
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		operation = add
	case "list":
		operation = list
	case "remove":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		operation = remove
	case "findById":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		operation = findById

	}

	err := operation(args, writer)
	if err != nil {
		return err
	}
	return nil
}

func add(args Arguments, writer io.Writer) error {
	var user User
	json.Unmarshal([]byte(args["item"]), &user)

	file, err := os.Open(args["fileName"])
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var users []User
	json.Unmarshal(data, &users)
	users = append(users, user)

	data, err = json.Marshal(users)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(args["filename"], data, 0644)
	if err != nil {
		panic(err)
	}
	return err
}

func list(args Arguments, writer io.Writer) error {
	filename := args["fileName"]

	content, err := readFile(filename)
	if err != nil {
		return err
	}

	writer.Write(content)

	return nil
}
func remove(args Arguments, writer io.Writer) error {

	file, err := os.Open(args["filename"])
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return err
	}

	for idx, el := range users {
		if el.Id == args["id"] {
			users = append(users[:idx], users[idx+1:]...)
			data, err := json.Marshal(users)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(args["filename"], data, 0644)
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}

func findById(args Arguments, writer io.Writer) error {
	file, err := os.Open(args["id"])
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var users []User
	json.Unmarshal(data, &users)

	for _, el := range users {
		if el.Id == args["id"] {
			data, err := json.Marshal(el)
			if err != nil {
				panic(err)
			}
			fmt.Fprintln(writer, string(data))
			return err
		}
	}
	return err
}

func readFile(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeFile(filename string, fileMode fs.FileMode, message []byte) error {
	err := os.WriteFile(filename, message, fileMode)
	if err != nil {
		return err
	}

	return nil
}

func parseUsers(fileContent []byte) error {
	var users User

	if len(fileContent) > 0 {
		err := json.Unmarshal(fileContent, &users)
		if err != nil {
			return fmt.Errorf("json parse error")
		}
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
