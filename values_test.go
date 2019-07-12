package request_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elemc/request"
)

func testMustInt() error {
	expected := 12345
	value := request.Value("12345")
	result := value.MustInt()
	if result != expected {
		return fmt.Errorf("unexpected result: %d. Expected: %d", result, expected)
	}

	resultUint64 := value.MustUint64()
	if resultUint64 != uint64(expected) {
		return fmt.Errorf("unexpected result (uint64): %d. Expected: %d", resultUint64, expected)
	}

	resultInt64 := value.MustInt64()
	if resultInt64 != int64(expected) {
		return fmt.Errorf("unexpected result (int64): %d. Expected: %d", resultInt64, expected)
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
	value := request.Value("true")
	result := value.MustBool()
	if !result {
		return fmt.Errorf("unexpected result: %t. Expected: true", result)
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
		return fmt.Errorf("unexpected result: %s. Expected: %s", result, expected)
	}
	return nil
}

func testNullMAC() error {
	expected := ""
	value := request.Value("00:00:00:00:00:00")
	result := value.MustMAC()
	if result != expected {
		return fmt.Errorf("unexpected result: %s. Expected: empty string", result)
	}
	return nil
}

func TestMustMAC(t *testing.T) {
	if err := testMustBool(); err != nil {
		t.Fatal(err)
	}
	if err := testNullMAC(); err != nil {
		t.Fatal(err)
	}
	if err := testMustMAC(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMustMAC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := testMustInt(); err != nil {
			b.Fatal(err)
		}
		if err := testNullMAC(); err != nil {
			b.Fatal(err)
		}
		if err := testMustMAC(); err != nil {
			b.Fatal(err)
		}
	}
}

func testMustTime() error {
	// тест на RFC3389
	expected := time.Date(2019, 3, 1, 12, 9, 0, 0, time.Local)
	value := request.Value("2019-03-01T12:09:00+03:00")
	result := value.MustTime()
	if !expected.Equal(result) {
		return fmt.Errorf("unexpected result: %s. Expected %s", result, expected)
	}

	value = request.Value("1551431340")
	result = value.MustTime()
	if !expected.Equal(result) {
		return fmt.Errorf("unexpected result: %s. Expected %d", result, expected.Unix())
	}

	// тест на YYYY-MM-DD
	loc, _ := time.LoadLocation("UTC")
	expected = time.Date(2019, 3, 1, 0, 0, 0, 0, loc)
	value = request.Value("2019-03-01")
	result = value.MustTime()
	if !expected.Equal(result) {
		return fmt.Errorf("unexpected result: %s. Expected %s", result, expected)
	}

	return nil
}

func TestMustTime(t *testing.T) {
	if err := testMustTime(); err != nil {
		t.Fatal(err)
	}
}
