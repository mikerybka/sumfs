package main

import (
	"encoding/json"
	"fmt"

	"github.com/mikerybka/sumfs"
)

func main() {
	dir := "."
	fsys, err := sumfs.Read(dir)
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(fsys.Sums, "", "  ")
	fmt.Println(string(b))
}
