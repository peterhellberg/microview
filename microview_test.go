package microview

import (
	"bytes"
	"testing"
)

func TestClear(t *testing.T) {
	for _, tt := range []struct {
		name string
		mode uint8
		want []byte
	}{
		{"ALL", ALL, []byte("0,1")},
		{"PAGE", PAGE, []byte("0,0")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clear(tt.mode); !bytes.Equal(got, tt.want) {
				t.Fatalf("Clear(%s) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
