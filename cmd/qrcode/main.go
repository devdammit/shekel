package main

import (
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"os"
)

func main() {
	file, err := os.Open("./qrcode.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	img, _, _ := image.Decode(file)

	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	qrCode, _ := qrReader.Decode(bmp, nil)

	fmt.Println(qrCode.GetBarcodeFormat(), qrCode.GetNumBits())

	encoder := qrcode.NewQRCodeWriter()
	img, _ = encoder.Encode(qrCode.GetText(), gozxing.BarcodeFormat_QR_CODE, 256, 256, nil)

	file, _ = os.Create("./qrcode2.jpg")
	defer file.Close()

	_ = jpeg.Encode(file, img, &jpeg.Options{100})
}
