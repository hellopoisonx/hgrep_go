package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

const HELPDOC = `
Usage: hgrep query path [-options]
	Options:
		-sensitive case sensitive
	Example:
		hgrpe test ./test.txt  -s
`

var (
	Sensitive *bool
)

type Config struct {
	FilePath  string
	Query     string
	Sensitive bool
	File      *os.File
}

func init() {
	Sensitive = flag.Bool("sensitive", false, "case sensitive")
}

func NewConfig() (*Config, error) {
	c := Config{}
	flag.Parse()
	c.Sensitive = *Sensitive
	args := os.Args
	if len(args) < 3 {
		return &c, errors.New("Error parsing arguments, lack arguments " + HELPDOC)
	}
	c.FilePath, c.Query = args[2], args[1]
	var err error
	c.File, err = os.OpenFile(c.FilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return &c, err
	}
	return &c, nil
}
func (c *Config) SearchInsensitive() error {
	f := c.File
	defer f.Close()
	bf := bufio.NewReader(f)
	for {
		buf := make([]byte, 4096)
		n, err := bf.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("Error reading the file: " + err.Error())
		}
		content := string(buf[:n])
		n = strings.Index(content, c.Query)
		if n == -1 {
			continue
		}
		fmt.Fprintf(os.Stdout, "%s", fmt.Sprintf("%s%s%s", content[:n], color.HiRedString(content[n:n+len(c.Query)]), content[n+len(c.Query):]))
	}
	return nil
}

func (c *Config) SearchSensitive() error {
	f := c.File
	defer f.Close()
	bf := bufio.NewReader(f)
	for {
		buf := make([]byte, 4096)
		n, err := bf.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.New("Error reading the file: " + err.Error())
		}
		content := string(buf[:n])
		n = strings.Index(strings.ToLower(string(content)), strings.ToLower(c.Query))
		if n == -1 {
			continue
		}
		fmt.Fprintf(os.Stdout, "%s", fmt.Sprintf("%s%s%s", content[:n], color.HiRedString(content[n:n+len(c.Query)]), content[n+len(c.Query):]))
	}
	return nil
}
func Run() {
	c, err := NewConfig()
	if err != nil {
		panic(err)
	}
	if !c.Sensitive {
		if err = c.SearchInsensitive(); err != nil {
			panic(err)
		}
	} else {
		if err = c.SearchSensitive(); err != nil {
			panic(err)
		}
	}
}
