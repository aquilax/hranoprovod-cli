// Hranoprovod is command line tracking tool. It supports nested recipes and
// custom defined tracking elements, which makes it perfect for tracking calories,
// nutrition data, exercises and other accumulative data.

package main

import (
	"log"
	"os"
)

func main() {
	if err := GetApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
