// codec pakage defines the interface of codec and the header of message.
package codec

import "io"

// Header represents the header of a message.
// ServerMethod is the service and method to call.
// Seq is the sequence number chosen by the client.
// Error is the error, if any, for the request.
type Header struct {
	ServiceMethod string
	Seq           uint64
	Error         string
}

// Codec is a interface that defines the methods to read and write messages.
// ReadHeader reads the header of a message.
// ReadBody reads the body of a message.
// Write writes the header and body of a message.
// Close closes the codec.
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// NewCodecFunc is a function that creates a new codec.
type NewCodecFunc func(io.ReadWriteCloser) Codec

// Type is the type of the message.
// The client and server can get the constructor through the Type of Codec
type Type string

// The types of Codec
const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

// NewCodecFuncMap is a map from Type to NewCodecFunc.
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec   // User Gob
	NewCodecFuncMap[JsonType] = NewJsonCodec // User Json
}
