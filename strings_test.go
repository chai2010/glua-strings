// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package strings_test

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"

	lua_strings "github.com/chai2010/glua-strings"
)

func setupLuaTest(t *testing.T, funcName string) *lua.LState {
	t.Helper()

	// create new Lua state with empty stack
	L := lua.NewState()
	// register 'strings' module in package.preload table
	L.PreloadModule("strings", lua_strings.Loader)

	// execute Lua code string
	if err := L.DoString(fmt.Sprintf(`
		local strings = require("strings") -- load module and push to stack
		%s = strings.%s -- create global reference to function
	`, funcName, funcName)); err != nil {
		t.Fatal(err)
	}

	return L
}

func setupLuaFuncTest(t *testing.T, funcName string, luaFunc string) *lua.LState {
	t.Helper()

	// create new Lua state
	L := lua.NewState()

	// register strings module
	L.PreloadModule("strings", lua_strings.Loader)

	// create the test wrapper function
	err := L.DoString(fmt.Sprintf(`
		local strings = require("strings")
		local fn = %s
		function test_%s(s)
			return strings.%s(s, fn)
		end
	`, luaFunc, funcName, funcName))

	require.NoError(t, err)
	return L
}

func toBool(L *lua.LState, idx int) bool {
	return L.ToBool(idx)
}

func toInt(L *lua.LState, idx int) int {
	return L.ToInt(idx)
}

func toStringSlice(table *lua.LTable) []string {
	result := make([]string, 0, table.Len())

	table.ForEach(func(_, value lua.LValue) {
		result = append(result, value.String())
	})

	return result
}

func callLuaFunc[T any](
	t *testing.T, L *lua.LState,
	funcName string, args []lua.LValue,
	converter func(*lua.LState, int) T,
) T {
	t.Helper()

	// push function onto the stack
	L.Push(L.GetGlobal(funcName))
	for _, arg := range args {
		// push each argument onto the stack
		L.Push(arg)
	}

	// execute function with len(args) arguments and 1 return value
	L.Call(len(args), 1)
	// convert value at top of stack to Go type
	result := converter(L, -1)
	// remove value from stack
	L.Pop(1)

	return result
}

func TestCompare(t *testing.T) {
	const luaFuncName = "Compare"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		a string
		b string
	}{
		{"hello", "hello"},
		{"hello", "world"},
		{"world", "hello"},
		{"", ""},
		{"", "hello"},
		{"世界", "世界"},
		{"hello", "HELLO"},
		{"\u0000", "\u0000"},
		{"αβ", "αγ"},
	}

	for i := range tests {
		expected := strings.Compare(tests[i].a, tests[i].b)

		args := []lua.LValue{
			lua.LString(tests[i].a),
			lua.LString(tests[i].b),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (a: %q, b: %q)",
			i, got, expected, tests[i].a, tests[i].b)
	}
}

func TestContains(t *testing.T) {
	const luaFuncName = "Contains"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		substr string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello world", "world"},
		{"hello world", "golang"},
		{"hello", "hell"},
		{"Hello", "hello"},
		{"hello hello", "hello"},
		{"你好世界", "世界"},
		{"你好世界", "goodbye"},
		{"\u0000hello", "\u0000"},
		{"", "\u0000"},
		{"aaa", "aa"},
		{"aaaaaaaaaaaaaaaa", "aaa"},
		{"\u0430\u0306", "\u0306"},
	}

	for i := range tests {
		expected := strings.Contains(tests[i].s, tests[i].substr)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].substr),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toBool)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, substr: %q)",
			i, got, expected, tests[i].s, tests[i].substr)
	}
}

func TestContainsAny(t *testing.T) {
	const luaFuncName = "ContainsAny"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s     string
		chars string
	}{
		{"", ""},
		{"hello", ""},
		{"", "abc"},
		{"hello", "h"},
		{"hello", "abc"},
		{"hello", "abch"},
		{"hello", "helo"},
		{"hello", "xyz"},
		{"你好世界", "好世"},
		{"你好世界", "abc"},
		{"hello!", "!@#"},
		{"hello", "!@#"},
		{"hello world", " "},
		{"hello", "ll"},
		{"\u0000hello", "\u0000abc"},
		{"αβγδ", "δζη"},
		{"hello", "\u0000"},
	}

	for i := range tests {
		expected := strings.ContainsAny(tests[i].s, tests[i].chars)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].chars),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toBool)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, chars: %q)",
			i, got, expected, tests[i].s, tests[i].chars)
	}
}

func TestContainsRune(t *testing.T) {
	const luaFuncName = "ContainsRune"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
		r int
	}{
		{"", 'a'},
		{"hello", 'h'},
		{"hello", 'x'},
		{"hello", 'l'},
		{"hello", 'o'},
		{"你好世界", '好'},
		{"你好世界", '谢'},
		{"hello!", '!'},
		{"hello\x00world", 0},
		{"hello world", ' '},
		{"hello世界", '界'},
		{"hello", 'l'},
		{"h3llo", '3'},
		{"hello", '9'},
		{"hello\U0010FFFF", 0x10FFFF},
		{"hello", 0xD800},
		{"hello", -1},
		{"hello", 0x110000},
	}

	for i := range tests {
		expected := strings.ContainsRune(tests[i].s, rune(tests[i].r))

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LNumber(tests[i].r),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toBool)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, rune: %d)",
			i, got, expected, tests[i].s, tests[i].r)
	}
}

func TestCount(t *testing.T) {
	const luaFuncName = "Count"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		substr string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "l"},
		{"hello", "ll"},
		{"hello hello", "hello"},
		{"hellohellohello", "hello"},
		{"你好世界你好", "你好"},
		{"aaa", "aa"},
		{"banana", "ana"},
		{"hello world", "xyz"},
		{" hello hello ", " "},
		{"....", "."},
		{"hello\nhello\n", "\n"},
		{"héllo héllo", "é"},
		{"aaaaa", "aa"},
		{"世界世界", "世界"},
	}

	for i := range tests {
		expected := strings.Count(tests[i].s, tests[i].substr)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].substr),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, substr: %q)",
			i, got, expected, tests[i].s, tests[i].substr)
	}
}

func TestEqualFold(t *testing.T) {
	const luaFuncName = "EqualFold"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
		t string
	}{
		{"", ""},
		{"hello", "hello"},
		{"Hello", "hello"},
		{"HELLO", "hello"},
		{"HeLLo", "hEllO"},
		{"world", "WORLD"},
		{"Go", "go"},
		{"σπίτι", "ΣΠΊΤΙ"},
		{"társasház", "TÁRSASHÁZ"},
		{"İstanbul", "istanbul"},
		{"München", "MÜNCHEN"},
		{"flambe\u0301", "FLAMBE\u0301"},
		{"hello world", "HELLO WORLD"},
		{"hello123", "HELLO123"},
		{"hello!", "HELLO!"},
		{"", "hello"},
		{"hello", ""},
		{"hello", "world"},
		{"hello", "hi"},
		{"στίγμα", "ΣΤΊΓΜΑ"},
	}

	for i := range tests {
		expected := strings.EqualFold(tests[i].s, tests[i].t)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].t),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toBool)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (s: %q, t: %q)",
			i, got, expected, tests[i].s, tests[i].t)
	}
}

func TestFields(t *testing.T) {
	const luaFuncName = "Fields"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
	}{
		{""},
		{" "},
		{"  "},
		{"\t"},
		{"\n"},
		{"\v"},
		{"\f"},
		{"\r"},
		{"hello"},
		{"hello world"},
		{"  hello   world  "},
		{"hello\tworld"},
		{"hello\nworld"},
		{"hello\rworld"},
		{"hello\vworld"},
		{"hello\fworld"},
		{"hello\t\n\v\f\r world"},
		{"世界 你好"},
		{"\u2000Hello\u2001World"},
		{"The\u00A0quick brown\u2000fox"},
		{"	 multiple	 spaces	 between	 words	 "},
		{"hello   world   how   are   you"},
		{"line1\n\nline2\n\n\nline3"},
		{"hello世界 你好world"},
		{"...hello...world..."},
		{"\u2028\u2029"},
		{"	 "},
		{"1 2\t3\n4"},
		{"日本 語"},
	}

	for i := range tests {
		expected := strings.Fields(tests[i].s)

		args := []lua.LValue{
			lua.LString(tests[i].s),
		}
		got := toStringSlice(callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) *lua.LTable {
			return L.CheckTable(idx)
		}))

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (s: %q)",
			i, got, expected, tests[i].s)
	}
}

func TestFieldsFunc(t *testing.T) {
	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) bool
	}{
		{
			s: "",
			luaFunc: `
				function(r)
					return true
				end
			`,
			goFunc: func(r rune) bool { return true },
		},
		{
			s: "你好世界한글",
			luaFunc: `
				function(r)
					return r > 0x4E00 and r < 0xD7A3
				end
			`,
			goFunc: func(r rune) bool {
				return unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hangul, r)
			},
		},
		{
			s: "αβγ,δεζ,ηθι",
			luaFunc: `
				function(r)
					return r >= 0x0370 and r <= 0x03FF
				end
			`,
			goFunc: func(r rune) bool { return unicode.Is(unicode.Greek, r) },
		},
		{
			s: "abc,def,ghi",
			luaFunc: `
				function(r)
					return r == string.byte(",")
				end
			`,
			goFunc: func(r rune) bool { return r == ',' },
		},
		{
			s: "a b\tc\nd",
			luaFunc: `
				function(r)
					return r == string.byte(" ") or
							r == string.byte("\t") or
							r == string.byte("\n")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' },
		},
		{
			s: "one☺two☺three",
			luaFunc: `
				function(r)
					return r == 0x263A
				end
			`,
			goFunc: func(r rune) bool { return r == '☺' },
		},
		{
			s: "   a   b   c   ",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "a::b:::c::::d",
			luaFunc: `
				function(r)
					return r == string.byte(":")
				end
			`,
			goFunc: func(r rune) bool { return r == ':' },
		},
		{
			s: "12ab34cd56",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "世界_hello_世界_goodbye",
			luaFunc: `
				function(r)
					return r == string.byte("_")
				end
			`,
			goFunc: func(r rune) bool { return r == '_' },
		},
		{
			s: "no-splits-here",
			luaFunc: `
				function(r)
					return r == string.byte("x")
				end
			`,
			goFunc: func(r rune) bool { return r == 'x' },
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return r == string.byte(",")
				end
			`,
			goFunc: func(r rune) bool { return r == ',' },
		},
		{
			s: "  ",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "a\u0000b\u0000c",
			luaFunc: `
				function(r)
					return r == 0
				end
			`,
			goFunc: func(r rune) bool { return r == 0 },
		},
		{
			s: "a\rb\n\tc",
			luaFunc: `
				function(r)
					return r == string.byte("\r") or
							r == string.byte("\n") or
							r == string.byte("\t")
				end
			`,
			goFunc: func(r rune) bool { return r == '\r' || r == '\n' || r == '\t' },
		},
		{
			s: "αβγ,δεζ,ηθι",
			luaFunc: `
				function(r)
					return r == string.byte(",")
				end
			`,
			goFunc: func(r rune) bool { return r == ',' },
		},
		{
			s: "!!a!!b!!c!!",
			luaFunc: `
				function(r)
					return r == string.byte("!")
				end
			`,
			goFunc: func(r rune) bool { return r == '!' },
		},
		{
			s: "a\u2028b\u2029c", // Unicode line/paragraph separators
			luaFunc: `
				function(r)
					return r == 0x2028 or r == 0x2029
				end
			`,
			goFunc: func(r rune) bool { return r == '\u2028' || r == '\u2029' },
		},
		{
			s: "αaβbγ",
			luaFunc: `
				function(r)
					return r >= 0x03B1 and r <= 0x03B3
				end
			`,
			goFunc: func(r rune) bool { return r >= 'α' && r <= 'γ' },
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			testL := setupLuaFuncTest(t, "FieldsFunc", tt.luaFunc)
			defer testL.Close()

			expected := strings.FieldsFunc(tt.s, tt.goFunc)

			testL.Push(testL.GetGlobal("test_FieldsFunc"))
			testL.Push(lua.LString(tt.s))
			testL.Call(1, 1)

			resultTable := testL.CheckTable(-1)
			got := make([]string, 0, resultTable.Len())
			resultTable.ForEach(func(_, value lua.LValue) {
				got = append(got, value.String())
			})

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, func: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestHasPrefix(t *testing.T) {
	const luaFuncName = "HasPrefix"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		prefix string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "hell"},
		{"hello", "hello"},
		{"hello", "hello1"},
		{"hello", "world"},
		{"Hello", "hello"},
		{"你好世界", "你好"},
		{"hello\nworld", "hello"},
		{"  hello", "  "},
		{"hello", "h"},
		{"αβγ", "αβ"},
		{"hello world", "hello "},
		{"\u0000hello", "\u0000"},
	}

	for i := range tests {
		expected := strings.HasPrefix(tests[i].s, tests[i].prefix)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].prefix),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toBool)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, prefix: %q)",
			i, got, expected, tests[i].s, tests[i].prefix)
	}
}

func TestHasSuffix(t *testing.T) {
	const luaFuncName = "HasSuffix"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		suffix string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "llo"},
		{"hello", "hello"},
		{"hello", "1hello"},
		{"hello", "world"},
		{"Hello", "hello"},
		{"你好世界", "世界"},
		{"hello\nworld", "world"},
		{"hello  ", "  "},
		{"hello", "o"},
		{"αβγ", "βγ"},
		{"hello world", " world"},
		{"hello\u0000", "\u0000"},
	}

	for i := range tests {
		expected := strings.HasSuffix(tests[i].s, tests[i].suffix)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].suffix),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toBool)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, suffix: %q)",
			i, got, expected, tests[i].s, tests[i].suffix)
	}
}

func TestIndex(t *testing.T) {
	const luaFuncName = "Index"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		substr string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "h"},
		{"hello", "hell"},
		{"hello", "hello"},
		{"hello", "world"},
		{"Hello", "hello"},
		{"你好世界", "世界"},
		{"hello\nworld", "world"},
		{"hellohello", "hello"},
		{"αβγδ", "βγ"},
		{"hello hello", "hello"},
		{"hello", "\u0000"},
		{"\u0000hello", "hello"},
	}

	for i := range tests {
		expected := strings.Index(tests[i].s, tests[i].substr)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].substr),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, substr: %q)",
			i, got, expected, tests[i].s, tests[i].substr)
	}
}

func TestIndexAny(t *testing.T) {
	const luaFuncName = "IndexAny"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s    string
		char string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "h"},
		{"hello", "aeiou"},
		{"hello", "xyz"},
		{"Hello", "abcH"},
		{"你好世界", "世界"},
		{"hello\nworld", "\n"},
		{"hello", "ol"},
		{"αβγδ", "γβ"},
		{"\u0000hello", "\u0000"},
		{"hello", "123"},
		{"hello", " \t\n"},
		{"hello", "llo"},
	}

	for i := range tests {
		expected := strings.IndexAny(tests[i].s, tests[i].char)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].char),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, chars: %q)",
			i, got, expected, tests[i].s, tests[i].char)
	}
}

func TestIndexByte(t *testing.T) {
	const luaFuncName = "IndexByte"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
		b int
	}{
		{"", 'a'},
		{"hello", 'h'},
		{"hello", 'e'},
		{"hello", 'o'},
		{"hello", 'x'},
		{"hello", 0},
		{"\u0000hello", 0},
		{"hello", '\n'},
		{"hello\nworld", '\n'},
		{"hello", ' '},
		{"hello world", ' '},
		{"hello", 255},
		{"hello", -1},
		{"hello", 256},
		{"hello", 'H'},
	}

	for i := range tests {
		expected := strings.IndexByte(tests[i].s, byte(tests[i].b))

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LNumber(tests[i].b),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, byte: %d)",
			i, got, expected, tests[i].s, tests[i].b)
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) bool
	}{
		{
			s: "hello",
			luaFunc: `
				function(r)
					return r == string.byte("l")
				end
			`,
			goFunc: func(r rune) bool { return r == 'l' },
		},
		{
			s: "hello123",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		// { // Lua VM has no UTF8 support
		// 	s: "你好世界",
		// 	luaFunc: `
		// 		function(r)
		// 			local utf8 = require("utf8")
		// 			return r == utf8.codepoint("好")
		// 		end
		// 	`,
		// 	goFunc: func(r rune) bool { return r == '好' },
		// },
		{
			s: "你好世界",
			luaFunc: `
				function(r)
					return r == 0x597D
				end
			`,
			goFunc: func(r rune) bool { return r == '好' },
		},
		{
			s: "hello world",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return true
				end
			`,
			goFunc: func(r rune) bool { return true },
		},
		{
			s: "αβγδ",
			luaFunc: `
				function(r)
					return r == 0x03B3
				end
			`,
			goFunc: func(r rune) bool { return r == 'γ' },
		},
		{
			s: "hello\u0000world",
			luaFunc: `
				function(r)
					return r == 0
				end
			`,
			goFunc: func(r rune) bool { return r == 0 },
		},
		{
			s: "hello世界",
			luaFunc: `
				function(r)
					return r > 0x4E00
				end
			`,
			goFunc: func(r rune) bool { return unicode.Is(unicode.Han, r) },
		},
		{
			s: "no match",
			luaFunc: `
				function(r)
					return false
				end
			`,
			goFunc: func(r rune) bool { return false },
		},
		{
			s: "HELLO",
			luaFunc: `
				function(r)
					return r >= string.byte("A") and r <= string.byte("Z")
				end
			`,
			goFunc: func(r rune) bool { return r >= 'A' && r <= 'Z' },
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			testL := setupLuaFuncTest(t, "IndexFunc", tt.luaFunc)
			defer testL.Close()

			expected := strings.IndexFunc(tt.s, tt.goFunc)

			testL.Push(testL.GetGlobal("test_IndexFunc"))
			testL.Push(lua.LString(tt.s))
			testL.Call(1, 1)

			got := testL.ToInt(-1)
			testL.Pop(1)

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, func: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestIndexRune(t *testing.T) {
	const luaFuncName = "IndexRune"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
		r int
	}{
		{"", 'a'},
		{"hello", 'h'},
		{"hello", 'e'},
		{"hello", 'l'},
		{"hello", 'o'},
		{"hello", 'x'},
		{"你好世界", '好'},
		{"hello世界", '界'},
		{"\u0000hello", 0},
		{"hello", -1},
		{"hello", 0x10FFFF},
		{"hello\U0010FFFF", 0x10FFFF},
		{"hello", 0xD800},
		{"αβγδ", 'β'},
		{"hello", 256},
		{"hello\U0010FFFF", 0x10FFFF},
		{"", -1},
		{"\uFFFD", 0xFFFD},
	}

	for i := range tests {
		expected := strings.IndexRune(tests[i].s, rune(tests[i].r))

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LNumber(tests[i].r),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, rune: %d)",
			i, got, expected, tests[i].s, tests[i].r)
	}
}

func TestJoin(t *testing.T) {
	const luaFuncName = "Join"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		elems []string
		sep   string
	}{
		{[]string{}, ""},
		{[]string{"a"}, ","},
		{[]string{"a", "b"}, ","},
		{[]string{"a", "b", "c"}, ","},
		{[]string{"hello", "world"}, " "},
		{[]string{"hello", "", "world"}, ","},
		{[]string{"你好", "世界"}, ""},
		{[]string{"a", "b", "c"}, "\n"},
		{[]string{"", "", ""}, ","},
		{[]string{"a", "b"}, "\u0000"},
		{[]string{"αβ", "γδ"}, "|"},
		{[]string{"hello", "world"}, "..."},
		{[]string{"\n", "\t"}, ","},
		{[]string{"a", "b", "c"}, ""},
	}

	for i := range tests {
		expected := strings.Join(tests[i].elems, tests[i].sep)

		elems := lua.LTable{}
		for _, e := range tests[i].elems {
			elems.Append(lua.LString(e))
		}

		args := []lua.LValue{
			&elems,
			lua.LString(tests[i].sep),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, sep: %q)",
			i, got, expected, tests[i].elems, tests[i].sep)
	}
}

func TestLastIndex(t *testing.T) {
	const luaFuncName = "LastIndex"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		substr string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "h"},
		{"hello", "hell"},
		{"hello", "hello"},
		{"hello", "world"},
		{"hello hello", "hello"},
		{"你好世界世界", "世界"},
		{"hello\nworld\n", "\n"},
		{"hellohello", "hello"},
		{"αβγδαβ", "αβ"},
		{"\u0000hello\u0000", "\u0000"},
		{"hello world world", "world"},
		{"aaa", "aa"},
	}

	for i := range tests {
		expected := strings.LastIndex(tests[i].s, tests[i].substr)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].substr),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, substr: %q)",
			i, got, expected, tests[i].s, tests[i].substr)
	}
}

func TestLastIndexAny(t *testing.T) {
	const luaFuncName = "LastIndexAny"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s     string
		chars string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "h"},
		{"hello", "ol"},
		{"hello", "xyz"},
		{"hello hello", "helo"},
		{"你好世界", "界世"},
		{"hello\nworld", "\n"},
		{"αβγδαβ", "βα"},
		{"\u0000hello\u0000", "\u0000"},
		{"hello world", " "},
		{"aeiou", "aeiou"},
		{"hello", "123"},
		{"hello", "lo"},
	}

	for i := range tests {
		expected := strings.LastIndexAny(tests[i].s, tests[i].chars)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].chars),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, chars: %q)",
			i, got, expected, tests[i].s, tests[i].chars)
	}
}

func TestLastIndexByte(t *testing.T) {
	const luaFuncName = "LastIndexByte"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
		b int
	}{
		{"", 'a'},
		{"hello", 'h'},
		{"hello", 'l'},
		{"hello", 'o'},
		{"hello", 'x'},
		{"hello hello", 'h'},
		{"\u0000hello\u0000", 0},
		{"hello\nworld\n", '\n'},
		{"hello world", ' '},
		{"aaa", 'a'},
		{"hello", 255},
		{"hello", -1},
		{"hello", 256},
		{"hello", 'H'},
		{"abc\xff", 255},
	}

	for i := range tests {
		expected := strings.LastIndexByte(tests[i].s, byte(tests[i].b))

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LNumber(tests[i].b),
		}
		got := callLuaFunc(t, L, luaFuncName, args, toInt)

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, byte: %d)",
			i, got, expected, tests[i].s, tests[i].b)
	}
}

func TestLastIndexFunc(t *testing.T) {
	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) bool
	}{
		{
			s: "hello",
			luaFunc: `
				function(r)
					return r == string.byte("l")
				end
			`,
			goFunc: func(r rune) bool { return r == 'l' },
		},
		{
			s: "123hello123",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "世界你好世界",
			luaFunc: `
				function(r)
					return r == 0x4E16
				end
			`,
			goFunc: func(r rune) bool { return r == '世' },
		},
		{
			s: "hello world hello",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return true
				end
			`,
			goFunc: func(r rune) bool { return true },
		},
		{
			s: "αβγδαβγ",
			luaFunc: `
				function(r)
					return r == 0x03B3
				end
			`,
			goFunc: func(r rune) bool { return r == 'γ' },
		},
		{
			s: "hello\u0000world\u0000",
			luaFunc: `
				function(r)
					return r == 0
				end
			`,
			goFunc: func(r rune) bool { return r == 0 },
		},
		{
			s: "hello世界你好",
			luaFunc: `
				function(r)
					return r > 0x4E00
				end
			`,
			goFunc: func(r rune) bool { return r > 0x4E00 },
		},
		{
			s: "no match",
			luaFunc: `
				function(r)
					return false
				end
			`,
			goFunc: func(r rune) bool { return false },
		},
		{
			s: "HELLOhello",
			luaFunc: `
				function(r)
					return r >= string.byte("A") and r <= string.byte("Z")
				end
			`,
			goFunc: func(r rune) bool { return r >= 'A' && r <= 'Z' },
		},
		{
			s: "aaa",
			luaFunc: `
				function(r)
					return r == string.byte("a")
				end
			`,
			goFunc: func(r rune) bool { return r == 'a' },
		},
		{
			s: "hello\nhello\n",
			luaFunc: `
				function(r)
					return r == string.byte("\n")
				end
			`,
			goFunc: func(r rune) bool { return r == '\n' },
		},
		{
			s: "  hello  world  ",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "αβγαβγ",
			luaFunc: `
				function(r)
					return r == 0x03B2
				end
			`,
			goFunc: func(r rune) bool { return r == 'β' },
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			testL := setupLuaFuncTest(t, "LastIndexFunc", tt.luaFunc)
			defer testL.Close()

			expected := strings.LastIndexFunc(tt.s, tt.goFunc)

			testL.Push(testL.GetGlobal("test_LastIndexFunc"))
			testL.Push(lua.LString(tt.s))
			testL.Call(1, 1)

			got := testL.ToInt(-1)
			testL.Pop(1)

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, func: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestMap(t *testing.T) {
	const luaFuncName = "Map"

	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) rune
	}{
		{
			s: "",
			luaFunc: `
				function(r)
					return r
				end
			`,
			goFunc: func(r rune) rune { return r },
		},

		{
			s: "hello",
			luaFunc: `
				function(r)
					if r > 0x10FFFF then
						error("invalid rune")
					end
					return r
				end
		    `,
			goFunc: func(r rune) rune {
				if r > unicode.MaxRune {
					panic("invalid rune")
				}
				return r
			},
		},
		{
			s: string(unicode.MaxRune),
			luaFunc: `
				function(r)
					return r
				end
		    `,
			goFunc: func(r rune) rune { return r },
		},
		{
			s: "hello",
			luaFunc: `
				function(r)
					return r + 1
				end
			`,
			goFunc: func(r rune) rune { return r + 1 },
		},
		{
			s: "HELLO",
			luaFunc: `
				function(r)
					return r + 32
				end
			`,
			goFunc: func(r rune) rune { return r + 32 },
		},
		{
			s: "hello",
			luaFunc: `
				function(r)
					return r - 32
				end
			`,
			goFunc: func(r rune) rune { return r - 32 },
		},
		{
			s: "12345",
			luaFunc: `
				function(r)
					if r >= string.byte("0") and r <= string.byte("9") then
						return r - string.byte("0") + string.byte("a")
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r >= '0' && r <= '9' {
					return r - '0' + 'a'
				}
				return r
			},
		},
		{
			s: "hello世界",
			luaFunc: `
				function(r)
					if r > 0x4E00 then
						return r + 1
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r > 0x4E00 {
					return r + 1
				}
				return r
			},
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return r
				end
			`,
			goFunc: func(r rune) rune { return r },
		},
		{
			s: "αβγδ",
			luaFunc: `
				function(r)
					if r >= 0x03B1 and r <= 0x03B4 then
						return r - 0x03B1 + string.byte("a")
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r >= 'α' && r <= 'δ' {
					return r - 'α' + 'a'
				}
				return r
			},
		},
		{
			s: "hello\u0000world",
			luaFunc: `
				function(r)
					if r == 0 then
						return string.byte("-")
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r == 0 {
					return '-'
				}
				return r
			},
		},
		{
			s: "a1b2c3",
			luaFunc: `
				function(r)
					if r >= string.byte("0") and r <= string.byte("9") then
						return string.byte("*")
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r >= '0' && r <= '9' {
					return '*'
				}
				return r
			},
		},
		{
			s: "  hello  ",
			luaFunc: `
				function(r)
					if r == string.byte(" ") then
						return string.byte("_")
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r == ' ' {
					return '_'
				}
				return r
			},
		},
		{
			s: "\u0000hello\u0000",
			luaFunc: `
				function(r)
					if r == 0 then
						return string.byte("x")
					end
					return r
				end
			`,
			goFunc: func(r rune) rune {
				if r == 0 {
					return 'x'
				}
				return r
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			L := setupLuaTest(t, luaFuncName)
			defer L.Close()

			expected := strings.Map(tt.goFunc, tt.s)

			// Create the mapping function
			err := L.DoString("mapFunc = " + tt.luaFunc)
			if err != nil {
				t.Fatal(err)
			}

			// Call Map with 2 arguments and 1 return value
			if err := L.CallByParam(lua.P{
				Fn:      L.GetGlobal("Map"),
				NRet:    1,
				Protect: true,
			}, L.GetGlobal("mapFunc"), lua.LString(tt.s)); err != nil {
				t.Fatal(err)
			}

			got := L.ToString(-1)
			L.Pop(1)

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, mapping: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestRepeat(t *testing.T) {
	const luaFuncName = "Repeat"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
		n int
	}{
		{"", 0},
		{"", 5},
		{"a", 0},
		{"a", 1},
		{"a", 5},
		{"abc", 3},
		{"你好", 2},
		{"\n", 3},
		{"\u0000", 2},
		{"αβ", 4},
		{"hello", 0},
		{"hello", 1},
		{"hello", 2},
		{" ", 5},
		{"a", math.MaxInt32 / 2},
	}

	for i := range tests {
		expected := strings.Repeat(tests[i].s, tests[i].n)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LNumber(tests[i].n),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, count: %d)",
			i, got, expected, tests[i].s, tests[i].n)
	}
}

func TestReplace(t *testing.T) {
	const luaFuncName = "Replace"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s   string
		old string
		new string
		n   int
	}{
		{"", "", "", 0},
		{"hello", "", "x", 0},
		{"hello", "", "x", -1},
		{"hello", "l", "L", 1},
		{"hello", "l", "L", 2},
		{"hello", "l", "L", -1},
		{"hello hello", "hello", "hi", 1},
		{"hello hello", "hello", "hi", -1},
		{"你好世界", "世界", "朋友", 1},
		{"αβγαβγ", "αβ", "δε", 1},
		{"hello\nworld", "\n", " ", 1},
		{"aaa", "a", "b", 2},
		{"", "a", "b", -1},
		{"hello", "hello", "", 1},
		{"hello", "e", "ee", -1},
		{"hello", "l", "L", 0},
		{"hello", "l", "L", -2},
	}

	for i := range tests {
		expected := strings.Replace(tests[i].s, tests[i].old, tests[i].new, tests[i].n)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].old),
			lua.LString(tests[i].new),
			lua.LNumber(tests[i].n),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, old: %q, new: %q, n: %d)",
			i, got, expected, tests[i].s, tests[i].old, tests[i].new, tests[i].n)
	}
}

func TestSplit(t *testing.T) {
	const luaFuncName = "Split"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s   string
		sep string
	}{
		{"", ""},
		{"", ","},
		{"a", ""},
		{"a,b", ","},
		{"a,b,c", ","},
		{"abc", ""},
		{"hello world", " "},
		{"你好,世界", ","},
		{"αβ,γδ", ","},
		{"a\nb\nc", "\n"},
		{"a,,b", ","},
		{",a,b,", ","},
		{"a:b:c", ":"},
		{"aaa", "a"},
		{"hello", "xyz"},
		{"a\u0000b", "\u0000"},
		{"a\u0000b\u0000c", "\u0000"},
		{"世界你好", "你"},
		{"α,β,γ", ","},
		{",,,,", ","},
	}

	for i := range tests {
		expected := strings.Split(tests[i].s, tests[i].sep)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].sep),
		}
		got := toStringSlice(callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) *lua.LTable {
			return L.CheckTable(idx)
		}))

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, sep: %q)",
			i, got, expected, tests[i].s, tests[i].sep)
	}
}

func TestSplitAfter(t *testing.T) {
	const luaFuncName = "SplitAfter"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s   string
		sep string
	}{
		{"", ""},
		{"", ","},
		{"a", ""},
		{"a,b", ","},
		{"a,b,c", ","},
		{"hello world", " "},
		{"你好,世界", ","},
		{"αβ,γδ", ","},
		{"a\nb\nc", "\n"},
		{"a,,b", ","},
		{",a,b,", ","},
		{"a:b:c", ":"},
		{"aaa", "a"},
		{"hello", "xyz"},
		{"a\u0000b", "\u0000"},
	}

	for i := range tests {
		expected := strings.SplitAfter(tests[i].s, tests[i].sep)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].sep),
		}
		got := toStringSlice(callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) *lua.LTable {
			return L.CheckTable(idx)
		}))

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, sep: %q)",
			i, got, expected, tests[i].s, tests[i].sep)
	}
}

func TestSplitAfterN(t *testing.T) {
	const luaFuncName = "SplitAfterN"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s   string
		sep string
		n   int
	}{
		{"", "", 0},
		{"", ",", 1},
		{"a", "", 1},
		{"a,b,c", ",", 2},
		{"hello world space", " ", 2},
		{"你好,世界,朋友", ",", 2},
		{"αβ,γδ,εζ", ",", 3},
		{"a\nb\nc\n", "\n", 2},
		{"a,,b,,c", ",", 3},
		{",a,b,c,", ",", 0},
		{"a:b:c:d", ":", -1},
		{"aaa", "a", 2},
		{"hello", "xyz", 2},
		{"a\u0000b\u0000c", "\u0000", 2},
		{"a,b,c", ",", 1},
	}

	for i := range tests {
		expected := strings.SplitAfterN(tests[i].s, tests[i].sep, tests[i].n)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].sep),
			lua.LNumber(tests[i].n),
		}

		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) []string {
			if L.Get(idx) == lua.LNil {
				return nil
			}
			tbl := L.CheckTable(idx)
			return toStringSlice(tbl)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, sep: %q, n: %d)",
			i, got, expected, tests[i].s, tests[i].sep, tests[i].n)
	}
}

func TestSplitN(t *testing.T) {
	const luaFuncName = "SplitN"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s   string
		sep string
		n   int
	}{
		{"", "", 0},
		{"", ",", 1},
		{"a", "", 1},
		{"a,b,c", ",", 2},
		{"hello world space", " ", 2},
		{"你好,世界,朋友", ",", 2},
		{"αβ,γδ,εζ", ",", 3},
		{"a\nb\nc\n", "\n", 2},
		{"a,,b,,c", ",", 3},
		{",a,b,c,", ",", 0},
		{"a:b:c:d", ":", -1},
		{"aaa", "a", 2},
		{"hello", "xyz", 2},
		{"a\u0000b\u0000c", "\u0000", 2},
		{"a,b,c", ",", 1},
	}

	for i := range tests {
		expected := strings.SplitN(tests[i].s, tests[i].sep, tests[i].n)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].sep),
			lua.LNumber(tests[i].n),
		}

		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) []string {
			if L.Get(idx) == lua.LNil {
				return nil
			}
			tbl := L.CheckTable(idx)
			return toStringSlice(tbl)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %v but Go returned %v (string: %q, sep: %q, n: %d)",
			i, got, expected, tests[i].s, tests[i].sep, tests[i].n)
	}
}

func TestTitle(t *testing.T) {
	const luaFuncName = "Title"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
	}{
		{""},
		{"hello"},
		{"hello world"},
		{"HELLO WORLD"},
		{"hELLO wORLD"},
		{"hello_world"},
		{"hello-world"},
		{"hello123world"},
		{"123hello"},
		{"你好世界"},
		{"αβγδ"},
		{"hello  world"},
		{"hello\nworld"},
		{"helloWorld"},
		{"hello's world"},
	}

	for i := range tests {
		expected := strings.Title(tests[i].s)

		args := []lua.LValue{
			lua.LString(tests[i].s),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q)",
			i, got, expected, tests[i].s)
	}
}

func TestToLower(t *testing.T) {
	const luaFuncName = "ToLower"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
	}{
		{""},
		{"hello"},
		{"HELLO"},
		{"Hello World"},
		{"123"},
		{"HELLO123"},
		{"Hello_World"},
		{"你好世界"},
		{"ΓΑΒΓΔ"},
		{"hElLo"},
		{"\u0041\u0042"},
		{"Hello\nWorld"},
		{"!@#$%"},
		{"   HELLO   "},
		{"HELLO'S WORLD"},
	}

	for i := range tests {
		expected := strings.ToLower(tests[i].s)

		args := []lua.LValue{
			lua.LString(tests[i].s),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q)",
			i, got, expected, tests[i].s)
	}
}

func TestToTitle(t *testing.T) {
	const luaFuncName = "ToTitle"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
	}{
		{""},
		{"hello"},
		{"hello world"},
		{"HELLO WORLD"},
		{"hELLO wORLD"},
		{"hello_world"},
		{"hello-world"},
		{"hello123world"},
		{"123hello"},
		{"你好世界"},
		{"αβγδ"},
		{"hello  world"},
		{"hello\nworld"},
		{"helloWorld"},
		{"hello's world"},
	}

	for i := range tests {
		expected := strings.ToTitle(tests[i].s)

		args := []lua.LValue{
			lua.LString(tests[i].s),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q)",
			i, got, expected, tests[i].s)
	}
}

func TestToUpper(t *testing.T) {
	const luaFuncName = "ToUpper"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
	}{
		{""},
		{"hello"},
		{"HELLO"},
		{"Hello World"},
		{"123"},
		{"hello123"},
		{"hello_world"},
		{"你好世界"},
		{"αβγδ"},
		{"hElLo"},
		{"\u0061\u0062"},
		{"hello\nworld"},
		{"!@#$%"},
		{"   hello   "},
		{"hello's world"},
	}

	for i := range tests {
		expected := strings.ToUpper(tests[i].s)

		args := []lua.LValue{
			lua.LString(tests[i].s),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q)",
			i, got, expected, tests[i].s)
	}
}

func TestTrim(t *testing.T) {
	const luaFuncName = "Trim"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		cutset string
	}{
		{"", ""},
		{"hello", ""},
		{"   hello   ", " "},
		{"\t\nhello\n\t", "\t\n"},
		{"...hello...", "."},
		{"hello", "o"},
		{"...hello...", "."},
		{"123hello123", "123"},
		{"你好世界", "界"},
		{"αβγδ", "αδ"},
		{" \t\n\r", " \t\n\r"},
		{"-=-hello-=-", "-="},
		{"abchelloabc", "abc"},
		{"\u0000hello\u0000", "\u0000"},
		{"   ", " "},
		{"", "\u0000"},
		{"αααhelloααα", "α"},
		{"\u2028\u2029", "\u2028\u2029"},
		{"...世界...", "."},
	}

	for i := range tests {
		expected := strings.Trim(tests[i].s, tests[i].cutset)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].cutset),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, cutset: %q)",
			i, got, expected, tests[i].s, tests[i].cutset)
	}
}

func TestTrimFunc(t *testing.T) {
	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) bool
	}{
		{
			s: "   hello   ",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "\t\n\rhello\t\n\r",
			luaFunc: `
				function(r)
					return r == string.byte("\t") or r == string.byte("\n") or r == string.byte("\r")
				end
			`,
			goFunc: func(r rune) bool { return r == '\t' || r == '\n' || r == '\r' },
		},
		{
			s: "123hello123",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "世界hello世界",
			luaFunc: `
				function(r)
					return r > 0x4E00
				end
			`,
			goFunc: func(r rune) bool { return unicode.Is(unicode.Han, r) },
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return true
				end
			`,
			goFunc: func(r rune) bool { return true },
		},
		{
			s: "αβγhelloγβα",
			luaFunc: `
				function(r)
					return r >= 0x03B1 and r <= 0x03B3
				end
			`,
			goFunc: func(r rune) bool { return r >= 'α' && r <= 'γ' },
		},
		{
			s: "\u0000hello\u0000",
			luaFunc: `
				function(r)
					return r == 0
				end
			`,
			goFunc: func(r rune) bool { return r == 0 },
		},
		{
			s: "no trim",
			luaFunc: `
				function(r)
					return false
				end
			`,
			goFunc: func(r rune) bool { return false },
		},
		{
			s: "...hello...",
			luaFunc: `
				function(r)
					return r == string.byte(".")
				end
			`,
			goFunc: func(r rune) bool { return r == '.' },
		},
		{
			s: "  \t \n hello \n \t  ",
			luaFunc: `
				function(r)
					local c = string.char(r)
					return string.match(c, "%s") ~= nil
				end
			`,
			goFunc: func(r rune) bool { return unicode.IsSpace(r) },
		},
		{
			s: "12345",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "Hello",
			luaFunc: `
				function(r)
					return r >= string.byte("a") and r <= string.byte("z")
				end
			`,
			goFunc: func(r rune) bool { return r >= 'a' && r <= 'z' },
		},
		{
			s: "--==hello==--",
			luaFunc: `
				function(r)
					return r == string.byte("-") or r == string.byte("=")
				end
			`,
			goFunc: func(r rune) bool { return r == '-' || r == '=' },
		},
		{
			s: "世界hello世界世界",
			luaFunc: `
				function(r)
					return r == 0x4E16 or r == 0x754C
				end
			`,
			goFunc: func(r rune) bool { return r == '世' || r == '界' },
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			testL := setupLuaFuncTest(t, "TrimFunc", tt.luaFunc)
			defer testL.Close()

			expected := strings.TrimFunc(tt.s, tt.goFunc)

			testL.Push(testL.GetGlobal("test_TrimFunc"))
			testL.Push(lua.LString(tt.s))
			testL.Call(1, 1)

			got := testL.ToString(-1)
			testL.Pop(1)

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, func: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestTrimLeft(t *testing.T) {
	const luaFuncName = "TrimLeft"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		cutset string
	}{
		{"", ""},
		{"hello", ""},
		{"   hello", " "},
		{"\t\nhello", "\t\n"},
		{"...hello", "."},
		{"hello", "h"},
		{"...hello...", "."},
		{"123hello123", "123"},
		{"你好世界", "你"},
		{"αβγδ", "α"},
		{" \t\n\r", " \t\n\r"},
		{"-=-hello-=-", "-="},
		{"abchelloabc", "abc"},
		{"\u0000hello\u0000", "\u0000"},
		{"   hello   ", " "},
	}

	for i := range tests {
		expected := strings.TrimLeft(tests[i].s, tests[i].cutset)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].cutset),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, cutset: %q)",
			i, got, expected, tests[i].s, tests[i].cutset)
	}
}

func TestTrimLeftFunc(t *testing.T) {
	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) bool
	}{
		{
			s: "   hello   ",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "\t\n\rhello",
			luaFunc: `
				function(r)
					return r == string.byte("\t") or r == string.byte("\n") or r == string.byte("\r")
				end
			`,
			goFunc: func(r rune) bool { return r == '\t' || r == '\n' || r == '\r' },
		},
		{
			s: "123hello456",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "世界hello你好",
			luaFunc: `
				function(r)
					return r > 0x4E00  -- Greater than first Han character
				end
			`,
			goFunc: func(r rune) bool { return r > 0x4E00 },
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return true
				end
			`,
			goFunc: func(r rune) bool { return true },
		},
		{
			s: "αβγhelloαβγ",
			luaFunc: `
				function(r)
					return r >= 0x03B1 and r <= 0x03B3  -- α to γ
				end
			`,
			goFunc: func(r rune) bool { return r >= 'α' && r <= 'γ' },
		},
		{
			s: "hello\u0000world",
			luaFunc: `
				function(r)
					return r == 0
				end
			`,
			goFunc: func(r rune) bool { return r == 0 },
		},
		{
			s: "...hello...",
			luaFunc: `
				function(r)
					return r == string.byte(".")
				end
			`,
			goFunc: func(r rune) bool { return r == '.' },
		},
		{
			s: "  \t \n hello \n \t  ",
			luaFunc: `
				function(r)
					local c = string.char(r)
					return string.match(c, "%s") ~= nil
				end
			`,
			goFunc: func(r rune) bool { return unicode.IsSpace(r) },
		},
		{
			s: "12345hello",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "--==hello==--",
			luaFunc: `
				function(r)
					return r == string.byte("-") or r == string.byte("=")
				end
			`,
			goFunc: func(r rune) bool { return r == '-' || r == '=' },
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			testL := setupLuaFuncTest(t, "TrimLeftFunc", tt.luaFunc)
			defer testL.Close()

			expected := strings.TrimLeftFunc(tt.s, tt.goFunc)

			testL.Push(testL.GetGlobal("test_TrimLeftFunc"))
			testL.Push(lua.LString(tt.s))
			testL.Call(1, 1)

			got := testL.ToString(-1)
			testL.Pop(1)

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, func: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestTrimPrefix(t *testing.T) {
	const luaFuncName = "TrimPrefix"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		prefix string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "hell"},
		{"hello", "hello"},
		{"hello", "hello1"},
		{"hello", "world"},
		{"hello world", "hello "},
		{"你好世界", "你好"},
		{"αβγδ", "αβ"},
		{"hello\nworld", "hello\n"},
		{"  hello", "  "},
		{"hello", "h"},
		{"\u0000hello", "\u0000"},
		{"prefixhello", "prefix"},
	}

	for i := range tests {
		expected := strings.TrimPrefix(tests[i].s, tests[i].prefix)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].prefix),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, prefix: %q)",
			i, got, expected, tests[i].s, tests[i].prefix)
	}
}

func TestTrimRight(t *testing.T) {
	const luaFuncName = "TrimRight"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		cutset string
	}{
		{"", ""},
		{"hello", ""},
		{"hello   ", " "},
		{"hello\t\n", "\t\n"},
		{"hello...", "."},
		{"hello", "o"},
		{"...hello...", "."},
		{"123hello123", "123"},
		{"你好世界", "界"},
		{"αβγδ", "δ"},
		{" \t\n\r", " \t\n\r"},
		{"-=-hello-=-", "-="},
		{"abchelloabc", "abc"},
		{"\u0000hello\u0000", "\u0000"},
		{"   hello   ", " "},
	}

	for i := range tests {
		expected := strings.TrimRight(tests[i].s, tests[i].cutset)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].cutset),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, cutset: %q)",
			i, got, expected, tests[i].s, tests[i].cutset)
	}
}

func TestTrimRightFunc(t *testing.T) {
	tests := []struct {
		s       string
		luaFunc string
		goFunc  func(rune) bool
	}{
		{
			s: "   hello   ",
			luaFunc: `
				function(r)
					return r == string.byte(" ")
				end
			`,
			goFunc: func(r rune) bool { return r == ' ' },
		},
		{
			s: "hello\t\n\r",
			luaFunc: `
				function(r)
					return r == string.byte("\t") or r == string.byte("\n") or r == string.byte("\r")
				end
			`,
			goFunc: func(r rune) bool { return r == '\t' || r == '\n' || r == '\r' },
		},
		{
			s: "hello123world456",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "hello世界你好",
			luaFunc: `
				function(r)
					return r > 0x4E00
				end
			`,
			goFunc: func(r rune) bool { return unicode.Is(unicode.Han, r) },
		},
		{
			s: "",
			luaFunc: `
				function(r)
					return true
				end
			`,
			goFunc: func(r rune) bool { return true },
		},
		{
			s: "helloαβγ",
			luaFunc: `
				function(r)
					return r >= 0x03B1 and r <= 0x03B3
				end
			`,
			goFunc: func(r rune) bool { return r >= 'α' && r <= 'γ' },
		},
		{
			s: "world\u0000",
			luaFunc: `
				function(r)
					return r == 0
				end
			`,
			goFunc: func(r rune) bool { return r == 0 },
		},
		{
			s: "hello...",
			luaFunc: `
				function(r)
					return r == string.byte(".")
				end
			`,
			goFunc: func(r rune) bool { return r == '.' },
		},
		{
			s: "  hello \n \t  ",
			luaFunc: `
				function(r)
					local c = string.char(r)
					return string.match(c, "%s") ~= nil
				end
			`,
			goFunc: func(r rune) bool { return unicode.IsSpace(r) },
		},
		{
			s: "hello12345",
			luaFunc: `
				function(r)
					return r >= string.byte("0") and r <= string.byte("9")
				end
			`,
			goFunc: func(r rune) bool { return r >= '0' && r <= '9' },
		},
		{
			s: "hello==--",
			luaFunc: `
				function(r)
					return r == string.byte("-") or r == string.byte("=")
				end
			`,
			goFunc: func(r rune) bool { return r == '-' || r == '=' },
		},
		{
			s: "test\r\n\t ",
			luaFunc: `
				function(r)
					return r == string.byte(" ") or
							r == string.byte("\t") or
							r == string.byte("\r") or
							r == string.byte("\n")
				end
			`,
			goFunc: func(r rune) bool {
				return r == ' ' || r == '\t' || r == '\r' || r == '\n'
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case=%d/string=%q", i, tt.s), func(t *testing.T) {
			testL := setupLuaFuncTest(t, "TrimRightFunc", tt.luaFunc)
			defer testL.Close()

			expected := strings.TrimRightFunc(tt.s, tt.goFunc)

			testL.Push(testL.GetGlobal("test_TrimRightFunc"))
			testL.Push(lua.LString(tt.s))
			testL.Call(1, 1)

			got := testL.ToString(-1)
			testL.Pop(1)

			require.Equal(t, expected, got,
				"case %d: Lua returned %v but Go returned %v (string: %q, func: %q)",
				i, got, expected, tt.s, tt.luaFunc)
		})
	}
}

func TestTrimSpace(t *testing.T) {
	const luaFuncName = "TrimSpace"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s string
	}{
		{""},
		{" "},
		{"  "},
		{"\t"},
		{"\n"},
		{"\r"},
		{"\v"},
		{"\f"},
		{"hello"},
		{" hello "},
		{"  hello  "},
		{"\thello\t"},
		{"\nhello\n"},
		{"\rhello\r"},
		{"\vhello\v"},
		{"\fhello\f"},
		{" \t\n\r\v\f hello \f\v\r\n\t "},
		{"hello world"},
		{"  hello  world  "},
		{"\u0085hello\u0085"},
		{"\u00A0hello\u00A0"},
		{"你好世界"},
		{"  你好世界  "},
		{"αβγδ"},
		{"  αβγδ  "},
	}

	for i := range tests {
		expected := strings.TrimSpace(tests[i].s)

		args := []lua.LValue{
			lua.LString(tests[i].s),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q)",
			i, got, expected, tests[i].s)
	}
}

func TestTrimSuffix(t *testing.T) {
	const luaFuncName = "TrimSuffix"

	L := setupLuaTest(t, luaFuncName)
	defer L.Close()

	tests := []struct {
		s      string
		suffix string
	}{
		{"", ""},
		{"hello", ""},
		{"", "hello"},
		{"hello", "llo"},
		{"hello", "hello"},
		{"hello", "1hello"},
		{"hello", "world"},
		{"hello world", " world"},
		{"你好世界", "世界"},
		{"αβγδ", "γδ"},
		{"hello\nworld", "\nworld"},
		{"hello  ", "  "},
		{"hello", "o"},
		{"hello\u0000", "\u0000"},
		{"hellosuffix", "suffix"},
	}

	for i := range tests {
		expected := strings.TrimSuffix(tests[i].s, tests[i].suffix)

		args := []lua.LValue{
			lua.LString(tests[i].s),
			lua.LString(tests[i].suffix),
		}
		got := callLuaFunc(t, L, luaFuncName, args, func(L *lua.LState, idx int) string {
			return L.ToString(idx)
		})

		require.Equal(t, expected, got,
			"case %d: Lua returned %q but Go returned %q (string: %q, suffix: %q)",
			i, got, expected, tests[i].s, tests[i].suffix)
	}
}
