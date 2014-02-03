package gouuid_test

import (
	"fmt"
	"github.com/nu7hatch/gouuid"

)

func ExampleNewV4() {
	u4, err := gouuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u4)
}

func ExampleNewV5() {
	u5, err := gouuid.NewV5(gouuid.NamespaceURL, []byte("nu7hat.ch"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u5)
}

func ExampleParseHex() {
	u, err := gouuid.ParseHex("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u)
}
