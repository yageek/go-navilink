package main

//#include <stdio.h>
//#include <string.h>
//
// struct Person {
//  char name[1024];
//  int age;
// };
//
// typedef struct Person Person_t;
//
// void print_person_name(Person_t *person) {
//  printf("[C] Name: %s \n", person->name);
// };
//
// void set_person_name(Person_t *person, const char *name) {
// strcpy(person->name, name);
//};
//
// Person_t boss = (Person_t){.name = "Tristan", .age = 42};
import "C"
import "fmt"

func main() {

	C.print_person_name(&C.boss) // stdout: [C] Name: Tristan

	name := "Linus"
	fmt.Printf("[Go]: Address: %p \n", &name) // stdout: [Go]: Address: 0xc000012050

	nameC := C.CString(name)
	fmt.Printf("[Go]: Address: %p \n", nameC) // stdout: [Go]: Address: 0x4700020

	C.set_person_name(&C.boss, nameC)
	C.print_person_name(&C.boss) // stdout: [C] Name: Linus

	nameGo := C.GoString(nameC)
	fmt.Printf("[Go]: Address: %p \n", &nameGo) // [Go]: Address: 0xc000012060

	p1 := C.struct_Person{age: 10}
	C.set_person_name(&p1, C.CString("Person 1"))

	p2 := C.Person_t{age: 12}
	C.set_person_name(&p2, C.CString("Person 2"))

	C.print_person_name(&p1) // stdout: [C] Name: Person 1

	C.print_person_name(&p2) // stdout: [C] Name: Person 2

}
