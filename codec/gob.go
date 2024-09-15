package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// 为了调用 init 函数
var _ Codec = (*GobCodec)(nil)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g *GobCodec) ReadBody(i interface{}) error {
	return g.dec.Decode(i)
}

func (g *GobCodec) Write(header *Header, i interface{}) (err error) {
	// 最后关闭
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()

	if err = g.enc.Encode(header); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = g.enc.Encode(i); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}

	return
}
