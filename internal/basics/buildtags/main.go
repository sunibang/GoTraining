package main

import (
	"fmt"
	"github.com/romangurevitch/go-training/internal/basics/buildtags/buildTags"
)

func main() {
	msg := buildTags.WelcomeMessage()
	fmt.Printf("%v", msg)
}
