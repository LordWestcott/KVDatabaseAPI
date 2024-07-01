package util

import (
	"bytes"
	"io"
	"testing"
)

func TestStreamToByte_NonEmptyStream(t *testing.T) {
	input := "Hello, World!"
	stream := bytes.NewBufferString(input)

	result, err := StreamToByte(stream)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(result) != input {
		t.Errorf("Expected %s, got %s", input, string(result))
	}
}

func TestStreamToByte_EmptyStream(t *testing.T) {
	var stream io.Reader = bytes.NewReader([]byte{})

	result, err := StreamToByte(stream)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty byte slice, got %v", result)
	}
}
