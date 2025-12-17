package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ManoloEsS/burrow_prototype/internal/engine"
)

const defaultPort = "8080"

func main() {
	ctx := context.Background()

	//get input
	fmt.Println("Input a method and url")
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		log.Fatalf("no input")
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	args := strings.Fields(line)

	var (
		method string
		url    string
	)

	if len(args) >= 2 {
		method = args[0]
		url = args[1]
	} else {
		method = args[0]
		url = "http://localhost:" + defaultPort
	}

	response, err := engine.GetRequest(&ctx, formatMethod(method), formatUrl(url), nil)
	if err != nil {
		log.Fatalf("invalid arguments: %v", err)
	}
	fmt.Println(response.StatusCode)

}

func formatMethod(method string) string {
	correctMethod := strings.ToUpper(method)
	return correctMethod

}

func formatUrl(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	if url == defaultPort {
		return "http://localhost:" + defaultPort
	}

	return "https://" + url
}
