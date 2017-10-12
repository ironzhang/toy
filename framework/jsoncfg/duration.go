package jsoncfg

import (
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) String() string {
	return time.Duration(d).String()
}

func (d Duration) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%s", d)
	return []byte(s), nil
}

func (d *Duration) UnmarshalText(b []byte) error {
	du, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	*d = Duration(du)
	return nil
}
