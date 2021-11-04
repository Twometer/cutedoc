package utils

import "log"

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintError(err error, msg string) {
	if err != nil {
		log.Println("error: ", msg, "-", err.Error())
	}
}
