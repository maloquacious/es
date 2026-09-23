package es_test

import (
	"fmt"

	"github.com/maloquacious/es"
)

func ExampleCrockford() {
	encoded := es.Crockford.EncodeToString([]byte("foobar"))
	decoded, err := es.Crockford.DecodeString("csqpy-rkle8")
	fmt.Println(encoded)
	fmt.Printf("%s, %v\n", decoded, err)
	// Output:
	// CSQPYRK1E8
	// foobar, <nil>
}

func ExampleZBase32() {
	encoded := es.ZBase32.EncodeToString([]byte("hello"))
	decoded, err := es.ZBase32.DecodeString("PB1SA5DX")
	fmt.Println(encoded)
	fmt.Printf("%s, %v\n", decoded, err)
	// Output:
	// pb1sa5dx
	// hello, <nil>
}
