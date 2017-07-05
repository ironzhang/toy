package report

import (
	"encoding/gob"
	"io"
	"time"
)

type Header struct {
	Name       string
	QPS        int
	Request    int
	Concurrent int
}

type Block struct {
	Time    time.Time
	Records []Record
}

type Encoder interface {
	EncodeHeader(h *Header) error
	EncodeBlock(b *Block) error
}

type Decoder interface {
	DecodeHeader(h *Header) error
	DecodeBlock(b *Block) error
}

type gobEncoder struct {
	enc *gob.Encoder
}

func NewGobEncoder(w io.Writer) Encoder {
	return &gobEncoder{enc: gob.NewEncoder(w)}
}

func (p *gobEncoder) EncodeHeader(h *Header) error {
	return p.enc.Encode(h)
}

func (p *gobEncoder) EncodeBlock(b *Block) error {
	return p.enc.Encode(b)
}

type gobDecoder struct {
	dec *gob.Decoder
}

func NewGobDecoder(r io.Reader) Decoder {
	return &gobDecoder{dec: gob.NewDecoder(r)}
}

func (p *gobDecoder) DecodeHeader(h *Header) error {
	return p.dec.Decode(h)
}

func (p *gobDecoder) DecodeBlock(b *Block) error {
	return p.dec.Decode(b)
}
