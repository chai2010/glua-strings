// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package strings

import (
	"strconv"
	"strings"

	"github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("strings", Loader)
}

func Loader(L *lua.LState) int {
	mod := L.NewTable()
	L.SetFuncs(mod, stringsFuncs)
	L.Push(mod)
	return 1
}

var stringsFuncs = map[string]lua.LGFunction{
	"Compare": func(L *lua.LState) int {
		a := L.CheckString(1)
		b := L.CheckString(2)

		ret := strings.Compare(a, b)
		return retInt(L, ret)
	},
	"Contains": func(L *lua.LState) int {
		s := L.CheckString(1)
		substr := L.CheckString(2)

		ret := strings.Contains(s, substr)
		return retBool(L, ret)
	},
	"ContainsAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		chars := L.CheckString(2)

		ret := strings.ContainsAny(s, chars)
		return retBool(L, ret)
	},
	"ContainsRune": func(L *lua.LState) int {
		s := L.CheckString(1)
		r := L.CheckInt(2)

		ret := strings.ContainsRune(s, rune(r))
		return retBool(L, ret)
	},
	"Count": func(L *lua.LState) int {
		s := L.CheckString(1)
		substr := L.CheckString(2)

		ret := strings.Count(s, substr)
		return retInt(L, ret)
	},
	"EqualFold": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.EqualFold(s, t)
		return retBool(L, ret)
	},
	"Fields": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.Fields(s)
		return retStringList(L, ret)
	},
	"FieldsFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.FieldsFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return retStringList(L, ret)
	},
	"HasPrefix": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.HasPrefix(s, t)
		return retBool(L, ret)
	},
	"HasSuffix": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.HasSuffix(s, t)
		return retBool(L, ret)
	},
	"Index": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.Index(s, t)
		return retInt(L, ret)
	},
	"IndexAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.IndexAny(s, t)
		return retInt(L, ret)
	},
	"IndexByte": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.IndexByte(s, byte(t))
		return retInt(L, ret)
	},
	"IndexFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.IndexFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return retInt(L, ret)
	},
	"IndexRune": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.IndexRune(s, rune(t))
		return retInt(L, ret)
	},
	"Join": func(L *lua.LState) int {
		s := checkStringList(L, 1)
		t := L.CheckString(2)

		ret := strings.Join(s, t)
		return retString(L, ret)
	},
	"LastIndex": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.LastIndex(s, t)
		return retInt(L, ret)
	},
	"LastIndexAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.LastIndexAny(s, t)
		return retInt(L, ret)
	},
	"LastIndexByte": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.LastIndexByte(s, byte(t))
		return retInt(L, ret)
	},
	"LastIndexFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.IndexFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return retInt(L, ret)
	},
	"Map": func(L *lua.LState) int {
		fn := L.CheckFunction(1)
		s := L.CheckString(2)

		ret := strings.Map(
			func(r rune) rune {
				return callFunc_Rune_ret_Rune(
					L, fn, lua.LNumber(r),
				)
			},
			s,
		)
		return retString(L, ret)
	},
	"Repeat": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.Repeat(s, t)
		return retString(L, ret)
	},
	"Replace": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		z := L.CheckString(3)
		n := L.CheckInt(4)

		ret := strings.Replace(s, t, z, n)
		return retString(L, ret)
	},
	"Split": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.Split(s, t)
		return retStringList(L, ret)
	},
	"SplitAfter": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.SplitAfter(s, t)
		return retStringList(L, ret)
	},
	"SplitAfterN": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		n := L.CheckInt(3)

		ret := strings.SplitAfterN(s, t, n)
		return retStringList(L, ret)
	},
	"SplitN": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		n := L.CheckInt(3)

		ret := strings.SplitN(s, t, n)
		return retStringList(L, ret)
	},
	"Title": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.Title(s)
		return retString(L, ret)
	},
	"ToLower": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToLower(s)
		return retString(L, ret)
	},
	"ToTitle": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToTitle(s)
		return retString(L, ret)
	},
	"ToUpper": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToUpper(s)
		return retString(L, ret)
	},
	"Trim": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.Trim(s, cutset)
		return retString(L, ret)
	},
	"TrimFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return retString(L, ret)
	},
	"TrimLeft": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.TrimLeft(s, cutset)
		return retString(L, ret)
	},
	"TrimLeftFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimLeftFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return retString(L, ret)
	},
	"TrimPrefix": func(L *lua.LState) int {
		s := L.CheckString(1)
		prefix := L.CheckString(2)

		ret := strings.TrimPrefix(s, prefix)
		return retString(L, ret)
	},
	"TrimRight": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.TrimRight(s, cutset)
		return retString(L, ret)
	},
	"TrimRightFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimRightFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return retString(L, ret)
	},
	"TrimSpace": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.TrimSpace(s)
		return retString(L, ret)
	},
	"TrimSuffix": func(L *lua.LState) int {
		s := L.CheckString(1)
		suffix := L.CheckString(2)

		ret := strings.TrimSuffix(s, suffix)
		return retString(L, ret)
	},
}

func checkIntList(L *lua.LState, n int) []int {
	v := L.Get(n)
	if tb, ok := v.(*lua.LTable); ok {
		var ret []int
		for i := 0; i < tb.Len(); i++ {
			item := tb.RawGetInt(i)
			if lv, ok := item.(lua.LNumber); ok {
				ret = append(ret, int(lv))
			} else {
				x, _ := strconv.Atoi(lv.String())
				ret = append(ret, x)
			}
		}
		return ret
	}
	L.TypeError(n, lua.LTTable)
	return nil
}
func checkStringList(L *lua.LState, n int) []string {
	v := L.Get(n)
	if tb, ok := v.(*lua.LTable); ok {
		var ret []string
		for i := 0; i < tb.Len(); i++ {
			item := tb.RawGetInt(i)
			if lv, ok := item.(lua.LString); ok {
				ret = append(ret, string(lv))
			} else {
				ret = append(ret, item.String())
			}
		}
		return ret
	}
	L.TypeError(n, lua.LTTable)
	return nil
}

func retBool(L *lua.LState, v bool) int {
	L.Push(lua.LBool(v))
	return 1
}

func retInt(L *lua.LState, v int) int {
	L.Push(lua.LNumber(v))
	return 1
}
func retIntList(L *lua.LState, vs []int) int {
	tb := L.NewTable()
	for _, v := range vs {
		tb.Append(lua.LNumber(v))
	}
	L.Push(tb)
	return 1
}

func retString(L *lua.LState, v string) int {
	L.Push(lua.LString(v))
	return 1
}
func retStringList(L *lua.LState, vs []string) int {
	tb := L.NewTable()
	for _, v := range vs {
		tb.Append(lua.LString(v))
	}
	L.Push(tb)
	return 1
}

// func(rune) bool
func callFunc_Rune_ret_Bool(L *lua.LState, lf *lua.LFunction, args ...lua.LValue) bool {
	err := L.CallByParam(lua.P{Protect: true, Fn: lf, NRet: 1}, args...)
	if err != nil {
		panic(err)
	}
	defer L.Pop(1)

	ret := L.CheckBool(-1)
	return ret
}

// func(rune) rune
func callFunc_Rune_ret_Rune(L *lua.LState, lf *lua.LFunction, args ...lua.LValue) rune {
	err := L.CallByParam(lua.P{Protect: true, Fn: lf, NRet: 1}, args...)
	if err != nil {
		panic(err)
	}
	defer L.Pop(1)

	ret := L.CheckInt(-1)
	return rune(ret)
}
