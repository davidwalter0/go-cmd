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
	"github.com/davidwalter0/tracer"
)

var trace = tracer.New()
var once = sync.Once{}

func SubCmdArgs(args []string) (subCmd string, rest []string) {
	defer trace.ScopedTrace(fmt.Sprintf("SubCmdArgs: args %s", args))()
	rest = []string{args[0]}
	if len(args) > 1 {
		subCmd = args[1]
		if len(args) > 2 {
			rest = append(rest, args[2:]...)
		}
	}
	return
}

func Args(args []string) (subCmd string, rest []string) {
	once.Do(func() {
		subCmd, rest = SubCmdArgs(args)
	})
	return
}
func Run(XsubCmd string, pgm interface{}) (err error) {
	return
}

var NilInterfaceErr = fmt.Errorf("Xargs is nil")
var NonStructErr = fmt.Errorf("arg is not a struct")

func IsStruct(i interface{}) bool {
	if i == nil {
		panic(NilInterfaceErr)
	}
	var v = reflect.ValueOf(i)
	return v.Kind() == reflect.Struct
}

func IsStructPtr(i interface{}) bool {
	if i == nil {
		panic(NilInterfaceErr)
	}
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct
}

func PointerFromInterface(name string, i interface{}) (r interface{}) {
	if i == nil {
		panic(NilInterfaceErr)
	}
	var err error
	var v = reflect.ValueOf(i)
	if IsStructPtr(i) {
		// defer trace.ScopedTrace("IsStructPtr")()
		err = fmt.Errorf("object [%s] kind [%s] interface is a struct ptr", name, v.Kind())
		fmt.Fprintln(os.Stderr, err)
		panic(err)
		r = v.Interface()
	} else if IsStruct(i) {
		var typeOfS = v.Type()
		fmt.Printf("%+v\n", typeOfS)
		err = fmt.Errorf("object [%s] kind [%s] interface is a struct not a struct ptr", name, v.Kind())
		fmt.Fprintln(os.Stderr, err)
		panic(err)
		r = i
	} else {
		var typeOfS = v.Type()
		var name = typeOfS.Name()
		// defer trace.ScopedTrace(name, v.Kind())
		err = fmt.Errorf("object [%s] kind [%s] interface is not a struct or struct ptr", name, v.Kind())
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
	return
}

func Do() (err error) {
	return
}

// Tag from struct field spec info
func Tag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func Init(program, subCmd string, pgm interface{}) (cfgd interface{}, err error) {
	defer trace.ScopedTrace()()
	// flag.CommandLine = flag.NewFlagSet("pgm", flag.ContinueOnError)
	// err = cfg.Flags(pgm)
	// return pgm, err
	if true {
		cfg.Reset(subCmd)
		var v = reflect.ValueOf(pgm).Elem()
		var typeOfS = v.Type()
		if !IsStructPtr(pgm) {
			err = fmt.Errorf("object [%s] kind [%s] interface is a struct ptr", typeOfS.Name(), v.Kind())
			fmt.Fprintln(os.Stderr, err)
			panic(err)
		}
		for i := 0; i < v.NumField(); i++ {
			// c, e := func() (interface{}, error) {
			defer trace.ScopedTrace(fmt.Sprintf("for loop %d", i), subCmd)()
			var name = strings.ToLower(typeOfS.Field(i).Name)
			fmt.Printf("subCmd %s name %s %+v %T\n", subCmd, name, cfgd, cfgd)
			fmt.Fprintln(os.Stderr, strings.ToLower(subCmd), strings.ToLower(name))
			if strings.ToLower(subCmd) == strings.ToLower(name) {
				fmt.Fprintf(os.Stderr, "> subCmd %s name %s %+v %T\n", subCmd, name, cfgd, cfgd)
				cfgd = v.Field(i).Addr().Interface()
				// fmt.Fprintf(os.Stderr, ">> subCmd %s name %s %+v %T\n", subCmd, name, cfgd, cfgd)
				if !IsStructPtr(cfgd) {
					cfgd = v.Field(i).Elem().Addr().Interface()
					// fmt.Fprintf(os.Stderr, ">>> subCmd %s name %s %+v %T\n", subCmd, name, cfgd, cfgd)
				}
				// err = cfg.NestWrap("cmd", cfgd)
				fmt.Fprintf(os.Stderr, "\n\n\n\n")
				err = cfg.Flags(cfgd)
				// err = cfg.Nest(cfgd)
				if err != nil {
					fmt.Println(err)
					Help(program, pgm)
					return nil, err
				}
				// var text = []byte{}
				// text, err = json.MarshalIndent(cfgd, "", "  ")
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// fmt.Println(string(text))
				return
			}
			// 	return cfgd, err
			// }()
			// if c != nil || e != nil {
			// 	cfgd, err = c, e
			// 	return
			// }
		}
		Help(program, pgm)
	}
	return
}

func Help(name string, pgm interface{}) {
	var program = name
	var err error
	var v = reflect.ValueOf(pgm).Elem()
	var typeOfS = v.Type()
	cfg.Usage = func() {}
	fmt.Fprintln(os.Stderr, "Usage", program)
	for i := 0; i < v.NumField(); i++ {
		//		cfg.Reset()
		var name = typeOfS.Field(i).Name
		var cfgd interface{}
		cfgd = v.Field(i).Addr().Interface()
		if !IsStructPtr(cfgd) {
			cfgd = v.Field(i).Elem().Addr().Interface()
		}
		err = cfg.Simple(cfgd)
		if err != nil {
			panic(err)
		}
		var doc string
		var f, ok = reflect.TypeOf(pgm).Elem().FieldByName(name)
		if ok {
			doc = Tag(f, "doc")
		}
		fmt.Fprintf(os.Stderr, "%-15s\n  %s\n", strings.ToLower(name), doc)
		flag.PrintDefaults()
	}
	return
}
func JSON(i interface{}) string {
	var err error
	var text = []byte{}
	text, err = json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(text)
}
