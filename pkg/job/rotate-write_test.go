package job

import "testing"

func TestNewfile(t *testing.T) {
	ctx := "123123\n"
	writer := NewRotateWriter("logs/alfa/set66.txt")
	for i := 1; i < 100; i += 1 {
		writer.Write([]byte(ctx))
	}
}
