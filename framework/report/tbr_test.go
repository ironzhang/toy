package report

import (
	"bytes"
	"fmt"
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
	for {
		var b Block
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

func (p *TData) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Header: %s\n", p.header.String())
	for i, b := range p.blocks {
		fmt.Fprintf(&buf, "Block[%d]: %s\n", i, b.String())
	}
	return buf.String()
}

func TestGobCodec(t *testing.T) {
	ts := time.Now()
	var a = TData{
		header: Header{
			Name:       "TestGobCodec",
			QPS:        100,
			Request:    200,
			Concurrent: 10,
		},
		blocks: []Block{
			{Time: ts.Add(100 * time.Second), Records: MakeRandomRecords(1)},
			{Time: ts.Add(120 * time.Second), Records: MakeRandomRecords(10)},
		},
	}
	var b TData

	var err error
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	dec := NewDecoder(&buf)
	if err = a.Encode(enc); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err = b.Decode(dec); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%s != %s", a.String(), b.String())
	}
}
