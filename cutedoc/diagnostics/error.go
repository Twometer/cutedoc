package diagnostics

import "log"

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrintError(err error, msg string) {
	if err != nil {
		log.Printf("error: %s: %s\n", msg, err.Error())
	}
}
