package jsoncfg

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	type Config struct {
		A Duration
		B Duration
	}

	want := Config{
		A: Duration(60 * time.Second),
		B: Duration(90 * time.Minute),
	}
	var got Config

	buf, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	} else {
		t.Logf("buf: %s", buf)
	}
	if err = json.Unmarshal(buf, &got); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if got != want {
		t.Errorf("%v != %v", got, want)
	} else {
		t.Logf("%v == %v", got, want)
	}
}

func TestDurationUnmarshal(t *testing.T) {
	type Config struct {
		A Duration
		B Duration
	}
	var got Config

	buf := []byte(`{"A":"1m0s"}`)
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	} else {
		t.Logf("got: %v", got)
	}
}
