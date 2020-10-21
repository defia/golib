package golib

import (
	"bufio"
	"os"
	"strings"
)

func ScanFile(fn string, trim bool) ([]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	arr := make([]string, 0)
	for scanner.Scan() {
		str := scanner.Text()
		if trim {
			str = strings.TrimSpace(str)
		}
		arr = append(arr, str)
	}
	return arr, nil
}
