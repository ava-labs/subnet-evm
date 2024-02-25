package utils

import (
	"os"
)

func AppendToFile(file string, data []byte) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Add a newline to the beginning of the data
	data = append([]byte("\n"), data...)
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
