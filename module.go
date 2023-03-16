package main

import (
	lua "github.com/yuin/gopher-lua"
	"log"
)

const MyModuleName = "MyModule"

func luaModuleLoader(L *lua.LState) int {
	// func placeholder
	L.SetGlobal(MyModuleName, L.NewFunction(emptyFunc))

	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)
	// returns the module
	L.Push(mod)
	return 1
}

func emptyFunc(L *lua.LState) int {
	log.Printf("[LuaVM] use empty function")
	return 0
}

var exports = map[string]lua.LGFunction{
	"sayHi": sayHi,
}

func sayHi(L *lua.LState) int {
	text := L.ToString(1)
	log.Printf("[LuaVM] %s %s", L.String(), text)
	return 0
}
