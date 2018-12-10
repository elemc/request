package request_test

import (
	"fmt"
	"testing"

	"github.com/elemc/request"
)

func testMustInt() error {
	expected := 12345
	value := request.Value("12345")
	result := value.MustInt()
	if result != expected {
		return fmt.Errorf("Unexpected result: %d. Expected: %d", result, expected)
	}
	return nil
}

func TestMustInt(t *testing.T) {
	if err := testMustInt(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMustInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := testMustInt(); err != nil {
			b.Fatal(err)
		}
	}
}

func testMustBool() error {
	expected := true
	value := request.Value("true")
	result := value.MustBool()
	if result != expected {
		return fmt.Errorf("Unexpected result: %t. Expected: %t", result, expected)
	}
	return nil
}

func TestMustBool(t *testing.T) {
	if err := testMustBool(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMustBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := testMustInt(); err != nil {
			b.Fatal(err)
		}
	}
}

func testMustMAC() error {
	expected := "74-e1-b6-6d-1d-58"
	value := request.Value("74:E1:B6:6D:1D:58")
	result := value.MustMAC()
	if result != expected {
		return fmt.Errorf("Unexpected result: %s. Expected: %s", result, expected)
	}
	return nil
}

func TestMustMAC(t *testing.T) {
	if err := testMustBool(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMustMAC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := testMustInt(); err != nil {
			b.Fatal(err)
		}
	}
}
