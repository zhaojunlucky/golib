package text

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"testing"
)

func Utf8ToEnc(str string, enc encoding.Encoding) []byte {
	r := transform.NewReader(bytes.NewReader([]byte(str)), enc.NewEncoder())
	b, _ := ioutil.ReadAll(r)
	return b
}

func TestDetectFileEncoding(t *testing.T) {
	var encDetect = NewBytesEncodingDetect()
	var bytesArr = [][]byte{
		[]byte("/Users/jun/Documents/magicworldz/github/go-lib/pkg/text/bytes_encoding_detect_test.go"),
		Utf8ToEnc("哈啥事", simplifiedchinese.GB18030),
		Utf8ToEnc("哈啥事", simplifiedchinese.GBK),
		Utf8ToEnc("哈啥事", simplifiedchinese.HZGB2312),
		Utf8ToEnc("哈啥事", traditionalchinese.Big5),
	}
	for _, bytes := range bytesArr {
		enc := encDetect.DetectBytesEncoding(bytes)

		fmt.Printf("%v\n", enc)
	}

}
