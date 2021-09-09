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
var XsubCmd string
var Xargs []string
var once = sync.Once{}

// Args from command line remove sub command and returning XsubCmd name
// and Xargs
var _ = func() {
	Args()
}

func SubCmdArgs() (string, []string) {
	defer trace.ScopedTrace()()
	fmt.Fprintln(os.Stderr, "SubCmdArgs: Xargs", Xargs)
	if len(os.Args) > 1 {
		XsubCmd = os.Args[1]
		if len(os.Args) > 2 {
			Xargs = os.Args[2:]
			os.Args = []string{os.Args[0]}
			os.Args = append(os.Args, Xargs[:]...)
		}
		Xargs = os.Args
	}
	// fmt.Fprintln(os.Stderr, "XsubCmd", XsubCmd, "Xargs", Xargs)
	return XsubCmd, Xargs
}

func Args() (string, []string) {
	once.Do(func() {
		SubCmdArgs()
	})
	return XsubCmd, Xargs
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

/*
func PointerFromReflectedValue(i interface{}) (r interface{}) {
	var err error
	var v = reflect.ValueOf(i)
	if IsStructPtr(v.Interface()) {
		r = v.Interface()
	} else if IsStruct(v.Addr().Interface()) {
		var typeOfS = v.Type()
		var name = typeOfS.Name()
		err = fmt.Errorf("object [%s] kind [%s] interface is not a struct ptr", name, v.Kind())
		fmt.Fprintln(os.Stderr, err)
		// r = v.Addr().Interface()
	} else {
		var typeOfS = v.Type()
		var name = typeOfS.Name()
		err = fmt.Errorf("object [%s] kind [%s] interface is not a struct ptr", name, v.Kind())
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
	return
}
*/
var x = 0

// func Init(XsubCmd string, pgm interface{}) (cfgd interface{}, err error) {
func Init(pgm interface{}) (cfgd interface{}, err error) {
	var v = reflect.ValueOf(pgm).Elem()
	typeOfS := v.Type()
	if !IsStructPtr(pgm) {
		err = fmt.Errorf("object [%s] kind [%s] interface is a struct ptr", typeOfS.Name(), v.Kind())
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
	for i := 0; i < v.NumField(); i++ {
		c, e := func() (interface{}, error) {
			defer trace.ScopedTrace(fmt.Sprintf("for loop %d", i))()
			var name = typeOfS.Field(i).Name
			fmt.Printf("XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
			if strings.ToLower(XsubCmd) == strings.ToLower(name) {
				//				fmt.Printf("XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
				fmt.Fprintf(os.Stderr, "> XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
				cfgd = v.Field(i).Elem().Addr().Interface()
				fmt.Fprintf(os.Stderr, ">> XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
				if !IsStructPtr(cfgd) {
					cfgd = v.Field(i).Elem().Addr().Interface()
					fmt.Fprintf(os.Stderr, ">>> XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
				}
				cfg.Reset()
				err = cfg.Flags(cfgd)
				if err != nil {
					fmt.Println(err)
					Help(pgm)
					return nil, err
				}
				var text = []byte{}
				text, err = json.MarshalIndent(cfgd, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(text))
				// if x > 0 {
				// 	panic(fmt.Sprintf("%+v %T\n", cfgd, cfgd))
				// }
				x++
			}
			return cfgd, err
		}()
		if c != nil || e != nil {
			cfgd, err = c, e
			return
		}
	}
	return
}

func Do() (err error) {
	return
}

func Help(pgm interface{}) {
	var parts = strings.Split(os.Args[0], "/")
	var l = len(parts)
	var program = parts[l-1]
	var err error
	var v = reflect.ValueOf(pgm).Elem()
	var typeOfS = v.Type()
	fmt.Fprintln(os.Stderr, "Usage", program)
	for i := 0; i < v.NumField(); i++ {
		var name = typeOfS.Field(i).Name
		var cfgd interface{}
		// fmt.Printf("XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
		// fmt.Fprintf(os.Stderr, "> XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
		cfgd = v.Field(i).Elem().Addr().Interface()
		// fmt.Fprintf(os.Stderr, ">> XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
		if !IsStructPtr(cfgd) {
			cfgd = v.Field(i).Elem().Addr().Interface()
			// fmt.Fprintf(os.Stderr, ">>> XsubCmd %s name %s %+v %T\n", XsubCmd, name, cfgd, cfgd)
		}
		cfg.Reset()
		err = cfg.Flags(cfgd)
		if err != nil {
			panic(err)
		}
		var doc string
		var f, ok = reflect.TypeOf(pgm).Elem().FieldByName(name)
		if ok {
			doc = Tag(f, "doc")
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
