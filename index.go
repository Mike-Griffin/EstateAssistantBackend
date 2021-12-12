package main

import (
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println(DB_USER)
	fmt.Println(DB_NAME)
	fmt.Println(DB_PASSWORD)
}
