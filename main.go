package main

import (
	"fmt"

	"github.com/socketspace-jihad/omatdb/engine"
)

func main() {
	fmt.Println("welcome to omatdb")

	kv := engine.NewKVStore()

	kv.Store("some caching", 10)
	kv.Store("test", 200)

	fmt.Println(kv.Get("test"))
	kv.Delete("test")

	fmt.Println(kv.Get("test"))
}
