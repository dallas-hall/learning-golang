package main

import (
	"fmt"
	"keyvaluestore"
	"log"
)

func main() {
	kvs, err := keyvaluestore.OpenStore("test/data/data.bin")
	if err != nil {
		log.Fatal(err)
	}

	err = kvs.Save()
	if err != nil {
		log.Fatal(err)
	}

	kvs.Set("key", "value")
	v, ok := kvs.Get("key")
	if !ok {
		log.Fatal(err)
	}
	fmt.Println(v, ok)

}
