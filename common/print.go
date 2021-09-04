package common

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintPretty(obj interface{}) {

	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	if _, err = os.Stdout.Write(b); err != nil {
		fmt.Println(err)
	}
}
