package main

import (
	"fmt"

	"github.com/skeptycal/stringpool"
)

func main() {
	sb := stringpool.Get()
	defer stringpool.Release(sb)

	for i := 0; i < 255; i++ {

		// this is just an example ...
		//
		// using fmt.Sprintf is slow and mostly defeats
		// the purpose of strings.Builder objects ...
		//
		// but it does provide a familiar example context ...
		s := fmt.Sprintf("%d + %q\n", i, i)
		sb.WriteString(s)
	}

	fmt.Print("stringpool example:\n\n")
	fmt.Print(sb)
	fmt.Print("\nstringpool example\n\n")

}
