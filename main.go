// Hranoprovod is command line tracking tool. It supports nested recipes and
// custom defined tracking elements, which makes it perfect for tracking calories,
// nutrition data, exercises and other accumulative data.

package main

import (
	"log"
	"os"

	"github.com/aquilax/hranoprovod-cli/v2/cmd"
)

func main() {
	if err := cmd.GetApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
