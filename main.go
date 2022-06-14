package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type user struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var users []user

func addOperation(fileName string, item string) ([]byte, error) {

	if item == "" {
		return nil, fmt.Errorf("-item flag has to be specified")
	}

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		f, _ := os.Create(fileName)
		defer f.Close()
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var newUser user

	json.Unmarshal(data, &users)
	json.Unmarshal([]byte(item), &newUser)

	for _, user := range users {
		if user.Id == newUser.Id {
			return []byte(fmt.Sprintf("Item with id %s already exists", newUser.Id)), nil
		}
	}

	users = append(users, newUser)

	result, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}


	return result, ioutil.WriteFile(fileName, result, 0777)

}

func listOperation(fileName string) ([]byte, error) {

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		file, _ = os.Create(fileName)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileContent, nil

}
func findByIdOperation(fileName, id string) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("-id flag has to be specified")
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &users)

	found := false

	for _, user := range users {
		if user.Id == id {
			res, err := json.Marshal(user)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}

	if !found {
		return []byte(""), nil
	}
	return nil, nil

}

func removeOperation(fileName, id string) ([]byte, error) {

	if id == "" {
		return nil, fmt.Errorf("-id flag has to be specified")
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &users)

	found := false

	for i, user := range users {
		if user.Id == id {
			users[i] = users[len(users)-1]
			users = users[:len(users)-1]
			found = true
			break
		}
	}
	if !found {
		return []byte(fmt.Sprintf("Item with id %s not found", id)), nil
	}

	result, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return result, ioutil.WriteFile(fileName, result, 0644)

}

func executeOperation(args Arguments) (result []byte, err error) {

	if args["fileName"] == "" {
		return nil, fmt.Errorf("-fileName flag has to be specified")
	}

	if args["operation"] == "" {
		return nil, fmt.Errorf("-operation flag has to be specified")
	}

	switch args["operation"] {
	case "add":
		return addOperation(args["fileName"], args["item"])
	case "list":
		return listOperation(args["fileName"])
	case "findById":
		return findByIdOperation(args["fileName"], args["id"])
	case "remove":
		return removeOperation(args["fileName"], args["id"])
	default:
		return nil, fmt.Errorf("Operation %s not allowed!", args["operation"])
	}

}

func Perform(args Arguments, writer io.Writer) error {

	data, err := executeOperation(args)
	if err != nil {
		return err
	}
	writer.Write(data)
	return nil
}

func parseArgs() (args Arguments) {

	var idFlag = flag.String("id", "", "id of the user")
	var itemFlag = flag.String("item", "", "json obj of user")
	var operationFlag = flag.String("operation", "", "operation to perform.")
	var fileNameFlag = flag.String("fileName", "", "the name of the file.")

	flag.Parse()

	return Arguments{
		"id":        *idFlag,
		"item":      *itemFlag,
		"operation": *operationFlag,
		"fileName":  *fileNameFlag,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
