package main

import (
	"flag"
	"io"
)




func Perform(args Arguments, writer io.Writer) error {

	operator := NewOperator(args)
	data, err := operator.ExecuteOperation()
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








