package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

var _ Codec = (*JsonCodec)(nil)

func (c *JsonCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *JsonCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// Write writes the header and body of a message.
func (c *JsonCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		// flush the buffer
		// automatically flushes to call sys-writer when buffer is full
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc codec: Json error encoding header: ", err)
		return err
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec: Json error encoding body: ", err)
		return err
	}
	return nil
}

func (c *JsonCodec) Close() error {
	return c.conn.Close()
}
