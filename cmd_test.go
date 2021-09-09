// -*- mode:go;mode:go-playground -*-
// snippet of code @ 2021-09-06 14:11:03

// === Go Playground ===
// Execute the snippet with Ctl-Return
// Provide custom arguments to compile with Alt-Return
// Remove the snippet completely with its dir and all files M-x `go-playground-rm`

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func testSetup(name string) *Pgm {
	resetArgs()
	resetEnv()
	newArgs(name)
	var pgm = newPgm()
	Init(pgm)
	return pgm
}
func testSetupEnv(name string) *Pgm {
	resetArgs()
	resetEnv()
	newEnv(name)
	var pgm = newPgm()
	Init(pgm)
	return pgm
}
func Test_SubCmd(t *testing.T) {
	var pgm *Pgm
	defer trace.ScopedTrace()()
	pgm = testSetup("T0a")
	fmt.Fprintf(os.Stderr, "subCmd %s args %s %s\n", XsubCmd, Xargs, JSON(pgm))
	pgm = testSetup("T1a")
	fmt.Fprintf(os.Stderr, "subCmd %s args %s %s\n", XsubCmd, Xargs, JSON(pgm))
	// resetArgs()
	// resetEnv()
	// newArgs("T1a")
	// pgm = newPgm()
	// fmt.Fprintf(os.Stderr, "subCmd %s args %s %s\n", XsubCmd, Xargs, JSON(pgm))
}

func Test_IsStruct(t *testing.T) {
	trace.Off()
	var err error = NilInterfaceErr
	var pgm = newPgm()
	if !IsStruct(err) {
		t.Log("*err is not a struct")
	}
	if IsStruct(*pgm) && !IsStructPtr(*pgm) {
		t.Log("*pgm is a struct not a pointer")
	}
	if IsStructPtr(pgm) {
		t.Log("pgm &Pgm is a struct pointer")
	}
}

func Test_Args(t *testing.T) {
	var pgm = testSetup("t0a")
	if pgm == nil {
		t.Error(NilInterfaceErr)
	}
	if pgm.T0a.Int != test_IntArg {
		t.Error(pgm.T0a.Int, "wrong args value", test_IntArg)
	}
}

func Test_Env(t *testing.T) {
	var pgm = testSetupEnv("t1a")
	fmt.Println("pgm", JSON(pgm))
	fmt.Println("pgm.T1a", JSON(pgm.T1a))
	fmt.Println(JSON(pgm))
	fmt.Println("pgm.T1a", JSON(pgm.T1a))
	if pgm == nil {
		t.Error(NilInterfaceErr)
	}
	if pgm.T1a.Int != test_IntEnv {
		t.Error(pgm.T1a.Int, "env variable value set", test_IntEnv, JSON(pgm))
		//		panic(JSON(pgm))
	}
	if pgm.T1a.S != test_InformationEnv {
		t.Error(pgm.T1a.S, "env variable value set", test_InformationEnv)
		//		panic(JSON(cfgd))
	}
}

func Test_ResolveAddr(t *testing.T) {}

var test_IntArg = 42
var test_IntEnv = 3
var test_InformationEnv = "*env-set*"
var test_InformationArg = "*arg-set*"

var holdArgs = os.Args

// resetArgs sets the os.Args to the original list
func resetArgs() {
	os.Args = holdArgs
	// flag.CommandLine = new(flag.FlagSet)
	// flag.CommandLine.Reset()
}

// newArgs configures argument list and sets internal variables subcmd
// and args for tests
func newArgs(name string) {
	os.Args = []string{"cmd", name, "--int", fmt.Sprintf("%d", test_IntArg), "--information", test_InformationArg}
	fmt.Fprintln(os.Stderr, os.Args)
	SubCmdArgs()
}

func newEnv(name string) {
	os.Setenv("INT", fmt.Sprintf("%d", test_IntEnv))
	os.Setenv("INFORMATION", test_InformationEnv)
	os.Args = []string{"cmd", name}
}

func resetEnv() {
	os.Unsetenv("INT")
	os.Unsetenv("INFORMATION")
}

type Pgm struct {
	T0a *T0 `json:"t0a" doc:"Pgm is a command line program object"`
	// T0b T0  `doc:"Pgm is a command line program object"`
	// T1a *T1 `doc:"Pgm is a command line program object"`
	T1a *T1 `json:"t1a" doc:"Pgm is a command line program object"`
	// T1b T1  `doc:"Pgm is a command line program object"`
}

func newPgm() *Pgm {
	return &Pgm{
		T0a: &T0{},
		// T0b: T0{},
		T1a: &T1{},
		// T1b: T1{},
	}
}

type T0 struct {
	Int int    `json:"int"`
	S   string `json:"information" doc:"T1 subcommand"`
}

func (t *T0) Do() error {
	return nil
}
func (t *T0) Init() error {
	return nil
}
func (t *T0) Help() string {
	return fmt.Sprintf("Help: %T", t)
}

type T1 struct {
	Int int    `json:"int"`
	S   string `json:"information" doc:"T1 subcommand"`
}

func (t *T1) Do() error {
	return nil
}
func (t *T1) Init() error {
	return nil
}
func (t *T1) Help() string {
	return fmt.Sprintf("Help: %T", t)
}

func (pgm *Pgm) Init() (subcmd string, err error) {
	return
}
func JSON(i interface{}) string {
	var err error
	var text = []byte{}
	text, err = json.MarshalIndent(i, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(text)
}
