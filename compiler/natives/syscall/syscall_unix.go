// +build js,!windows

package syscall

import (
	"github.com/gopherjs/gopherjs/js"
	"unsafe"
)

func init() {
	process := js.Global.Get("process")
	if !process.IsUndefined() {
		jsEnv := process.Get("env")
		envkeys := js.Global.Get("Object").Call("keys", jsEnv)
		envs = make([]string, envkeys.Length())
		for i := 0; i < envkeys.Length(); i++ {
			key := envkeys.Index(i).Str()
			envs[i] = key + "=" + jsEnv.Get(key).Str()
		}
	}
}

var syscallModule js.Object
var alreadyTriedToLoad = false
var minusOne = -1

func syscall(name string) js.Object {
	defer recover() // return nil
	if syscallModule == nil {
		if alreadyTriedToLoad {
			return nil
		}
		alreadyTriedToLoad = true
		require := js.Global.Get("require")
		if require.IsUndefined() {
			syscallHandler := js.Global.Get("$syscall")
			if !syscallHandler.IsUndefined() {
				return syscallHandler
			}
			panic("")
		}
		syscallModule = require.Invoke("syscall")
	}
	return syscallModule.Get(name)
}

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	if f := syscall("Syscall"); f != nil {
		r := f.Invoke(trap, a1, a2, a3)
		return uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())
	}
	if trap == SYS_WRITE && (a1 == 1 || a1 == 2) {
		b := js.Global.Call("go$sliceType", js.Global.Get("Go$Uint8")).New(js.InternalObject(a2)).Interface().([]byte)
		printToConsole(b)
		return uintptr(len(b)), 0, 0
	}
	printWarning()
	return uintptr(minusOne), 0, EACCES
}

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	if f := syscall("Syscall6"); f != nil {
		r := f.Invoke(trap, a1, a2, a3, a4, a5, a6)
		return uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())
	}
	printWarning()
	return uintptr(minusOne), 0, EACCES
}

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	if f := syscall("Syscall"); f != nil {
		r := f.Invoke(trap, a1, a2, a3)
		return uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())
	}
	printWarning()
	return uintptr(minusOne), 0, EACCES
}

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	if f := syscall("Syscall6"); f != nil {
		r := f.Invoke(trap, a1, a2, a3, a4, a5, a6)
		return uintptr(r.Index(0).Int()), uintptr(r.Index(1).Int()), Errno(r.Index(2).Int())
	}
	printWarning()
	return uintptr(minusOne), 0, EACCES
}

func BytePtrFromString(s string) (*byte, error) {
	return (*byte)(unsafe.Pointer(js.Global.Call("go$stringToBytes", s, true).Unsafe())), nil
}
