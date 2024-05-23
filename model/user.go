package model

import (
	"fmt"
)

type ViewData struct {
	Title   string
	Message string
}

func (V *ViewData) test() {
	fmt.Println("Hello World")
}
