package utils

import (
	"bufio"
	"os"
	"strings"
)

// ReadFileLines reads a file line by line and returns a slice of strings
func ReadFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func FileEndsWithNewline(filePath string) (bool, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return false, err
	}
	if stat.Size() == 0 {
		return true, nil // Empty file, consider it ends with newline
	}
	buf := make([]byte, 1)
	_, err = f.ReadAt(buf, stat.Size()-1)
	if err != nil {
		return false, err
	}
	return buf[0] == '\n', nil
}
