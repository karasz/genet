package genetlib

import (
	"bufio"
	"os"
)

func ReadLines(filename string) ([]string, error) {

	file, err := os.Open(filename)

	if err != nil {
		return []string{""}, err
	}

	defer file.Close()

	var ret []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	return ret, scanner.Err()
}
