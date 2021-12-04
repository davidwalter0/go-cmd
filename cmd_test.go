// -*- mode:go;mode:go-playground -*-
// snippet of code @ 2021-09-06 14:11:03

// === Go Playground ===
// Execute the snippet with Ctl-Return
// Provide custom arguments to compile with Alt-Return
// Remove the snippet completely with its dir and all files M-x `go-playground-rm`

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davidwalter0/go-cfg"
)

func testSetup(name string, which int) (subCmd string, pgm *Pgm, args []string) {
	resetArgs()
	resetEnv()
	pgm = newPgm()
	var err error
	// fmt.Fprintln(os.Stderr, cfg.Nest(pgm))
	// var err = cfg.Flags(pgm)
	// fmt.Fprintln(os.Stderr, err)
	// fmt.Fprintln(os.Stderr, JSON(pgm))
	var i interface{}
	var program string
	program = filepath.Base(os.Args[0])
	cfg.Reset(program)
	var holdArgs []string
	holdArgs = os.Args
	subCmd, args = newArgs(name, which)
	os.Args = args
	i, err = Init(program, subCmd, pgm)
	fmt.Fprintln(os.Stderr, "pgm", JSON(pgm))
	os.Args = holdArgs
	// if false {
	// 	pgm = i.(*Pgm)
	// }
	switch subCmd {
	case "t0a":
		pgm.T0a = i.(T0)
	case "t1a":
		pgm.T1a = i.(T1)
	case "t2":
		pgm.T2 = *i.(*T2)
	}
	fmt.Fprintln(os.Stderr, err)

	fmt.Fprintln(os.Stderr, JSON(pgm))
	cfg.Usage()
	return
}
func testSetupEnv(name string, which int) (subCmd string, pgm *Pgm) {
	// resetArgs()
	// resetEnv()
	//	newEnv(name)
	var i interface{}
	var err error
	pgm = newPgm()
	var program string
	program = filepath.Base(os.Args[0])
	var holdArgs, rest []string
	holdArgs = os.Args
	subCmd, rest = newEnv(name, which)
	os.Args = rest
	i, err = Init(program, subCmd, pgm)
	fmt.Fprintln(os.Stderr, "testSetupEnv", err)
	os.Args = holdArgs
	// if false {
	// 	pgm = i.(*Pgm)
	// }
	switch subCmd {
	case "t0a":
		pgm.T0a = *i.(*T0)
	case "t1a":
		pgm.T1a = *i.(*T1)
	case "t2":
		pgm.T2 = *i.(*T2)
	}

	return subCmd, pgm
}

func Test_Args(t *testing.T) {
	var subCmd, pgm, args = testSetup("t2", 2)
	defer trace.ScopedTrace("subCmd", fmt.Sprintf("subCmd %s args %s\n", subCmd, args))()

	if pgm == nil {
		t.Error(NilInterfaceErr)
	}
	if pgm.T2.Int2 != test_IntArg {
		t.Error(subCmd, pgm.T2.Int2, "wrong args value", test_IntArg)
	}
	if pgm.T2.S2 != test_InformationArg {
		t.Error(subCmd, pgm.T2.S2, "wrong args value", test_InformationArg)
	}

	if false {
		subCmd, pgm, args = testSetup("t0a", 0)
		defer trace.ScopedTrace("subCmd", fmt.Sprintf("subCmd %s args %s\n", subCmd, args))()
		if pgm == nil {
			t.Error(NilInterfaceErr)
		}
		if pgm.T0a.Int0 != test_IntArg {
			t.Error(subCmd, pgm.T0a.Int0, "wrong args value", test_IntArg)
		}
	}
}

func Test_SubCmd(t *testing.T) {
	var pgm *Pgm
	var subCmd string
	var args []string
	subCmd, pgm, args = testSetup("T2", 2)
	// cfg.Usage()
	// subCmd, pgm, args = testSetup("T1a", 1)
	defer trace.ScopedTrace("subCmd", fmt.Sprintf("subCmd %s args %s\n", subCmd, args))()
	fmt.Fprintf(os.Stderr, "subCmd %s args %s %s\n", subCmd, args, JSON(pgm))

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

func Test_Env(t *testing.T) {
	var subCmd, pgm = testSetupEnv("t1a", 1)
	fmt.Println("pgm", JSON(pgm))
	fmt.Println("pgm.T1a", JSON(pgm.T1a))
	fmt.Println(JSON(pgm))
	fmt.Println("pgm.T1a", JSON(pgm.T1a))
	if pgm == nil {
		t.Error(NilInterfaceErr)
	}
	if pgm.T1a.Int1 != test_IntEnv {
		t.Error(subCmd, pgm.T1a.Int1, "env variable value set", test_IntEnv, JSON(pgm))
		//		panic(JSON(pgm))
	}
	if pgm.T1a.S1 != test_InformationEnv {
		t.Error(subCmd, pgm.T1a.S1, "env variable value set", test_InformationEnv)
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
func newArgs(name string, which int) (subCmd string, rest []string) {
	var args []string
	name = strings.ToLower(name)
	switch which {
	case 0:
		args = []string{"cmd",
			name,
			"--int0", fmt.Sprintf("%d", test_IntArg),
			"--information0", test_InformationArg,
		}
	case 1:
		args = []string{"cmd",
			name,
			"--int1", fmt.Sprintf("%d", test_IntArg),
			"--information1", test_InformationArg,
		}
	default:
		fallthrough
	case 2:
		args = []string{"cmd",
			name,
			// "--int0", fmt.Sprintf("%d", test_IntArg),
			// "--information0", test_InformationArg,
			// "--int1", fmt.Sprintf("%d", test_IntArg),
			// "--information1", test_InformationArg,
			"--int2", fmt.Sprintf("%d", test_IntArg),
			"--information2", test_InformationArg,
			"--int3", fmt.Sprintf("%d", test_IntArg),
			"--information3", test_InformationArg,
			"--int4", fmt.Sprintf("%d", test_IntArg),
			"--information3", test_InformationArg,
		} // os.Args = []string{
		// 	"cmd",
		// 	name,
		// }
		fmt.Fprintln(os.Stderr, os.Args)
		subCmd, rest = SubCmdArgs(args)
		// fmt.Fprintln(os.Stderr, os.Args)
	}
	return
}

func newEnv(name string, which int) (subCmd string, rest []string) {
	os.Setenv("INT0", fmt.Sprintf("%d", test_IntEnv))
	os.Setenv("INFORMATION0", test_InformationEnv)
	os.Setenv("INT1", fmt.Sprintf("%d", test_IntEnv))
	os.Setenv("INFORMATION1", test_InformationEnv)
	os.Setenv("INT2", fmt.Sprintf("%d", test_IntEnv))
	os.Setenv("INFORMATION2", test_InformationEnv)
	os.Setenv("INT3", fmt.Sprintf("%d", test_IntEnv))
	os.Setenv("INFORMATION3", test_InformationEnv)
	os.Setenv("INT4", fmt.Sprintf("%d", test_IntEnv))
	os.Setenv("INFORMATION4", test_InformationEnv)

	subCmd, rest = SubCmdArgs([]string{"cmd", name})
	return
}

func resetEnv() {
	os.Unsetenv("INT")
	os.Unsetenv("INFORMATION")
}

type Pgm struct {
	T0a T0 `json:"t0a" doc:"Pgm is a command line program object"`
	T1a T1 `json:"t1a" doc:"Pgm is a command line program object"`
	T2  T2 `json:"t2" doc:"Pgm is a command line program object"`
}

func newPgm() *Pgm {
	return &Pgm{
		T0a: T0{},
		// T0b: T0{},
		T1a: T1{},
		// T1b: T1{},
		T2: T2{},
	}
}

type T0 struct {
	Int0 int    `json:"int0"`
	S0   string `json:"information0" doc:"T1 subcommand"`
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
	Int1 int    `json:"int1"`
	S1   string `json:"information1" doc:"T1 subcommand"`
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

type T3 struct {
	Int3 int    `json:"int3"`
	S3   string `json:"information3" doc:"T1 subcommand"`
}
type T4 struct {
	Int4 int    `json:"int4"`
	S4   string `json:"information4" doc:"T1 subcommand"`
}

type T2 struct {
	Int2 int    `json:"int2"`
	S2   string `json:"information2" doc:"T2 subcommand"`
	T3
	T4
}

func (t *T2) Do() error {
	return nil
}
func (t *T2) Init() error {
	return nil
}
func (t *T2) Help() string {
	return fmt.Sprintf("Help: %T", t)
}
