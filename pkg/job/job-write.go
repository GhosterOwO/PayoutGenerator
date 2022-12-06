package job

import (
	"os"
)

func write(n string, s string) error {
	const flags = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	f, err := os.OpenFile(n, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(s)
	return err
}
