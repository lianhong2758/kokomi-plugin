package kokomi

import (
	"image"
	"image/color"
	m "math"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/zbputils/img"
)

// Polygon 画多边形
func Polygon(n int) []gg.Point {
	result := make([]gg.Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*m.Pi/float64(n) - m.Pi/2
		result[i] = gg.Point{m.Cos(a), m.Sin(a)}
	}
	return result
}

// Drawstars 画星星
func Drawstars(side, all string, num int) image.Image {
	dc := gg.NewContext(500, 80)
	n := 5
	points := Polygon(n)
	for x, i := 40, 0; i < num; x += 80 {
		dc.Push()
		//s := rand.Float64()*S/4 + S/4
		dc.Translate(float64(x), 45)
		//	dc.Rotate(rand.Float64() * 1.5 * math.Pi) //旋转
		dc.Scale(30, 30) //大小
		for i := 0; i < n+1; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}
		dc.SetLineWidth(10)
		dc.SetHexColor(side) //线
		dc.StrokePreserve()
		dc.SetHexColor(all)
		dc.Fill()
		dc.Pop()
		i++
	}
	return dc.Image()
}

// AdjustOpacity 更改透明度
func AdjustOpacity(m image.Image, percentage float64) image.Image {
	bounds := m.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()
	nimg := image.NewRGBA64(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			r, g, b, a := m.At(i, j).RGBA()
			opacity := uint16(float64(a) * percentage)
			r, g, b, a = nimg.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: opacity}).RGBA()
			nimg.SetRGBA64(i, j, color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
		}
	}
	return nimg
}

// Yinying 绘制阴影 圆角矩形
func Yinying(x int, y int, r float64, c color.Color) image.Image {
	ctx := gg.NewContext(x, y)
	ctx.SetColor(c)
	ctx.DrawRoundedRectangle(0, 0, float64(x), float64(y), r)
	ctx.Fill()
	return ctx.Image()
}

// SetMark 绘制马赛克
func SetMark(pic image.Image) (picture []byte) {
	dst := img.Size(pic, 256*5, 256*5)
	b := dst.Im.Bounds()
	markSize := 32

	for y0fMarknum := 0; y0fMarknum <= math.Ceil(b.Max.Y, markSize); y0fMarknum++ {
		for x0fMarknum := 0; x0fMarknum <= math.Ceil(b.Max.X, markSize); x0fMarknum++ {
			a := dst.Im.At(x0fMarknum*markSize+markSize/2, y0fMarknum*markSize+markSize/2)
			cc := color.NRGBAModel.Convert(a).(color.NRGBA)
			for y := 0; y < markSize; y++ {
				for x := 0; x < markSize; x++ {
					x0fPic := x0fMarknum*markSize + x
					y0fPic := y0fMarknum*markSize + y
					dst.Im.Set(x0fPic, y0fPic, cc)
				}
			}
		}
	}
	picture, cl := writer.ToBytes(dst.Im)
	defer cl()
	return
}
