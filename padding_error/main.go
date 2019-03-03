package main

/*
#include <stdio.h>
#include <stdint.h>

typedef struct {
    uint16_t value1;
    uint32_t value2;
    uint64_t value3;
    uint8_t value4;
} Padded;

struct NonPadded{
    uint16_t value1;
    uint32_t value2;
    uint64_t value3;
    uint8_t value4;
} __attribute__((packed));
*/
import "C"

import "fmt"

func main() {
	fmt.Printf("Taille de Padded: %d octet(s)\n", C.sizeof_Padded)
	fmt.Printf("Taille de NonPadded: %d octet(s)\n", C.sizeof_struct_NonPadded)
	a := C.Padded{
		value1: 1,
		value2: 2,
		value3: 3,
		value4: 4,
	}

	b := C.struct_NonPadded{
		value1: 1,
		value2: 2,
		value3: 3,
		value4: 4,
	}

	fmt.Printf("a: %#v \n", a)
	fmt.Printf("b: %#v \n", b)

}
