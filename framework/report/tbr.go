package report

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"time"
)

type Header struct {
	Time       time.Time
	Name       string
	QPS        int
	Request    int
	Concurrent int
}

func (h *Header) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Time:%s, Name:%s, QPS:%d, Request:%d, Concurrent:%d", h.Time, h.Name, h.QPS, h.Request, h.Concurrent)
	return buf.String()
}

type Block struct {
	Time    time.Time
	Records []Record
}

func (b *Block) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Time:%s, Records:%v", b.Time, b.Records)
	return buf.String()
}

type Encoder interface {
	EncodeHeader(h *Header) error
	EncodeBlock(b *Block) error
}

type Decoder interface {
	DecodeHeader(h *Header) error
	DecodeBlock(b *Block) error
}

func NewEncoder(w io.Writer) Encoder {
	return &gobEncoder{enc: gob.NewEncoder(w)}
}

type gobEncoder struct {
	enc *gob.Encoder
}

func (p *gobEncoder) EncodeHeader(h *Header) error {
	return p.enc.Encode(h)
}

func (p *gobEncoder) EncodeBlock(b *Block) error {
	return p.enc.Encode(b)
}

func NewDecoder(r io.Reader) Decoder {
	return &gobDecoder{dec: gob.NewDecoder(r)}
}

type gobDecoder struct {
	dec *gob.Decoder
}

func (p *gobDecoder) DecodeHeader(h *Header) error {
	return p.dec.Decode(h)
}

func (p *gobDecoder) DecodeBlock(b *Block) error {
	return p.dec.Decode(b)
}
