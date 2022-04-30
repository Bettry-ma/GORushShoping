package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewEqual(t *testing.T) {
	// Different allocations should not be equal.
	// 虽然都是errors的包new的相同的字符串,但是这两个error对象并不相等,error底层的比较是对象地址,而不是简单的字符串等值
	if errors.New("abc") == errors.New("abc") {
		t.Errorf(`New("abc") == New("abc")`)
	}
	fmt.Println("this is stdout")
	//error的比较和字符串无关
	if errors.New("abc") == errors.New("xyz") {
		t.Errorf(`New("abc") == New("xyz")`)
	}

	// Same allocation should be equal to itself (not crash).相同的对象应等于自己
	err := errors.New("jkl")
	if err != err {
		t.Errorf(`err != err`)
	}

}

func TestErrorMethod(t *testing.T) {
	err := errors.New("abc")
	// 同一个对象,调用字符串比较时会相等,即errors.New("xxx")==err.Error()
	if err.Error() != "abc" {
		t.Errorf(`New("abc").Error() = %q, want %q`, err.Error(), "abc")
	}
}

func ExampleNew() {
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}
	// Output: emit macho dwarf: elf header corrupted
}

// The fmt package's Errorf function lets us use the package's formatting
// features to create descriptive error messages.
func ExampleNew_errorf() {
	const name, id = "bimmler", 17
	err := fmt.Errorf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Print(err)
	}
	// Output: user "bimmler" (id 17) not found
}
