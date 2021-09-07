package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/davidwalter0/go-cfg"
	"github.com/davidwalter0/go-flag"
)

var subCmd string
var args []string
var once = sync.Once{}

// Args from command line remove sub command and returning subCmd name
// and args
var _ = func() {
	Args()
}

func Args() (string, []string) {
	once.Do(func() {
		if len(os.Args) > 1 {
			subCmd = os.Args[1]
			if len(os.Args) > 2 {
				var args = os.Args[2:]
				os.Args = []string{os.Args[0]}
				os.Args = append(os.Args, args[:]...)
			}
			args = os.Args
		}
	})
	return subCmd, args
}
func Run(subCmd string, pgm interface{}) (err error) {
	return
}

func IsStructPtr(i interface{}) (ptr bool) {
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct
}
func Init(subCmd string, pgm interface{}) (cfgd interface{}, err error) {
	v := reflect.ValueOf(pgm)
	if !IsStructPtr(pgm) {
		fmt.Println(fmt.Errorf("interface is not a struct ptr"))
		os.Exit(1)
	}
	v = v.Elem()
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		var name = typeOfS.Field(i).Name
		if subCmd == name {

			err = cfg.Flags(v.Field(i).Interface())

			if err != nil {
				fmt.Println(err)
				Help(pgm)
				os.Exit(1)
			}
			var text = []byte{}
			cfgd = v.Field(i).Interface()
			text, err = json.MarshalIndent(cfgd, "", "  ")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(text))
		}
	}
	return
}

func Do() (err error) {
	return
}

func Help(pgm interface{}) (text string) {
	var err error
	var parts = strings.Split(os.Args[0], "/")
	var l = len(parts)
	var Program = parts[l-1]
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", Program)
	v := reflect.ValueOf(pgm)
	if !IsStructPtr(pgm) {
		fmt.Fprintln(os.Stderr, fmt.Errorf("interface is not a struct ptr"))
		os.Exit(1)
	}
	v = v.Elem()
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		flag.CommandLine.Reset()
		var name = typeOfS.Field(i).Name
		err = cfg.Flags(v.Field(i).Interface())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		var doc string
		field, ok := reflect.TypeOf(pgm).Elem().FieldByName(name)
		if ok {
			doc = Tag(field, "doc")
		}
		fmt.Fprintf(os.Stderr, "%-15s\n  %s\n", name, doc)
		flag.PrintDefaults()
	}
	return
}

// Tag from struct field spec info
func Tag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}
