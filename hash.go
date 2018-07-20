// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hash

import (
	"math/rand"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

const (
	c0 = uintptr((8-ptrSize)/4*2860486313 +
		(ptrSize-4)/4*33054211828000289)
	c1 = uintptr((8-ptrSize)/4*3267000013 +
		(ptrSize-4)/4*23344194077549503)
)

var (
	hashkey [4]uintptr
)

func init() {
	hashkeyInit(time.Now().UnixNano())
}

func hashkeyInit(seed int64) {
	rand.Seed(seed)
	rand.Read((*[uintptr(len(hashkey)) *
		ptrSize]byte)(unsafe.Pointer(&hashkey))[:])
	hashkey[0] |= 1 // make sure these numbers are odd
	hashkey[1] |= 1
	hashkey[2] |= 1
	hashkey[3] |= 1
}

type Hasher interface {
	Hash() uintptr
}

type SeededHasher interface {
	SeededHash(seed uintptr) uintptr
}

func Any(v interface{}, seed uintptr) uintptr {
	return hash(v, seed)
}

func hash(v interface{}, seed uintptr) uintptr {
	switch val := v.(type) {
	case SeededHasher:
		return val.SeededHash(seed)
	case Hasher:
		return val.Hash()
	case string:
		return String(val, seed)
	case []byte:
		return Bytes(val, seed)
	case int8:
		return Int8(val, seed)
	case uint8:
		return Uint8(val, seed)
	case int16:
		return Int16(val, seed)
	case uint16:
		return Uint16(val, seed)
	case int32:
		return Int32(val, seed)
	case uint32:
		return Uint32(val, seed)
	case int:
		return Int(val, seed)
	case uint:
		return Uint(val, seed)
	case int64:
		return Int64(val, seed)
	case uint64:
		return Uint64(val, seed)
	case uintptr:
		return memhash(unsafe.Pointer(&val), seed, ptrSize)
	case float32:
		return Float32(val, seed)
	case float64:
		return Float64(val, seed)
	case complex64:
		return Complex64(val, seed)
	case complex128:
		return Complex128(val, seed)
	case struct{}:
		return memhash0(unsafe.Pointer(&val), seed)
	default:
		return reflecthash(val, seed)
	}
}

func Unsafe(ptr unsafe.Pointer, size uintptr, seed uintptr) uintptr {
	return memhash(ptr, seed, size)
}

func Int8(val int8, seed uintptr) uintptr {
	return memhash8(unsafe.Pointer(&val), seed)
}

func Uint8(val uint8, seed uintptr) uintptr {
	return memhash8(unsafe.Pointer(&val), seed)
}

func Int16(val int16, seed uintptr) uintptr {
	return memhash16(unsafe.Pointer(&val), seed)
}

func Uint16(val uint16, seed uintptr) uintptr {
	return memhash16(unsafe.Pointer(&val), seed)
}

func Int32(val int32, seed uintptr) uintptr {
	return memhash32(unsafe.Pointer(&val), seed)
}

func Uint32(val uint32, seed uintptr) uintptr {
	return memhash32(unsafe.Pointer(&val), seed)
}

func Int(val int, seed uintptr) uintptr {
	return memhash32(unsafe.Pointer(&val), seed)
}

func Uint(val uint, seed uintptr) uintptr {
	return memhash32(unsafe.Pointer(&val), seed)
}

func Int64(val int64, seed uintptr) uintptr {
	return memhash64(unsafe.Pointer(&val), seed)
}

func Uint64(val uint64, seed uintptr) uintptr {
	return memhash64(unsafe.Pointer(&val), seed)
}

func String(s string, h uintptr) uintptr {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	out := memhash(unsafe.Pointer(hdr.Data), h, uintptr(hdr.Len))
	runtime.KeepAlive(&s)
	return out
}

func Bytes(b []byte, seed uintptr) uintptr {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	out := memhash(unsafe.Pointer(hdr.Data), seed, uintptr(hdr.Len))
	runtime.KeepAlive(&b)
	return out
}

func Float32(f float32, h uintptr) uintptr {
	switch {
	case f == 0:
		return c1 * (c0 ^ h) // +0, -0
	case f != f:
		return c1 * (c0 ^ h ^ uintptr(rand.Uint32())) //any kind of NaN
	default:
		return memhash(unsafe.Pointer(&f), h, 4)
	}
}

func Float64(f float64, h uintptr) uintptr {
	switch {
	case f == 0:
		return c1 * (c0 ^ h) // +0, -0
	case f != f:
		return c1 * (c0 ^ h ^ uintptr(rand.Uint32())) //any kind of NaN
	default:
		return memhash(unsafe.Pointer(&f), h, 8)
	}
}

func Complex64(c complex64, h uintptr) uintptr {
	p := unsafe.Pointer(&c)
	x := (*[2]float32)(p)
	return Float32(x[1], Float32(x[0], h))
}

func Complex128(c complex128, h uintptr) uintptr {
	p := unsafe.Pointer(&c)
	x := (*[2]float64)(p)
	return Float64(x[1], Float64(x[0], h))
}

func reflecthash(v interface{}, seed uintptr) uintptr {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		vp := reflect.New(val.Type())
		vp.Elem().Set(val)
		val = vp
	}
	valType := val.Type().Elem()
	ptr := unsafe.Pointer(val.Pointer())
	return c1 * memhash(ptr, seed^c0, valType.Size())
}

func memhash0(p unsafe.Pointer, h uintptr) uintptr {
	return h
}

func memhash8(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 1)
}

func memhash16(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 2)
}

func memhash128(p unsafe.Pointer, h uintptr) uintptr {
	return memhash(p, h, 16)
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}
