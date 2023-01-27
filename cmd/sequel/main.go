package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	p "github.com/ethanmidgley/the-sequel/pkg/parser"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("the-sequel> ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		if input == ".exit" {
			break
		}

		c := p.Parse(input)

		switch c {
		case p.DELETE:
			fmt.Println("We want to delete here")
		case p.INSERT:
			fmt.Println("We want to insert here")
		case p.FETCH:
			fmt.Println("We want to fetch here")
		case p.UPDATE:
			fmt.Println("We want to update here")
		default:
			fmt.Printf("Unable to parse command: %s\n", input)
		}

		// exectute
		p.Extracter.Insert()

	}

}
