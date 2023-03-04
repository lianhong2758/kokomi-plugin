package kokomi

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	//"github.com/fogleman/gg"//原版gg
	"github.com/FloatTech/gg"
)

// Polygon 画多边形
func Polygon(n int) []gg.Point {
	result := make([]gg.Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = gg.Point{X: math.Cos(a), Y: math.Sin(a)}
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

// 绘制带底色的字块,color为字体颜色,x,y左上角位置
func DrawStringRec(dc *gg.Context, str, color string, x, y float64) (w float64) {
	w, h := dc.MeasureString(str)
	dc.DrawRoundedRectangle(x, y, w+8, h+10, 8)
	dc.Fill()
	dc.SetHexColor(color)
	dc.DrawString(str, x+3, h+y+4)
	return
}

// 随机颜色
func randfill() (c [3]int) {
	c[0] = rand.Intn(195) + 15 //r
	c[1] = rand.Intn(195) + 15 //g
	c[2] = rand.Intn(195) + 15 //b
	return
}
