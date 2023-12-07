package cli

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	filePath := "./test.txt"
	c := NewConfig([]string{"ggrep", filePath})
	fmt.Println(c)
}
