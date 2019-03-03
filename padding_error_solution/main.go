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


#if __BYTE_ORDER__ == __ORDER_LITTLE_ENDIAN__
const uint8_t endianness = 0;
#else
const uint8_t endianness = 255;
#endif

struct NonPadded SomeNonPaddedInstance = {.value1 = 2, .value2 = 4, .value3 = 56, .value4 = 255};
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"fmt"
)

type PaddingSolution struct {
	Value1 uint16
	Value2 uint32
	Value3 uint64
	Value4 uint8
}

func NewPaddingSolution(value *C.struct_NonPadded) (PaddingSolution, error) {

	ptr := unsafe.Pointer(value)
	buff := C.GoBytes(ptr, 15)
	buffer := bytes.NewBuffer(buff)

	var orderByte binary.ByteOrder
	if C.endianness == 0 {
		orderByte = binary.LittleEndian
	} else {
		orderByte = binary.BigEndian
	}

	var value1 uint16
	if err := binary.Read(buffer, orderByte, &value1); err != nil {
		return PaddingSolution{}, err
	}

	var value2 uint32
	if err := binary.Read(buffer, orderByte, &value2); err != nil {
		return PaddingSolution{}, err
	}

	var value3 uint64
	if err := binary.Read(buffer, orderByte, &value3); err != nil {
		return PaddingSolution{}, err
	}

	var value4 uint8 = buff[len(buff)-1]

	return PaddingSolution{
		Value1: value1,
		Value2: value2,
		Value3: value3,
		Value4: value4,
	}, nil
}

func main() {
	fmt.Printf("Taille de Padded: %d octet(s)\n", C.sizeof_Padded)
	fmt.Printf("Taille de NonPadded: %d octet(s)\n", C.sizeof_struct_NonPadded)
	a := C.Padded{
		value1: 1,
		value2: 2,
		value3: 3,
		value4: 4,
	}

	fmt.Printf("a: %#v \n", a)

	// Récupération des valeurs de `SomeNonPaddedInstance` depuis `unsafe.Pointer`
	if solution, err := NewPaddingSolution(&C.SomeNonPaddedInstance); err != nil {
		fmt.Printf("Erreur: %v\n", err)
	} else {
		fmt.Printf("NewPaddingSolution: %#v \n", solution)
	}

}
