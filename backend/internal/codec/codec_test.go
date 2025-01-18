package codec

import (
	"bytes"
	"io"
	"testing"
)

func TestJsonCodec(t *testing.T) {
	var buf bytes.Buffer
	conn := struct {
		io.Reader
		io.Writer
		io.Closer
	}{&buf, &buf, io.NopCloser(nil)}

	codec := NewJsonCodec(conn)

	header := &Header{ServiceMethod: "Test.Method", Seq: 1}
	body := map[string]string{"key": "value"}

	// Test Write
	if err := codec.Write(header, body); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Test ReadHeader
	readHeader := &Header{}
	if err := codec.ReadHeader(readHeader); err != nil {
		t.Fatalf("ReadHeader failed: %v", err)
	}
	if readHeader.ServiceMethod != header.ServiceMethod || readHeader.Seq != header.Seq {
		t.Fatalf("ReadHeader got %v, want %v", readHeader, header)
	}

	// Test ReadBody
	readBody := make(map[string]string)
	if err := codec.ReadBody(&readBody); err != nil {
		t.Fatalf("ReadBody failed: %v", err)
	}
	if readBody["key"] != body["key"] {
		t.Fatalf("ReadBody got %v, want %v", readBody, body)
	}
}

func TestGobCodec(t *testing.T) {
	var buf bytes.Buffer
	conn := struct {
		io.Reader
		io.Writer
		io.Closer
	}{&buf, &buf, io.NopCloser(nil)}

	codec := NewGobCodec(conn)

	header := &Header{ServiceMethod: "Test.Method", Seq: 1}
	body := map[string]string{"key": "value"}

	// Test Write
	if err := codec.Write(header, body); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Test ReadHeader
	readHeader := &Header{}
	if err := codec.ReadHeader(readHeader); err != nil {
		t.Fatalf("ReadHeader failed: %v", err)
	}
	if readHeader.ServiceMethod != header.ServiceMethod || readHeader.Seq != header.Seq {
		t.Fatalf("ReadHeader got %v, want %v", readHeader, header)
	}

	// Test ReadBody
	readBody := make(map[string]string)
	if err := codec.ReadBody(&readBody); err != nil {
		t.Fatalf("ReadBody failed: %v", err)
	}
	if readBody["key"] != body["key"] {
		t.Fatalf("ReadBody got %v, want %v", readBody, body)
	}
}
