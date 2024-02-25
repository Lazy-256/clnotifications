package main

import (
	"clnotifications/clnotifications"
	"fmt"
)

func main() {
	//err := clnotifications.GetKeys()
	err := clnotifications.ClearValues()
	if err != nil {
		fmt.Printf("Error during GetKeys: %v", err)
	}

}
