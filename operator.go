package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Arguments map[string]string


type Operator struct {
	args Arguments
}

func NewOperator(args Arguments) *Operator {
	return &Operator{
		args: args,
	}
}

func (o *Operator) List() ([]byte, error) {

	fileName := o.args["fileName"]

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
func (o *Operator) Add() ([]byte, error) {

	fileName := o.args["fileName"]
	item := o.args["item"]

	if item == "" {
		return nil, fmt.Errorf("-item flag has to be specified")
	}

	//check if file exists
	if _, err := os.Stat(fileName); err != nil {
		file, _ := os.Create(fileName)
		defer file.Close()
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var newUser User
	var newUsers []User

	json.Unmarshal(data, &newUsers)
	json.Unmarshal([]byte(item), &newUser)

	for _, user := range newUsers {
		if user.Id == newUser.Id {
			return []byte(fmt.Sprintf("Item with id %s already exists", newUser.Id)), nil
		}
	}

	users = append(newUsers, newUser)

	result, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(fileName, result, 0644)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (o *Operator) FindById() ([]byte, error) {

	fileName := o.args["fileName"]
	id := o.args["id"]

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

func (o *Operator) Remove() ([]byte, error) {

	id := o.args["id"]
	fileName := o.args["fileName"]

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

func (o *Operator) Default() ([]byte, error) {
	return nil, fmt.Errorf("Operation %s not allowed!", o.args["operation"])
}

func (o *Operator) ExecuteOperation() ([]byte, error) {

	if o.args["fileName"] == "" {
		return nil, fmt.Errorf("-fileName flag has to be specified")
	}

	if o.args["operation"] == "" {
		return nil, fmt.Errorf("-operation flag has to be specified")
	}

	switch o.args["operation"] {

	case "add":
		return o.Add()
	case "list":
		return o.List()
	case "findById":
		return o.FindById()
	case "remove":
		return o.Remove()
	default:
		return o.Default()
	}
}
