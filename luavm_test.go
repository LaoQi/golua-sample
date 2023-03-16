package main

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestCallLuaFunction(t *testing.T) {
	luaVM, err := initLuaVM()
	if err != nil {
		panic(err)
	}

	if err := luaVM.simpleCallLuaFunction("Reply", lua.LString("Hi go")); err != nil {
		panic(err)
	}
}

func TestLuaPreloadScript(t *testing.T) {
	_, err := initLuaVM()
	if err != nil {
		panic(err)
	}
}

func TestDoFile(t *testing.T) {
	luaVM, err := initLuaVM()
	if err != nil {
		panic(err)
	}
	if err := luaVM.simpleDoLuaFile("scripts/exec.lua"); err != nil {
		panic(err)
	}
}
