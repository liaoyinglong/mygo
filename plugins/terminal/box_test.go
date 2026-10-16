package terminal

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestRoundedPowerlineCaps(t *testing.T) {
	for _, height := range []float32{12, 16, 20} {
		for _, scale := range []float32{1, 2} {
			for _, r := range []rune{0xe0b4, 0xe0b5, 0xe0b6, 0xe0b7} {
				t.Run(fmt.Sprintf("%04X/height=%g/scale=%g", r, height, scale), func(t *testing.T) {
					if !drawn(r) {
						t.Fatalf("U+%04X uses font rendering instead of cell geometry", r)
					}
					cell := ui.Rect{X: 4, Y: 4, W: 10, H: height}
					tt := ui.NewTester(func(c *ui.Context) {
						ui.Box(c).Fill().Draw(func(p *ui.Painter, _ ui.Rect) {
							p.Fill(ui.Rect{W: 20, H: height + 8}, ui.RGB(0, 0, 0), 0)
							drawBox(p, r, cell, ui.RGB(255, 255, 255))
						})
					}, 20, int(height)+8)
					tt.SetScale(scale)
					img := tt.Image()
					x0, y0 := int(cell.X*scale), int(cell.Y*scale)
					x1, y1 := int((cell.X+cell.W)*scale), int((cell.Y+cell.H)*scale)
					right := r == 0xe0b4 || r == 0xe0b5
					solid := r == 0xe0b4 || r == 0xe0b6
					edge, tip := x0, x1-1
					if !right {
						edge, tip = tip, edge
					}
					for y := y0; y < y1; y++ {
						// Allow antialiasing where the curve meets the top and bottom.
						if solid && img.RGBAAt(edge, y).R < 230 {
							t.Errorf("gap at background join (%d, %d): %v", edge, y, img.RGBAAt(edge, y))
						}
					}
					if img.RGBAAt(tip, (y0+y1)/2).R == 0 {
						t.Error("cap does not reach the outer edge at mid-height")
					}
					center := img.RGBAAt((x0+x1)/2, (y0+y1)/2).R
					if solid && center < 240 || !solid && center != 0 {
						t.Errorf("cap center = %d, solid = %v", center, solid)
					}
					for y := range img.Bounds().Dy() {
						for x := range img.Bounds().Dx() {
							if (x < x0 || x >= x1 || y < y0 || y >= y1) && img.RGBAAt(x, y).R != 0 {
								t.Fatalf("cap spills outside its cell at (%d, %d)", x, y)
							}
						}
					}
					if img.RGBAAt(tip, y0).R != 0 || img.RGBAAt(tip, y1-1).R != 0 {
						t.Error("outer corners must remain empty")
					}
				})
			}
		}
	}
}

func TestPowerlineDrawingRange(t *testing.T) {
	for r := rune(0xe0b0); r <= 0xe0b7; r++ {
		if !drawn(r) {
			t.Errorf("U+%04X is not drawn", r)
		}
	}
	for _, r := range []rune{0xe0af, 0xe0b8, 0xf07b, 0xf120} {
		if drawn(r) {
			t.Errorf("U+%04X must stay on the font path", r)
		}
	}
}
