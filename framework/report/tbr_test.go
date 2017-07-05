package report

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

type TData struct {
	header Header
	blocks []Block
}

func (p *TData) Encode(enc Encoder) (err error) {
	if err = enc.EncodeHeader(&p.header); err != nil {
		return err
	}
	for _, b := range p.blocks {
		if err = enc.EncodeBlock(&b); err != nil {
			return err
		}
	}
	return nil
}

func (p *TData) Decode(dec Decoder) (err error) {
	if err = dec.DecodeHeader(&p.header); err != nil {
		return err
	}
	var b Block
	for {
		if err = dec.DecodeBlock(&b); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		p.blocks = append(p.blocks, b)
	}
	return nil
}

func TestGobCodec(t *testing.T) {
	var a = TData{
		header: Header{
			Name:       "TestGobCodec",
			QPS:        100,
			Request:    200,
			Concurrent: 10,
		},
		blocks: []Block{
			{Total: 100 * time.Second, Records: MakeRandomRecords(10)},
			{Total: 120 * time.Second, Records: MakeRandomRecords(20)},
		},
	}
	var b TData

	var err error
	var buf bytes.Buffer
	enc := NewGobEncoder(&buf)
	dec := NewGobDecoder(&buf)
	if err = a.Encode(enc); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err = b.Decode(dec); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%v != %v", a, b)
	}
}
