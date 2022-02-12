package Reflection_test

import (
	"CitadelCore/Shared/Reflection"
	"crypto/rand"
	"fmt"
	"testing"
)

type testpointerstruct struct {
	Name  string
	Value int
	Bytes []byte
}

func TestGetFields(t *testing.T) {
	sut := testpointerstruct{}

	result := Reflection.GetFields(&sut)

	name := result[0].(*string)
	value := result[1].(*int)

	fmt.Printf("Name: %s, Value: %d\n", *name, *value)
	if *name != "" {
		t.Fatalf("Expected empty string, got %s\n", *name)
	}
	if *value != 0 {
		t.Fatalf("Expcted 0, got %d\n", *value)
	}
	// TODO: Test bytes?
}

func TestCreateResultFromFields(t *testing.T) {
	testname := "test"
	testvalue := 42
	testbytes := make([]byte, 42)
	rand.Read(testbytes)
	fields := []interface{}{&testname, &testvalue, &testbytes}

	sut := &testpointerstruct{}
	sut = Reflection.CreateResultlFromFields(fields, sut).(*testpointerstruct)

	if testname != sut.Name {
		t.Fatalf("Expected %s, got %s\n", testname, sut.Name)
	}
	if testvalue != sut.Value {
		t.Fatalf("Expcted %d, got %d\n", testvalue, sut.Value)
	}
}
