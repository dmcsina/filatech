package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (app *application) setEnvironmentVariables(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	linenumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ":")
		key := strings.TrimSpace(parts[0])
		value := parts[1]
		os.Setenv(key, value)
		_, exists := os.LookupEnv(key)
		if exists {
			fmt.Printf("Successfully Set environment value of %s with the value of %s\n", key, value)
		} else {
			fmt.Println("Something went wrong")
		}

		linenumber++
	}

}
