package terminal

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

// The render-state path must cache caps as geometry, not font glyphs.
func TestRoundedPowerlineTerminal(t *testing.T) {
	loadLib(t)
	for _, size := range []float32{12, 16, 20} {
		for _, scale := range []float32{1, 2} {
			t.Run(fmt.Sprintf("size=%g/scale=%g", size, scale), func(t *testing.T) {
				term, err := New(Options{Conn: newPipe(), Font: Font{Size: size}})
				if err != nil {
					t.Fatal(err)
				}
				defer term.Close()
				term.Feed([]byte("\x1b[?25l\x1b[38;2;255;255;255m\ue0b6\x1b[48;2;255;255;255m   \x1b[49m\ue0b4\x1b[0m\r\n\ue0b7   \ue0b5"))
				tt := ui.NewTester(func(c *ui.Context) { View(c, term).Fill() }, 200, 100)
				tt.SetScale(scale)
				img := tt.Image()
				v := term.v
				for y, want := range [][2]rune{{0xe0b6, 0xe0b4}, {0xe0b7, 0xe0b5}} {
					line := v.lines[y]
					if len(line.boxes) != 2 || line.boxes[0].r != want[0] || line.boxes[1].r != want[1] || line.boxes[0].x != 0 || line.boxes[1].x != 4 {
						t.Fatalf("row %d caps = %v, want %U in columns 0 and 4", y, line.boxes, want)
					}
					if len(line.runs) != 0 {
						t.Errorf("row %d still shapes caps as font glyphs", y)
					}
				}
				// Both solid caps meet the neighboring white background at every row.
				for y := v.oy; y < v.oy+v.cellH; y++ {
					for _, x := range []int{v.ox + v.cellW - 1, v.ox + 4*v.cellW} {
						if got := img.RGBAAt(x, y).R; got < 230 {
							t.Errorf("gap at (%d,%d): brightness %d", x, y, got)
						}
					}
				}
			})
		}
	}
}
