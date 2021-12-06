package utils

import (
	"io"
	"math"
	"strconv"
	"strings"
)

var (
	suffixes [5]string
)

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func ConvertByte(size float64) string {
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}

func CharLimiter(s string, limit int) string {
	reader := strings.NewReader(s)
	buff := make([]byte, limit)
	n, _ := io.ReadAtLeast(reader, buff, limit)

	if n != 0 {
		return string(buff) + "..."
	} else {
		return s
	}
}
