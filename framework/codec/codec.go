package codec

import "encoding/gob"

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}

//var NewEncoder = json.NewEncoder
//var NewDeocder = json.NewDecoder

var NewEncoder = gob.NewEncoder
var NewDecoder = gob.NewDecoder
