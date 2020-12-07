package uitls

import (
	"fmt"
	"testing"
)

func Test_copy(t *testing.T) {
	a := A{Name: "zs", Age: 18, C: C{KK: "123"}}
	b := B{}
	err := Copy(&a, &b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
	fmt.Println(b)
	b.KK = ""
	fmt.Println(a)
	fmt.Println(b)
}

type A struct {
	Age  int64
	Name string
	C
}

type B struct {
	Name string
	Age  int
	Sex  int
	C
}

type C struct {
	KK string
}
