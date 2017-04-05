package microview

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"testing"
)

func TestNewMicroView(t *testing.T) {
	buf := bytes.NewBufferString("")

	mv := NewMicroView(struct {
		io.ReadWriter
		io.Closer
	}{buf, ioutil.NopCloser(nil)}, Delay(0))
	defer mv.Close()

	if got, want := mv.Bounds(), image.Rect(0, 0, 64, 48); !got.Eq(want) {
		t.Fatalf("m.Bounds() = %v, want %v", got, want)
	}

	mv.Run(Rect(5, 10, 15, 20))

	if got, want := buf.String(), "9,5,10,15,20"; got != want {
		t.Fatalf("buf.String() = %q, want %q", got, want)
	}
}

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
