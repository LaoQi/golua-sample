package main

import lua "github.com/yuin/gopher-lua"

func main() {
	// initialize lua vm pool with default script
	luaVM, err := initLuaVM()
	if err != nil {
		panic(err)
	}

	// call lua function
	if err := luaVM.simpleCallLuaFunction("Reply", lua.LString("")); err != nil {
		panic(err)
	}

	// reload luaVM when scripts has changed
	if err := luaVM.reload(); err != nil {
		panic(err)
	}

	luaVM.close()
}
