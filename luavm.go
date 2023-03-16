package main

import (
	"bufio"
	luajson "github.com/layeh/gopher-json"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"log"
	"os"
	"sync"
)

const DefaultScript = "scripts/main.lua"

type LuaVM struct {
	m       sync.Mutex
	saved   []*lua.LState
	version int
	proto   *lua.FunctionProto
}

func (lv *LuaVM) get() (*lua.LState, int, error) {
	lv.m.Lock()
	defer lv.m.Unlock()
	n := len(lv.saved)
	if n == 0 {
		return lv.new()
	}
	x := lv.saved[n-1]
	lv.saved = lv.saved[0 : n-1]
	return x, lv.version, nil
}

func (lv *LuaVM) new() (*lua.LState, int, error) {
	L := lua.NewState()
	L.PreloadModule(MyModuleName, luaModuleLoader)
	luajson.Preload(L)
	lf := L.NewFunctionFromProto(lv.proto)
	L.Push(lf)
	err := L.PCall(0, lua.MultRet, nil)
	return L, lv.version, err
}

func (lv *LuaVM) put(L *lua.LState, version int) {
	lv.m.Lock()
	defer lv.m.Unlock()
	if version != lv.version {
		L.Close()
		return
	}
	lv.saved = append(lv.saved, L)
}

func (lv *LuaVM) loadScript() error {
	file, err := os.Open(DefaultScript)
	defer file.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, DefaultScript)
	if err != nil {
		return err
	}
	proto, err := lua.Compile(chunk, DefaultScript)
	if err != nil {
		return err
	}
	lv.proto = proto
	return nil
}

func (lv *LuaVM) reload() error {
	if err := lv.loadScript(); err != nil {
		return err
	}
	lv.close()
	lv.version += 1
	log.Printf("Reload LuaVM over!")
	return nil
}

func (lv *LuaVM) close() {
	lv.m.Lock()
	defer lv.m.Unlock()
	for _, L := range lv.saved {
		L.Close()
	}
	lv.saved = lv.saved[0:0]
}

var luaVM *LuaVM

func (lv *LuaVM) simpleCallLuaFunction(functionName string, args ...lua.LValue) error {
	l, version, err := lv.get()
	if err != nil {
		return err
	}
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(functionName),
		NRet:    0,
		Protect: true,
	}, args...); err != nil {
		l.Close()
		return err
	}
	lv.put(l, version)
	return nil
}

func (lv *LuaVM) simpleDoLuaFile(path string) error {
	l, version, err := lv.get()
	if err != nil {
		return err
	}
	if err := l.DoFile(path); err != nil {
		l.Close()
		return err
	}
	lv.put(l, version)
	return nil
}

func initLuaVM() (*LuaVM, error) {
	luaVM = &LuaVM{
		saved:   make([]*lua.LState, 0, 4),
		version: 0,
	}

	if err := luaVM.loadScript(); err != nil {
		return luaVM, err
	}
	return luaVM, nil
}
