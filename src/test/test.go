package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

func main() {
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)

	u1 = uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)

	u1 = uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)
}
