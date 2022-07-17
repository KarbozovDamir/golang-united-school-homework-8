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
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}

	err := operation(args, writer)
	if err != nil {
		return err
	}
	return nil
}

func add(args Arguments, writer io.Writer) error {
	item := []byte(args["item"])
	filename := args["fileName"]
	users := make([]User, 0)
	newUser := User{}

	err := json.Unmarshal(item, &newUser)
	content, err := readFile(filename)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		err = json.Unmarshal(content, &users)
		if err != nil {
			return nil
		}
	}

	for _, v := range users {
		if v.Id == newUser.Id {
			message := fmt.Sprintf("Item with id %s already exists", v.Id)
			writer.Write([]byte(message))
			return nil
		}
	}

	users = append(users, newUser)
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = writeFile(filename, 0644, data)
	if err != nil {
		return err
	}

	return nil
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
	userId := args["id"]
	filename := args["fileName"]

	users := make([]User, 0)

	content, err := readFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &users)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, t := range users {
		if t.Id == userId {
			message := fmt.Sprintf("{\"id\":\"%s\",\"email\":\"%s\",\"age\":%v}", t.Id, t.Email, t.Age)
			writer.Write([]byte(message))
			return nil
		}
	}

	return err
}

func readFile(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
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
