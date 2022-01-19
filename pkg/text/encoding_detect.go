package text

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

const (
	GB2312        = 0
	GBK           = 1
	GB18030       = 2
	HZ            = 3
	BIG5          = 4
	CNS11643      = 5
	UTF8          = 6
	UTF8T         = 7
	UTF8S         = 8
	UNICODE       = 9
	UNICODET      = 10
	UNICODES      = 11
	ISO2022CN     = 12
	ISO2022CN_CNS = 13
	ISO2022CN_GB  = 14
	EUC_KR        = 15
	CP949         = 16
	ISO2022KR     = 17
	JOHAB         = 18
	SJIS          = 19
	EUC_JP        = 20
	ISO2022JP     = 21
	ASCII         = 22
	OTHER         = 23
	TOTALTYPES    = 24
)

type EncodingDetect struct {
	// Names of the encodings as understood by Java
	gonames []encoding.Encoding

	// Names of the encodings for human viewing
	nicename []string

	// Names of charsets as used in charset parameter of HTML Meta tag
	htmlname []string
}

func NewEncoding() *EncodingDetect {
	encoding := &EncodingDetect{}
	encoding.initialize()
	return encoding
}

func (enc *EncodingDetect) initialize() {
	enc.gonames = make([]encoding.Encoding, TOTALTYPES)
	enc.nicename = make([]string, TOTALTYPES)
	enc.htmlname = make([]string, TOTALTYPES)
	// Assign encoding names
	enc.gonames[GB2312] = simplifiedchinese.HZGB2312
	enc.gonames[GBK] = simplifiedchinese.GBK
	enc.gonames[GB18030] = simplifiedchinese.GB18030
	enc.gonames[HZ] = encoding.Nop
	enc.gonames[ISO2022CN_GB] = encoding.Nop
	enc.gonames[BIG5] = traditionalchinese.Big5
	enc.gonames[CNS11643] = encoding.Nop
	enc.gonames[ISO2022CN_CNS] = encoding.Nop
	enc.gonames[ISO2022CN] = encoding.Nop
	enc.gonames[UTF8] = unicode.UTF8
	enc.gonames[UTF8T] = unicode.UTF8
	enc.gonames[UTF8S] = unicode.UTF8
	enc.gonames[UNICODE] = encoding.Nop
	enc.gonames[UNICODET] = encoding.Nop
	enc.gonames[UNICODES] = encoding.Nop
	enc.gonames[EUC_KR] = korean.EUCKR
	enc.gonames[CP949] = encoding.Nop
	enc.gonames[ISO2022KR] = korean.EUCKR
	enc.gonames[JOHAB] = encoding.Nop
	enc.gonames[SJIS] = japanese.ShiftJIS
	enc.gonames[EUC_JP] = japanese.EUCJP
	enc.gonames[ISO2022JP] = japanese.ISO2022JP
	enc.gonames[ASCII] = charmap.ISO8859_1
	enc.gonames[OTHER] = charmap.ISO8859_1
	// Assign encoding names
	enc.htmlname[GB2312] = "GB2312"
	enc.htmlname[GBK] = "GBK"
	enc.htmlname[GB18030] = "GB18030"
	enc.htmlname[HZ] = "HZ-GB-2312"
	enc.htmlname[ISO2022CN_GB] = "ISO-2022-CN-EXT"
	enc.htmlname[BIG5] = "BIG5"
	enc.htmlname[CNS11643] = "EUC-TW"
	enc.htmlname[ISO2022CN_CNS] = "ISO-2022-CN-EXT"
	enc.htmlname[ISO2022CN] = "ISO-2022-CN"
	enc.htmlname[UTF8] = "UTF-8"
	enc.htmlname[UTF8T] = "UTF-8"
	enc.htmlname[UTF8S] = "UTF-8"
	enc.htmlname[UNICODE] = "UTF-16"
	enc.htmlname[UNICODET] = "UTF-16"
	enc.htmlname[UNICODES] = "UTF-16"
	enc.htmlname[EUC_KR] = "EUC-KR"
	enc.htmlname[CP949] = "x-windows-949"
	enc.htmlname[ISO2022KR] = "ISO-2022-KR"
	enc.htmlname[JOHAB] = "x-Johab"
	enc.htmlname[SJIS] = "Shift_JIS"
	enc.htmlname[EUC_JP] = "EUC-JP"
	enc.htmlname[ISO2022JP] = "ISO-2022-JP"
	enc.htmlname[ASCII] = "ASCII"
	enc.htmlname[OTHER] = "ISO8859-1"
	// Assign Human readable names
	enc.nicename[GB2312] = "GB-2312"
	enc.nicename[GBK] = "GBK"
	enc.nicename[GB18030] = "GB18030"
	enc.nicename[HZ] = "HZ"
	enc.nicename[ISO2022CN_GB] = "ISO2022CN-GB"
	enc.nicename[BIG5] = "Big5"
	enc.nicename[CNS11643] = "CNS11643"
	enc.nicename[ISO2022CN_CNS] = "ISO2022CN-CNS"
	enc.nicename[ISO2022CN] = "ISO2022 CN"
	enc.nicename[UTF8] = "UTF-8"
	enc.nicename[UTF8T] = "UTF-8 (Trad)"
	enc.nicename[UTF8S] = "UTF-8 (Simp)"
	enc.nicename[UNICODE] = "Unicode"
	enc.nicename[UNICODET] = "Unicode (Trad)"
	enc.nicename[UNICODES] = "Unicode (Simp)"
	enc.nicename[EUC_KR] = "EUC-KR"
	enc.nicename[CP949] = "CP949"
	enc.nicename[ISO2022KR] = "ISO 2022 KR"
	enc.nicename[JOHAB] = "Johab"
	enc.nicename[SJIS] = "Shift-JIS"
	enc.nicename[EUC_JP] = "EUC-JP"
	enc.nicename[ISO2022JP] = "ISO 2022 JP"
	enc.nicename[ASCII] = "ASCII"
	enc.nicename[OTHER] = "OTHER"
}
