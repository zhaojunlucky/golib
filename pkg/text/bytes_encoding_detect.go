package text

import (
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"golang.org/x/text/encoding"
	"log"
	"os"
)

type BytesEncodingDetect struct {
	GBFreq    [][]int
	GBKFreq   [][]int
	Big5Freq  [][]int
	Big5PFreq [][]int
	EucTwfreq [][]int
	KRFreq    [][]int
	JPFreq    [][]int
	EncodingDetect
}

func NewBytesEncodingDetect() *BytesEncodingDetect {
	detect := &BytesEncodingDetect{
		EncodingDetect: *NewEncoding(),
	}
	detect.initialize()
	return detect
}

func (detect *BytesEncodingDetect) initialize() {
	detect.GBFreq = collection.CreateTwoDimArray(94, 94)
	detect.GBKFreq = collection.CreateTwoDimArray(126, 191)
	detect.Big5Freq = collection.CreateTwoDimArray(94, 158)
	detect.Big5PFreq = collection.CreateTwoDimArray(126, 191)
	detect.EucTwfreq = collection.CreateTwoDimArray(94, 94)
	detect.KRFreq = collection.CreateTwoDimArray(94, 94)
	detect.JPFreq = collection.CreateTwoDimArray(94, 94)

	detect.initializeFrequencies()
}

/**
 * Function : detectEncoding Aruguments: File Returns : One of the encodings
 * from the EncodingDetect enumeration (GB2312, HZ, BIG5, EUC_TW, ASCII, or OTHER)
 * Description: This function looks at the file and assigns it a probability
 * score for each encoding type. The encoding type with the highest
 * probability is returned.
 */
func (detect *BytesEncodingDetect) DetectFileEncoding(file string) (encoding.Encoding, error) {
	rawtext, err := detect.getFileBytes(file)
	if err != nil {
		return encoding.Nop, err
	}

	return detect.DetectBytesEncoding(rawtext), nil
}

func (detect *BytesEncodingDetect) DetectBytesEncoding(rawtext []byte) encoding.Encoding {
	bytesInt := detect.convertBytes2Int(rawtext)
	encCode := detect.DetectEncoding(bytesInt)
	enc := detect.gonames[encCode]
	if enc == encoding.Nop {
		log.Println("Warn: unable to detect encoding return null encoder")
	}
	return enc
}

func (detect *BytesEncodingDetect) getFileBytes(file string) ([]byte, error) {
	rawtext := make([]byte, 2000)
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("read fail")
		return nil, err
	}

	defer f.Close()

	_, err = f.Read(rawtext)
	if err != nil {
		return nil, err
	}
	return rawtext, nil
}

func (detect *BytesEncodingDetect) convertBytes2Int(bytes []byte) []int {
	rawtextInt := make([]int, 2000)
	for i := 0; i < len(bytes); i++ {
		b := int(bytes[i])
		if b > 127 {
			b -= 256
		}
		rawtextInt[i] = b
	}
	return rawtextInt
}

/**
 * Function : detectEncoding Aruguments: byte array Returns : One of the
 * encodings from the EncodingDetect enumeration (GB2312, HZ, BIG5, EUC_TW, ASCII,
 * or OTHER) Description: This function looks at the byte array and assigns
 * it a probability score for each encoding type. The encoding type with the
 * highest probability is returned.
 */
func (detect *BytesEncodingDetect) DetectEncoding(rawtext []int) int {
	scores := make([]int, TOTALTYPES)
	var maxscore int
	var encoding_guess = OTHER
	// Assign Scores
	scores[GB2312] = detect.gb2312_probability(rawtext)
	scores[GBK] = detect.gbk_probability(rawtext)
	scores[GB18030] = detect.gb18030_probability(rawtext)
	scores[HZ] = detect.hz_probability(rawtext)
	scores[BIG5] = detect.big5_probability(rawtext)
	scores[CNS11643] = detect.euc_tw_probability(rawtext)
	scores[ISO2022CN] = detect.iso_2022_cn_probability(rawtext)
	scores[UTF8] = detect.utf8_probability(rawtext)
	scores[UNICODE] = detect.utf16_probability(rawtext)
	scores[EUC_KR] = detect.euc_kr_probability(rawtext)
	scores[CP949] = detect.cp949_probability(rawtext)
	scores[JOHAB] = 0
	scores[ISO2022KR] = detect.iso_2022_kr_probability(rawtext)
	scores[ASCII] = detect.ascii_probability(rawtext)
	scores[SJIS] = detect.sjis_probability(rawtext)
	scores[EUC_JP] = detect.euc_jp_probability(rawtext)
	scores[ISO2022JP] = detect.iso_2022_jp_probability(rawtext)
	scores[UNICODET] = 0
	scores[UNICODES] = 0
	scores[ISO2022CN_GB] = 0
	scores[ISO2022CN_CNS] = 0
	scores[OTHER] = 0
	// Tabulate Scores
	for index := 0; index < TOTALTYPES; index++ {
		// if (debug)
		//     System.err.println("EncodingDetect " + nicename[index] + " score "
		//             + scores[index]);
		if scores[index] > maxscore {
			encoding_guess = index
			maxscore = scores[index]
		}
	}
	// Return OTHER if nothing scored above 50
	if maxscore <= 50 {
		encoding_guess = OTHER
	}
	return encoding_guess
}

/*
 * Function: gb2312_probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses GB-2312
 * encoding
 */
func (detect *BytesEncodingDetect) gb2312_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var gbchars = 1
	var gbfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xF7-256) && (0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
				gbchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				column = rawtext[i+1] + 256 - 0xA1
				if detect.GBFreq[row][column] != 0 {
					gbfreq += int64(detect.GBFreq[row][column])
				} else if 15 <= row && row < 55 {
					// In GB high-freq character range
					gbfreq += 200
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(gbchars)/float64(dbchars))
	freqval = 50 * float32(float64(gbfreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

/*
 * Function: gbk_probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses GBK
 * encoding
 */
func (detect *BytesEncodingDetect) gbk_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var gbchars = 1
	var gbfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xF7-256) && // Original GB range
				(0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
				gbchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				column = rawtext[i+1] + 256 - 0xA1
				// System.out.println("original row " + row + " column " +
				// column);
				if detect.GBFreq[row][column] != 0 {
					gbfreq += int64(detect.GBFreq[row][column])
				} else if 15 <= row && row < 55 {
					gbfreq += 200
				}
			} else if (0x81-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && // Extended GB range
				(((0x80-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256)) || ((0x40-256) <= rawtext[i+1] && rawtext[i+1] <= (0x7E-256))) {
				gbchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0x81
				if 0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E {
					column = rawtext[i+1] - 0x40
				} else {
					column = rawtext[i+1] + 256 - 0x40
				}
				// System.out.println("extended row " + row + " column " +
				// column + " rawtext[i] " + rawtext[i]);
				if detect.GBKFreq[row][column] != 0 {
					gbfreq += int64(detect.GBKFreq[row][column])
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(gbchars)/float64(dbchars))
	freqval = 50 * float32(float64(gbfreq)/float64(totalfreq))
	// For regular GB files, this would give the same score, so I handicap
	// it slightly
	return (int)(rangeval+freqval) - 1
}

/*
 * Function: gb18030_probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses GBK
 * encoding
 */
func (detect *BytesEncodingDetect) gb18030_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var gbchars = 1
	var gbfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xF7-256) && // Original GB range
				i+1 < rawtextlen && (0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
				gbchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				column = rawtext[i+1] + 256 - 0xA1
				// System.out.println("original row " + row + " column " +
				// column);
				if detect.GBFreq[row][column] != 0 {
					gbfreq += int64(detect.GBFreq[row][column])
				} else if 15 <= row && row < 55 {
					gbfreq += 200
				}
			} else if (0x81-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && // Extended GB range
				i+1 < rawtextlen && (((0x80-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256)) || ((0x40-256) <= rawtext[i+1] && rawtext[i+1] <= (0x7E-256))) {
				gbchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0x81
				if 0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E {
					column = rawtext[i+1] - 0x40
				} else {
					column = rawtext[i+1] + 256 - 0x40
				}
				// System.out.println("extended row " + row + " column " +
				// column + " rawtext[i] " + rawtext[i]);
				if detect.GBKFreq[row][column] != 0 {
					gbfreq += int64(detect.GBKFreq[row][column])
				}
			} else if (0x81-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && // Extended GB range
				i+3 < rawtextlen && (0x30-256) <= rawtext[i+1] && rawtext[i+1] <= (0x39-256) && (0x81-256) <= rawtext[i+2] && rawtext[i+2] <= (0xFE-256) && (0x30-256) <= rawtext[i+3] && rawtext[i+3] <= (0x39-256) {
				gbchars++
				/*
				 * totalfreq += 500; row = rawtext[i] + 256 - 0x81; if (0x40
				 * <= rawtext[i+1] && rawtext[i+1] <= 0x7E) { column =
				 * rawtext[i+1] - 0x40; } else { column = rawtext[i+1] + 256
				 * - 0x40; } //System.out.println("extended row " + row + "
				 * column " + column + " rawtext[i] " + rawtext[i]); if
				 * (GBKFreq[row][column] != 0) { gbfreq +=
				 * GBKFreq[row][column]; }
				 */
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(gbchars)/float64(dbchars))
	freqval = 50 * float32(float64(gbfreq)/float64(totalfreq))
	// For regular GB files, this would give the same score, so I handicap
	// it slightly
	return (int)(rangeval+freqval) - 1
}

/*
 * Function: hz_probability Argument: byte array Returns : number from 0 to
 * 100 representing probability text in array uses HZ encoding
 */
func (detect *BytesEncodingDetect) hz_probability(rawtext []int) int {
	var i, rawtextlen int
	var hzchars = 0
	var dbchars = 1
	var hzfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var hzstart = 0
	var hzend = 0
	var row, column int
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen; i++ {
		if rawtext[i] == '~' {
			if rawtext[i+1] == '{' {
				hzstart++
				i += 2
				for i < rawtextlen-1 {
					if rawtext[i] == 0x0A || rawtext[i] == 0x0D {
						break
					} else if rawtext[i] == '~' && rawtext[i+1] == '}' {
						hzend++
						i++
						break
					} else if (0x21 <= rawtext[i] && rawtext[i] <= 0x77) && (0x21 <= rawtext[i+1] && rawtext[i+1] <= 0x77) {
						hzchars += 2
						row = rawtext[i] - 0x21
						column = rawtext[i+1] - 0x21
						totalfreq += 500
						if detect.GBFreq[row][column] != 0 {
							hzfreq += int64(detect.GBFreq[row][column])
						} else if 15 <= row && row < 55 {
							hzfreq += 200
						}
					} else if (0xA1 <= rawtext[i] && rawtext[i] <= 0xF7) && (0xA1 <= rawtext[i+1] && rawtext[i+1] <= 0xF7) {
						hzchars += 2
						row = rawtext[i] + 256 - 0xA1
						column = rawtext[i+1] + 256 - 0xA1
						totalfreq += 500
						if detect.GBFreq[row][column] != 0 {
							hzfreq += int64(detect.GBFreq[row][column])
						} else if 15 <= row && row < 55 {
							hzfreq += 200
						}
					}
					dbchars += 2
					i += 2
				}
			} else if rawtext[i+1] == '}' {
				hzend++
				i++
			} else if rawtext[i+1] == '~' {
				i++
			}
		}
	}
	if hzstart > 4 {
		rangeval = 50
	} else if hzstart > 1 {
		rangeval = 41
	} else if hzstart > 0 { // Only 39 in case the sequence happened to
		// occur
		rangeval = 39 // in otherwise non-Hz text
	} else {
		rangeval = 0
	}
	freqval = 50 * float32(float64(hzfreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

/**
 * Function: big5_probability Argument: byte array Returns : number from 0
 * to 100 representing probability text in array uses Big5 encoding
 */
func (detect *BytesEncodingDetect) big5_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var bfchars = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var bffreq int64 = 0
	var totalfreq int64 = 1
	var row, column int
	// Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		if rawtext[i] >= 0 {
			// asciichars++
		} else {
			dbchars++
			if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xF9-256) && (((0x40-256) <= rawtext[i+1] && rawtext[i+1] <= (0x7E-256)) || ((0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256))) {
				bfchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				if 0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E {
					column = rawtext[i+1] - 0x40
				} else {
					column = rawtext[i+1] + 256 - 0x61
				}
				if detect.Big5Freq[row][column] != 0 {
					bffreq += int64(detect.Big5Freq[row][column])
				} else if 3 <= row && row <= 37 {
					bffreq += 200
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(bfchars)/float64(dbchars))
	freqval = 50 * float32(float64(bffreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

/*
 * Function: big5plus_probability Argument: pointer to unsigned char array
 * Returns : number from 0 to 100 representing probability text in array
 * uses Big5+ encoding
 */
func (detect *BytesEncodingDetect) big5plus_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var bfchars = 1
	var bffreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 128 {
			// asciichars++;
		} else {
			dbchars++
			if 0xA1 <= rawtext[i] && rawtext[i] <= 0xF9 && // Original Big5 range
				((0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E) || (0xA1 <= rawtext[i+1] && rawtext[i+1] <= 0xFE)) {
				bfchars++
				totalfreq += 500
				row = rawtext[i] - 0xA1
				if 0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E {
					column = rawtext[i+1] - 0x40
				} else {
					column = rawtext[i+1] - 0x61
				}
				// System.out.println("original row " + row + " column " +
				// column);
				if detect.Big5Freq[row][column] != 0 {
					bffreq += int64(detect.Big5Freq[row][column])
				} else if 3 <= row && row < 37 {
					bffreq += 200
				}
			} else if 0x81 <= rawtext[i] && rawtext[i] <= 0xFE && // Extended Big5 range
				((0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E) || (0x80 <= rawtext[i+1] && rawtext[i+1] <= 0xFE)) {
				bfchars++
				totalfreq += 500
				row = rawtext[i] - 0x81
				if 0x40 <= rawtext[i+1] && rawtext[i+1] <= 0x7E {
					column = rawtext[i+1] - 0x40
				} else {
					column = rawtext[i+1] - 0x40
				}
				// System.out.println("extended row " + row + " column " +
				// column + " rawtext[i] " + rawtext[i]);
				if detect.Big5PFreq[row][column] != 0 {
					bffreq += int64(detect.Big5PFreq[row][column])
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(bfchars)/float64(dbchars))
	freqval = 50 * float32(float64(bffreq)/float64(totalfreq))
	// For regular Big5 files, this would give the same score, so I handicap
	// it slightly
	return (int)(rangeval+freqval) - 1
}

/*
 * Function: euc_tw_probability Argument: byte array Returns : number from 0
 * to 100 representing probability text in array uses EUC-TW (CNS 11643)
 * encoding
 */
func (detect *BytesEncodingDetect) euc_tw_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var cnschars = 1
	var cnsfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Check to see if characters fit into acceptable ranges
	// and have expected frequency of use
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		if rawtext[i] >= 0 { // in ASCII range
			// asciichars++;
		} else { // high bit set
			dbchars++
			if i+3 < rawtextlen && (0x8E-256) == rawtext[i] && (0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xB0-256) && (0xA1-256) <= rawtext[i+2] && rawtext[i+2] <= (0xFE-256) && (0xA1-256) <= rawtext[i+3] && rawtext[i+3] <= (0xFE-256) { // Planes 1 - 16
				cnschars++
				// System.out.println("plane 2 or above CNS char");
				// These are all less frequent chars so just ignore freq
				i += 3
			} else if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && // Plane 1
				(0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
				cnschars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				column = rawtext[i+1] + 256 - 0xA1
				if detect.EucTwfreq[row][column] != 0 {
					cnsfreq += int64(detect.EucTwfreq[row][column])
				} else if 35 <= row && row <= 92 {
					cnsfreq += 150
				}
				i++
			}
		}
	}
	rangeval = 50 * float32(float64(cnschars)/float64(dbchars))
	freqval = 50 * float32(float64(cnsfreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

/*
 * Function: iso_2022_cn_probability Argument: byte array Returns : number
 * from 0 to 100 representing probability text in array uses ISO 2022-CN
 * encoding WORKS FOR BASIC CASES, BUT STILL NEEDS MORE WORK
 */
func (detect *BytesEncodingDetect) iso_2022_cn_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var isochars = 1
	var isofreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Check to see if characters fit into acceptable ranges
	// and have expected frequency of use
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		if rawtext[i] == (0x1B-256) && i+3 < rawtextlen { // Escape
			// char ESC
			if rawtext[i+1] == (0x24-256) && rawtext[i+2] == 0x29 && rawtext[i+3] == (0x41-256) { // GB Escape $ ) A
				i += 4
				for rawtext[i] != (0x1B - 256) {
					dbchars++
					if (0x21 <= rawtext[i] && rawtext[i] <= 0x77) && (0x21 <= rawtext[i+1] && rawtext[i+1] <= 0x77) {
						isochars++
						row = rawtext[i] - 0x21
						column = rawtext[i+1] - 0x21
						totalfreq += 500
						if detect.GBFreq[row][column] != 0 {
							isofreq += int64(detect.GBFreq[row][column])
						} else if 15 <= row && row < 55 {
							isofreq += 200
						}
						i++
					}
					i++
				}
			} else if i+3 < rawtextlen && rawtext[i+1] == (0x24-256) && rawtext[i+2] == (0x29-256) && rawtext[i+3] == (0x47-256) {
				// CNS Escape $ ) G
				i += 4
				for rawtext[i] != (0x1B - 256) {
					dbchars++
					if (0x21-256) <= rawtext[i] && rawtext[i] <= (0x7E-256) && (0x21-256) <= rawtext[i+1] && rawtext[i+1] <= (0x7E-256) {
						isochars++
						totalfreq += 500
						row = rawtext[i] - 0x21
						column = rawtext[i+1] - 0x21
						if detect.EucTwfreq[row][column] != 0 {
							isofreq += int64(detect.EucTwfreq[row][column])
						} else if 35 <= row && row <= 92 {
							isofreq += 150
						}
						i++
					}
					i++
				}
			}
			if rawtext[i] == (0x1B-256) && i+2 < rawtextlen && rawtext[i+1] == (0x28-256) && rawtext[i+2] == (0x42-256) { // ASCII:
				// ESC
				// ( B
				i += 2
			}
		}
	}
	rangeval = 50 * float32(float64(isochars)/float64(dbchars))
	freqval = 50 * float32(float64(isofreq)/float64(totalfreq))
	// System.out.println("isochars dbchars isofreq totalfreq " + isochars +
	// " " + dbchars + " " + isofreq + " " + totalfreq + "
	// " + rangeval + " " + freqval);
	return (int)(rangeval + freqval)
	// return 0;
}

/*
 * Function: utf8_probability Argument: byte array Returns : number from 0
 * to 100 representing probability text in array uses UTF-8 encoding of
 * Unicode
 */
func (detect *BytesEncodingDetect) utf8_probability(rawtext []int) int {
	var score = 0
	var i, rawtextlen int
	var goodbytes = 0
	var asciibytes = 0
	// Maybe also use UTF8 Byte Order Mark: EF BB BF
	// Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen; i++ {
		if (rawtext[i] & (0x7F - 256)) == rawtext[i] { // One byte
			asciibytes++
			// Ignore ASCII, can throw off count
		} else if -64 <= rawtext[i] && rawtext[i] <= -33 && // Two bytes
			i+1 < rawtextlen && -128 <= rawtext[i+1] && rawtext[i+1] <= -65 {
			goodbytes += 2
			i++
		} else if -32 <= rawtext[i] && rawtext[i] <= -17 && // Three bytes
			i+2 < rawtextlen && -128 <= rawtext[i+1] && rawtext[i+1] <= -65 && -128 <= rawtext[i+2] && rawtext[i+2] <= -65 {
			goodbytes += 3
			i += 2
		}
	}
	if asciibytes == rawtextlen {
		return 0
	}
	score = (int)(100 * (float64(goodbytes) / float64(rawtextlen-asciibytes)))
	// System.out.println("rawtextlen " + rawtextlen + " goodbytes " +
	// goodbytes + " asciibytes " + asciibytes + " score " +
	// score);
	// If not above 98, reduce to zero to prevent coincidental matches
	// Allows for some (few) bad formed sequences
	if score > 98 {
		return score
	} else if score > 95 && goodbytes > 30 {
		return score
	} else {
		return 0
	}
}

/*
 * Function: utf16_probability Argument: byte array Returns : number from 0
 * to 100 representing probability text in array uses UTF-16 encoding of
 * Unicode, guess based on BOM // NOT VERY GENERAL, NEEDS MUCH MORE WORK
 */
func (detect *BytesEncodingDetect) utf16_probability(rawtext []int) int {
	// int score = 0;
	// int i, rawtextlen = 0;
	// int goodbytes = 0, asciibytes = 0;
	if len(rawtext) > 1 && ((0xFE-256) == rawtext[0] && (0xFF-256) == rawtext[1]) || // Big-endian
		((0xFF-256) == rawtext[0] && (0xFE-256) == rawtext[1]) { // Little-endian
		return 100
	}
	return 0
	/*
	 * // Check to see if characters fit into acceptable ranges rawtextlen =
	 * len(rawtext) for (i = 0; i < rawtextlen; i++) { if ((rawtext[i] &
	 * (byte)0x7F) == rawtext[i]) { // One byte goodbytes += 1;
	 * asciibytes++; } else if ((rawtext[i] & (byte)0xDF) == rawtext[i]) {
	 * // Two bytes if (i+1 < rawtextlen && (rawtext[i+1] & (byte)0xBF) ==
	 * rawtext[i+1]) { goodbytes += 2; i++; } } else if ((rawtext[i] &
	 * (byte)0xEF) == rawtext[i]) { // Three bytes if (i+2 < rawtextlen &&
	 * (rawtext[i+1] & (byte)0xBF) == rawtext[i+1] && (rawtext[i+2] &
	 * (byte)0xBF) == rawtext[i+2]) { goodbytes += 3; i+=2; } } }
	 *
	 * score = (int)(100 * ((float)goodbytes/(float)rawtext.length)); // An
	 * all ASCII file is also a good UTF8 file, but I'd rather it // get
	 * identified as ASCII. Can delete following 3 lines otherwise if
	 * (goodbytes == asciibytes) { score = 0; } // If not above 90, reduce
	 * to zero to prevent coincidental matches if (score > 90) { return
	 * score; } else { return 0; }
	 */
}

/*
 * Function: ascii_probability Argument: byte array Returns : number from 0
 * to 100 representing probability text in array uses all ASCII Description:
 * Sees if array has any characters not in ASCII range, if so, score is
 * reduced
 */
func (detect *BytesEncodingDetect) ascii_probability(rawtext []int) int {
	var score = 75
	var i, rawtextlen int
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen; i++ {
		if rawtext[i] < 0 {
			score = score - 5
		} else if rawtext[i] == (0x1B - 256) { // ESC (used by ISO 2022)
			score = score - 5
		}
		if score <= 0 {
			return 0
		}
	}
	return score
}

/*
 * Function: euc_kr__probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses EUC-KR
 * encoding
 */
func (detect *BytesEncodingDetect) euc_kr_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var krchars = 1
	var krfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && (0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
				krchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				column = rawtext[i+1] + 256 - 0xA1
				if detect.KRFreq[row][column] != 0 {
					krfreq += int64(detect.KRFreq[row][column])
				} else if 15 <= row && row < 55 {
					krfreq += 0
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(krchars)/float64(dbchars))
	freqval = 50 * float32(float64(krfreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

/*
 * Function: cp949__probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses Cp949
 * encoding
 */
func (detect *BytesEncodingDetect) cp949_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var krchars = 1
	var krfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if (0x81-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && ((0x41-256) <= rawtext[i+1] && rawtext[i+1] <= (0x5A-256) || (0x61-256) <= rawtext[i+1] && rawtext[i+1] <= (0x7A-256) || (0x81-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256)) {
				krchars++
				totalfreq += 500
				if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && (0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
					row = rawtext[i] + 256 - 0xA1
					column = rawtext[i+1] + 256 - 0xA1
					if detect.KRFreq[row][column] != 0 {
						krfreq += int64(detect.KRFreq[row][column])
					}
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(krchars)/float64(dbchars))
	freqval = 50 * float32(float64(krfreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

func (detect *BytesEncodingDetect) iso_2022_kr_probability(rawtext []int) int {
	var i int
	for i = 0; i < len(rawtext); i++ {
		if i+3 < len(rawtext) && rawtext[i] == 0x1b && rawtext[i+1] == '$' && rawtext[i+2] == ')' && rawtext[i+3] == 'C' {
			return 100
		}
	}
	return 0
}

/*
 * Function: euc_jp_probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses EUC-JP
 * encoding
 */
func (detect *BytesEncodingDetect) euc_jp_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var jpchars = 1
	var jpfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xFE-256) && (0xA1-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFE-256) {
				jpchars++
				totalfreq += 500
				row = rawtext[i] + 256 - 0xA1
				column = rawtext[i+1] + 256 - 0xA1
				if detect.JPFreq[row][column] != 0 {
					jpfreq += int64(detect.JPFreq[row][column])
				} else if 15 <= row && row < 55 {
					jpfreq += 0
				}
			}
			i++
		}
	}
	rangeval = 50 * float32(float64(jpchars)/float64(dbchars))
	freqval = 50 * float32(float64(jpfreq)/float64(totalfreq))
	return (int)(rangeval + freqval)
}

func (detect *BytesEncodingDetect) iso_2022_jp_probability(rawtext []int) int {
	var i int
	for i = 0; i < len(rawtext); i++ {
		if i+2 < len(rawtext) && rawtext[i] == 0x1b && rawtext[i+1] == '$' && rawtext[i+2] == 'B' {
			return 100
		}
	}
	return 0
}

/*
 * Function: sjis_probability Argument: pointer to byte array Returns :
 * number from 0 to 100 representing probability text in array uses
 * Shift-JIS encoding
 */
func (detect *BytesEncodingDetect) sjis_probability(rawtext []int) int {
	var i, rawtextlen int
	var dbchars = 1
	var jpchars = 1
	var jpfreq int64 = 0
	var totalfreq int64 = 1
	var rangeval float32 = 0
	var freqval float32 = 0
	var row, column, adjust int
	// Stage 1: Check to see if characters fit into acceptable ranges
	rawtextlen = len(rawtext)
	for i = 0; i < rawtextlen-1; i++ {
		// System.err.println(rawtext[i]);
		if rawtext[i] >= 0 {
			// asciichars++;
		} else {
			dbchars++
			if i+1 < len(rawtext) && (((0x81-256) <= rawtext[i] && rawtext[i] <= (0x9F-256)) || ((0xE0-256) <= rawtext[i] && rawtext[i] <= (0xEF-256))) && (((0x40-256) <= rawtext[i+1] && rawtext[i+1] <= (0x7E-256)) || ((0x80-256) <= rawtext[i+1] && rawtext[i+1] <= (0xFC-256))) {
				jpchars++
				totalfreq += 500
				row = rawtext[i] + 256
				column = rawtext[i+1] + 256
				if column < 0x9f {
					adjust = 1
					if column > 0x7f {
						column -= 0x20
					} else {
						column -= 0x19
					}
				} else {
					adjust = 0
					column -= 0x7e
				}
				if row < 0xa0 {
					row = ((row - 0x70) << 1) - adjust
				} else {
					row = ((row - 0xb0) << 1) - adjust
				}
				row -= 0x20
				column = 0x20
				// System.out.println("original row " + row + " column " +
				// column);
				if row < len(detect.JPFreq) && column < len(detect.JPFreq[row]) && detect.JPFreq[row][column] != 0 {
					jpfreq += int64(detect.JPFreq[row][column])
				}
				i++
			} else if (0xA1-256) <= rawtext[i] && rawtext[i] <= (0xDF-256) {
				// half-width katakana, convert to full-width
			}
		}
	}
	rangeval = 50 * float32(float64(jpchars)/float64(dbchars))
	freqval = 50 * float32(float64(jpfreq)/float64(totalfreq))
	// For regular GB files, this would give the same score, so I handicap
	// it slightly
	return (int)(rangeval+freqval) - 1
}

func (detect *BytesEncodingDetect) initializeFrequencies() {
	for i := 93; i >= 0; i-- {
		for j := 93; j >= 0; j-- {
			detect.GBFreq[i][j] = 0
		}
	}
	for i := 125; i >= 0; i-- {
		for j := 190; j >= 0; j-- {
			detect.GBKFreq[i][j] = 0
		}
	}
	// for i := 0; i < 94; i++ {
	// for j := 0; j < 158; j++ {
	for i := 93; i >= 0; i-- {
		for j := 157; j >= 0; j-- {
			detect.Big5Freq[i][j] = 0
		}
	}
	// for i := 0; i < 126; i++ {
	// for j := 0; j < 191; j++ {
	for i := 125; i >= 0; i-- {
		for j := 190; j >= 0; j-- {
			detect.Big5PFreq[i][j] = 0
		}
	}
	// for i := 0; i < 94; i++ {
	// for j := 0; j < 94; j++ {
	for i := 93; i >= 0; i-- {
		for j := 93; j >= 0; j-- {
			detect.EucTwfreq[i][j] = 0
		}
	}
	for i := 93; i >= 0; i-- {
		for j := 93; j >= 0; j-- {
			detect.JPFreq[i][j] = 0
		}
	}
	detect.GBFreq[20][35] = 599
	detect.GBFreq[49][26] = 598
	detect.GBFreq[41][38] = 597
	detect.GBFreq[17][26] = 596
	detect.GBFreq[32][42] = 595
	detect.GBFreq[39][42] = 594
	detect.GBFreq[45][49] = 593
	detect.GBFreq[51][57] = 592
	detect.GBFreq[50][47] = 591
	detect.GBFreq[42][90] = 590
	detect.GBFreq[52][65] = 589
	detect.GBFreq[53][47] = 588
	detect.GBFreq[19][82] = 587
	detect.GBFreq[31][19] = 586
	detect.GBFreq[40][46] = 585
	detect.GBFreq[24][89] = 584
	detect.GBFreq[23][85] = 583
	detect.GBFreq[20][28] = 582
	detect.GBFreq[42][20] = 581
	detect.GBFreq[34][38] = 580
	detect.GBFreq[45][9] = 579
	detect.GBFreq[54][50] = 578
	detect.GBFreq[25][44] = 577
	detect.GBFreq[35][66] = 576
	detect.GBFreq[20][55] = 575
	detect.GBFreq[18][85] = 574
	detect.GBFreq[20][31] = 573
	detect.GBFreq[49][17] = 572
	detect.GBFreq[41][16] = 571
	detect.GBFreq[35][73] = 570
	detect.GBFreq[20][34] = 569
	detect.GBFreq[29][44] = 568
	detect.GBFreq[35][38] = 567
	detect.GBFreq[49][9] = 566
	detect.GBFreq[46][33] = 565
	detect.GBFreq[49][51] = 564
	detect.GBFreq[40][89] = 563
	detect.GBFreq[26][64] = 562
	detect.GBFreq[54][51] = 561
	detect.GBFreq[54][36] = 560
	detect.GBFreq[39][4] = 559
	detect.GBFreq[53][13] = 558
	detect.GBFreq[24][92] = 557
	detect.GBFreq[27][49] = 556
	detect.GBFreq[48][6] = 555
	detect.GBFreq[21][51] = 554
	detect.GBFreq[30][40] = 553
	detect.GBFreq[42][92] = 552
	detect.GBFreq[31][78] = 551
	detect.GBFreq[25][82] = 550
	detect.GBFreq[47][0] = 549
	detect.GBFreq[34][19] = 548
	detect.GBFreq[47][35] = 547
	detect.GBFreq[21][63] = 546
	detect.GBFreq[43][75] = 545
	detect.GBFreq[21][87] = 544
	detect.GBFreq[35][59] = 543
	detect.GBFreq[25][34] = 542
	detect.GBFreq[21][27] = 541
	detect.GBFreq[39][26] = 540
	detect.GBFreq[34][26] = 539
	detect.GBFreq[39][52] = 538
	detect.GBFreq[50][57] = 537
	detect.GBFreq[37][79] = 536
	detect.GBFreq[26][24] = 535
	detect.GBFreq[22][1] = 534
	detect.GBFreq[18][40] = 533
	detect.GBFreq[41][33] = 532
	detect.GBFreq[53][26] = 531
	detect.GBFreq[54][86] = 530
	detect.GBFreq[20][16] = 529
	detect.GBFreq[46][74] = 528
	detect.GBFreq[30][19] = 527
	detect.GBFreq[45][35] = 526
	detect.GBFreq[45][61] = 525
	detect.GBFreq[30][9] = 524
	detect.GBFreq[41][53] = 523
	detect.GBFreq[41][13] = 522
	detect.GBFreq[50][34] = 521
	detect.GBFreq[53][86] = 520
	detect.GBFreq[47][47] = 519
	detect.GBFreq[22][28] = 518
	detect.GBFreq[50][53] = 517
	detect.GBFreq[39][70] = 516
	detect.GBFreq[38][15] = 515
	detect.GBFreq[42][88] = 514
	detect.GBFreq[16][29] = 513
	detect.GBFreq[27][90] = 512
	detect.GBFreq[29][12] = 511
	detect.GBFreq[44][22] = 510
	detect.GBFreq[34][69] = 509
	detect.GBFreq[24][10] = 508
	detect.GBFreq[44][11] = 507
	detect.GBFreq[39][92] = 506
	detect.GBFreq[49][48] = 505
	detect.GBFreq[31][46] = 504
	detect.GBFreq[19][50] = 503
	detect.GBFreq[21][14] = 502
	detect.GBFreq[32][28] = 501
	detect.GBFreq[18][3] = 500
	detect.GBFreq[53][9] = 499
	detect.GBFreq[34][80] = 498
	detect.GBFreq[48][88] = 497
	detect.GBFreq[46][53] = 496
	detect.GBFreq[22][53] = 495
	detect.GBFreq[28][10] = 494
	detect.GBFreq[44][65] = 493
	detect.GBFreq[20][10] = 492
	detect.GBFreq[40][76] = 491
	detect.GBFreq[47][8] = 490
	detect.GBFreq[50][74] = 489
	detect.GBFreq[23][62] = 488
	detect.GBFreq[49][65] = 487
	detect.GBFreq[28][87] = 486
	detect.GBFreq[15][48] = 485
	detect.GBFreq[22][7] = 484
	detect.GBFreq[19][42] = 483
	detect.GBFreq[41][20] = 482
	detect.GBFreq[26][55] = 481
	detect.GBFreq[21][93] = 480
	detect.GBFreq[31][76] = 479
	detect.GBFreq[34][31] = 478
	detect.GBFreq[20][66] = 477
	detect.GBFreq[51][33] = 476
	detect.GBFreq[34][86] = 475
	detect.GBFreq[37][67] = 474
	detect.GBFreq[53][53] = 473
	detect.GBFreq[40][88] = 472
	detect.GBFreq[39][10] = 471
	detect.GBFreq[24][3] = 470
	detect.GBFreq[27][25] = 469
	detect.GBFreq[26][15] = 468
	detect.GBFreq[21][88] = 467
	detect.GBFreq[52][62] = 466
	detect.GBFreq[46][81] = 465
	detect.GBFreq[38][72] = 464
	detect.GBFreq[17][30] = 463
	detect.GBFreq[52][92] = 462
	detect.GBFreq[34][90] = 461
	detect.GBFreq[21][7] = 460
	detect.GBFreq[36][13] = 459
	detect.GBFreq[45][41] = 458
	detect.GBFreq[32][5] = 457
	detect.GBFreq[26][89] = 456
	detect.GBFreq[23][87] = 455
	detect.GBFreq[20][39] = 454
	detect.GBFreq[27][23] = 453
	detect.GBFreq[25][59] = 452
	detect.GBFreq[49][20] = 451
	detect.GBFreq[54][77] = 450
	detect.GBFreq[27][67] = 449
	detect.GBFreq[47][33] = 448
	detect.GBFreq[41][17] = 447
	detect.GBFreq[19][81] = 446
	detect.GBFreq[16][66] = 445
	detect.GBFreq[45][26] = 444
	detect.GBFreq[49][81] = 443
	detect.GBFreq[53][55] = 442
	detect.GBFreq[16][26] = 441
	detect.GBFreq[54][62] = 440
	detect.GBFreq[20][70] = 439
	detect.GBFreq[42][35] = 438
	detect.GBFreq[20][57] = 437
	detect.GBFreq[34][36] = 436
	detect.GBFreq[46][63] = 435
	detect.GBFreq[19][45] = 434
	detect.GBFreq[21][10] = 433
	detect.GBFreq[52][93] = 432
	detect.GBFreq[25][2] = 431
	detect.GBFreq[30][57] = 430
	detect.GBFreq[41][24] = 429
	detect.GBFreq[28][43] = 428
	detect.GBFreq[45][86] = 427
	detect.GBFreq[51][56] = 426
	detect.GBFreq[37][28] = 425
	detect.GBFreq[52][69] = 424
	detect.GBFreq[43][92] = 423
	detect.GBFreq[41][31] = 422
	detect.GBFreq[37][87] = 421
	detect.GBFreq[47][36] = 420
	detect.GBFreq[16][16] = 419
	detect.GBFreq[40][56] = 418
	detect.GBFreq[24][55] = 417
	detect.GBFreq[17][1] = 416
	detect.GBFreq[35][57] = 415
	detect.GBFreq[27][50] = 414
	detect.GBFreq[26][14] = 413
	detect.GBFreq[50][40] = 412
	detect.GBFreq[39][19] = 411
	detect.GBFreq[19][89] = 410
	detect.GBFreq[29][91] = 409
	detect.GBFreq[17][89] = 408
	detect.GBFreq[39][74] = 407
	detect.GBFreq[46][39] = 406
	detect.GBFreq[40][28] = 405
	detect.GBFreq[45][68] = 404
	detect.GBFreq[43][10] = 403
	detect.GBFreq[42][13] = 402
	detect.GBFreq[44][81] = 401
	detect.GBFreq[41][47] = 400
	detect.GBFreq[48][58] = 399
	detect.GBFreq[43][68] = 398
	detect.GBFreq[16][79] = 397
	detect.GBFreq[19][5] = 396
	detect.GBFreq[54][59] = 395
	detect.GBFreq[17][36] = 394
	detect.GBFreq[18][0] = 393
	detect.GBFreq[41][5] = 392
	detect.GBFreq[41][72] = 391
	detect.GBFreq[16][39] = 390
	detect.GBFreq[54][0] = 389
	detect.GBFreq[51][16] = 388
	detect.GBFreq[29][36] = 387
	detect.GBFreq[47][5] = 386
	detect.GBFreq[47][51] = 385
	detect.GBFreq[44][7] = 384
	detect.GBFreq[35][30] = 383
	detect.GBFreq[26][9] = 382
	detect.GBFreq[16][7] = 381
	detect.GBFreq[32][1] = 380
	detect.GBFreq[33][76] = 379
	detect.GBFreq[34][91] = 378
	detect.GBFreq[52][36] = 377
	detect.GBFreq[26][77] = 376
	detect.GBFreq[35][48] = 375
	detect.GBFreq[40][80] = 374
	detect.GBFreq[41][92] = 373
	detect.GBFreq[27][93] = 372
	detect.GBFreq[15][17] = 371
	detect.GBFreq[16][76] = 370
	detect.GBFreq[51][12] = 369
	detect.GBFreq[18][20] = 368
	detect.GBFreq[15][54] = 367
	detect.GBFreq[50][5] = 366
	detect.GBFreq[33][22] = 365
	detect.GBFreq[37][57] = 364
	detect.GBFreq[28][47] = 363
	detect.GBFreq[42][31] = 362
	detect.GBFreq[18][2] = 361
	detect.GBFreq[43][64] = 360
	detect.GBFreq[23][47] = 359
	detect.GBFreq[28][79] = 358
	detect.GBFreq[25][45] = 357
	detect.GBFreq[23][91] = 356
	detect.GBFreq[22][19] = 355
	detect.GBFreq[25][46] = 354
	detect.GBFreq[22][36] = 353
	detect.GBFreq[54][85] = 352
	detect.GBFreq[46][20] = 351
	detect.GBFreq[27][37] = 350
	detect.GBFreq[26][81] = 349
	detect.GBFreq[42][29] = 348
	detect.GBFreq[31][90] = 347
	detect.GBFreq[41][59] = 346
	detect.GBFreq[24][65] = 345
	detect.GBFreq[44][84] = 344
	detect.GBFreq[24][90] = 343
	detect.GBFreq[38][54] = 342
	detect.GBFreq[28][70] = 341
	detect.GBFreq[27][15] = 340
	detect.GBFreq[28][80] = 339
	detect.GBFreq[29][8] = 338
	detect.GBFreq[45][80] = 337
	detect.GBFreq[53][37] = 336
	detect.GBFreq[28][65] = 335
	detect.GBFreq[23][86] = 334
	detect.GBFreq[39][45] = 333
	detect.GBFreq[53][32] = 332
	detect.GBFreq[38][68] = 331
	detect.GBFreq[45][78] = 330
	detect.GBFreq[43][7] = 329
	detect.GBFreq[46][82] = 328
	detect.GBFreq[27][38] = 327
	detect.GBFreq[16][62] = 326
	detect.GBFreq[24][17] = 325
	detect.GBFreq[22][70] = 324
	detect.GBFreq[52][28] = 323
	detect.GBFreq[23][40] = 322
	detect.GBFreq[28][50] = 321
	detect.GBFreq[42][91] = 320
	detect.GBFreq[47][76] = 319
	detect.GBFreq[15][42] = 318
	detect.GBFreq[43][55] = 317
	detect.GBFreq[29][84] = 316
	detect.GBFreq[44][90] = 315
	detect.GBFreq[53][16] = 314
	detect.GBFreq[22][93] = 313
	detect.GBFreq[34][10] = 312
	detect.GBFreq[32][53] = 311
	detect.GBFreq[43][65] = 310
	detect.GBFreq[28][7] = 309
	detect.GBFreq[35][46] = 308
	detect.GBFreq[21][39] = 307
	detect.GBFreq[44][18] = 306
	detect.GBFreq[40][10] = 305
	detect.GBFreq[54][53] = 304
	detect.GBFreq[38][74] = 303
	detect.GBFreq[28][26] = 302
	detect.GBFreq[15][13] = 301
	detect.GBFreq[39][34] = 300
	detect.GBFreq[39][46] = 299
	detect.GBFreq[42][66] = 298
	detect.GBFreq[33][58] = 297
	detect.GBFreq[15][56] = 296
	detect.GBFreq[18][51] = 295
	detect.GBFreq[49][68] = 294
	detect.GBFreq[30][37] = 293
	detect.GBFreq[51][84] = 292
	detect.GBFreq[51][9] = 291
	detect.GBFreq[40][70] = 290
	detect.GBFreq[41][84] = 289
	detect.GBFreq[28][64] = 288
	detect.GBFreq[32][88] = 287
	detect.GBFreq[24][5] = 286
	detect.GBFreq[53][23] = 285
	detect.GBFreq[42][27] = 284
	detect.GBFreq[22][38] = 283
	detect.GBFreq[32][86] = 282
	detect.GBFreq[34][30] = 281
	detect.GBFreq[38][63] = 280
	detect.GBFreq[24][59] = 279
	detect.GBFreq[22][81] = 278
	detect.GBFreq[32][11] = 277
	detect.GBFreq[51][21] = 276
	detect.GBFreq[54][41] = 275
	detect.GBFreq[21][50] = 274
	detect.GBFreq[23][89] = 273
	detect.GBFreq[19][87] = 272
	detect.GBFreq[26][7] = 271
	detect.GBFreq[30][75] = 270
	detect.GBFreq[43][84] = 269
	detect.GBFreq[51][25] = 268
	detect.GBFreq[16][67] = 267
	detect.GBFreq[32][9] = 266
	detect.GBFreq[48][51] = 265
	detect.GBFreq[39][7] = 264
	detect.GBFreq[44][88] = 263
	detect.GBFreq[52][24] = 262
	detect.GBFreq[23][34] = 261
	detect.GBFreq[32][75] = 260
	detect.GBFreq[19][10] = 259
	detect.GBFreq[28][91] = 258
	detect.GBFreq[32][83] = 257
	detect.GBFreq[25][75] = 256
	detect.GBFreq[53][45] = 255
	detect.GBFreq[29][85] = 254
	detect.GBFreq[53][59] = 253
	detect.GBFreq[16][2] = 252
	detect.GBFreq[19][78] = 251
	detect.GBFreq[15][75] = 250
	detect.GBFreq[51][42] = 249
	detect.GBFreq[45][67] = 248
	detect.GBFreq[15][74] = 247
	detect.GBFreq[25][81] = 246
	detect.GBFreq[37][62] = 245
	detect.GBFreq[16][55] = 244
	detect.GBFreq[18][38] = 243
	detect.GBFreq[23][23] = 242
	detect.GBFreq[38][30] = 241
	detect.GBFreq[17][28] = 240
	detect.GBFreq[44][73] = 239
	detect.GBFreq[23][78] = 238
	detect.GBFreq[40][77] = 237
	detect.GBFreq[38][87] = 236
	detect.GBFreq[27][19] = 235
	detect.GBFreq[38][82] = 234
	detect.GBFreq[37][22] = 233
	detect.GBFreq[41][30] = 232
	detect.GBFreq[54][9] = 231
	detect.GBFreq[32][30] = 230
	detect.GBFreq[30][52] = 229
	detect.GBFreq[40][84] = 228
	detect.GBFreq[53][57] = 227
	detect.GBFreq[27][27] = 226
	detect.GBFreq[38][64] = 225
	detect.GBFreq[18][43] = 224
	detect.GBFreq[23][69] = 223
	detect.GBFreq[28][12] = 222
	detect.GBFreq[50][78] = 221
	detect.GBFreq[50][1] = 220
	detect.GBFreq[26][88] = 219
	detect.GBFreq[36][40] = 218
	detect.GBFreq[33][89] = 217
	detect.GBFreq[41][28] = 216
	detect.GBFreq[31][77] = 215
	detect.GBFreq[46][1] = 214
	detect.GBFreq[47][19] = 213
	detect.GBFreq[35][55] = 212
	detect.GBFreq[41][21] = 211
	detect.GBFreq[27][10] = 210
	detect.GBFreq[32][77] = 209
	detect.GBFreq[26][37] = 208
	detect.GBFreq[20][33] = 207
	detect.GBFreq[41][52] = 206
	detect.GBFreq[32][18] = 205
	detect.GBFreq[38][13] = 204
	detect.GBFreq[20][18] = 203
	detect.GBFreq[20][24] = 202
	detect.GBFreq[45][19] = 201
	detect.GBFreq[18][53] = 200
	/*
	 * GBFreq[39][0] = 199; GBFreq[40][71] = 198; GBFreq[41][27] = 197;
	 * GBFreq[15][69] = 196; GBFreq[42][10] = 195; GBFreq[31][89] = 194;
	 * GBFreq[51][28] = 193; GBFreq[41][22] = 192; GBFreq[40][43] = 191;
	 * GBFreq[38][6] = 190; GBFreq[37][11] = 189; GBFreq[39][60] = 188;
	 * GBFreq[48][47] = 187; GBFreq[46][80] = 186; GBFreq[52][49] = 185;
	 * GBFreq[50][48] = 184; GBFreq[25][1] = 183; GBFreq[52][29] = 182;
	 * GBFreq[24][66] = 181; GBFreq[23][35] = 180; GBFreq[49][72] = 179;
	 * GBFreq[47][45] = 178; GBFreq[45][14] = 177; GBFreq[51][70] = 176;
	 * GBFreq[22][30] = 175; GBFreq[49][83] = 174; GBFreq[26][79] = 173;
	 * GBFreq[27][41] = 172; GBFreq[51][81] = 171; GBFreq[41][54] = 170;
	 * GBFreq[20][4] = 169; GBFreq[29][60] = 168; GBFreq[20][27] = 167;
	 * GBFreq[50][15] = 166; GBFreq[41][6] = 165; GBFreq[35][34] = 164;
	 * GBFreq[44][87] = 163; GBFreq[46][66] = 162; GBFreq[42][37] = 161;
	 * GBFreq[42][24] = 160; GBFreq[54][7] = 159; GBFreq[41][14] = 158;
	 * GBFreq[39][83] = 157; GBFreq[16][87] = 156; GBFreq[20][59] = 155;
	 * GBFreq[42][12] = 154; GBFreq[47][2] = 153; GBFreq[21][32] = 152;
	 * GBFreq[53][29] = 151; GBFreq[22][40] = 150; GBFreq[24][58] = 149;
	 * GBFreq[52][88] = 148; GBFreq[29][30] = 147; GBFreq[15][91] = 146;
	 * GBFreq[54][72] = 145; GBFreq[51][75] = 144; GBFreq[33][67] = 143;
	 * GBFreq[41][50] = 142; GBFreq[27][34] = 141; GBFreq[46][17] = 140;
	 * GBFreq[31][74] = 139; GBFreq[42][67] = 138; GBFreq[54][87] = 137;
	 * GBFreq[27][14] = 136; GBFreq[16][63] = 135; GBFreq[16][5] = 134;
	 * GBFreq[43][23] = 133; GBFreq[23][13] = 132; GBFreq[31][12] = 131;
	 * GBFreq[25][57] = 130; GBFreq[38][49] = 129; GBFreq[42][69] = 128;
	 * GBFreq[23][80] = 127; GBFreq[29][0] = 126; GBFreq[28][2] = 125;
	 * GBFreq[28][17] = 124; GBFreq[17][27] = 123; GBFreq[40][16] = 122;
	 * GBFreq[45][1] = 121; GBFreq[36][33] = 120; GBFreq[35][23] = 119;
	 * GBFreq[20][86] = 118; GBFreq[29][53] = 117; GBFreq[23][88] = 116;
	 * GBFreq[51][87] = 115; GBFreq[54][27] = 114; GBFreq[44][36] = 113;
	 * GBFreq[21][45] = 112; GBFreq[53][52] = 111; GBFreq[31][53] = 110;
	 * GBFreq[38][47] = 109; GBFreq[27][21] = 108; GBFreq[30][42] = 107;
	 * GBFreq[29][10] = 106; GBFreq[35][35] = 105; GBFreq[24][56] = 104;
	 * GBFreq[41][29] = 103; GBFreq[18][68] = 102; GBFreq[29][24] = 101;
	 * GBFreq[25][84] = 100; GBFreq[35][47] = 99; GBFreq[29][56] = 98;
	 * GBFreq[30][44] = 97; GBFreq[53][3] = 96; GBFreq[30][63] = 95;
	 * GBFreq[52][52] = 94; GBFreq[54][1] = 93; GBFreq[22][48] = 92;
	 * GBFreq[54][66] = 91; GBFreq[21][90] = 90; GBFreq[52][47] = 89;
	 * GBFreq[39][25] = 88; GBFreq[39][39] = 87; GBFreq[44][37] = 86;
	 * GBFreq[44][76] = 85; GBFreq[46][75] = 84; GBFreq[18][37] = 83;
	 * GBFreq[47][42] = 82; GBFreq[19][92] = 81; GBFreq[51][27] = 80;
	 * GBFreq[48][83] = 79; GBFreq[23][70] = 78; GBFreq[29][9] = 77;
	 * GBFreq[33][79] = 76; GBFreq[52][90] = 75; GBFreq[53][6] = 74;
	 * GBFreq[24][36] = 73; GBFreq[25][25] = 72; GBFreq[44][26] = 71;
	 * GBFreq[25][36] = 70; GBFreq[29][87] = 69; GBFreq[48][0] = 68;
	 * GBFreq[15][40] = 67; GBFreq[17][45] = 66; GBFreq[30][14] = 65;
	 * GBFreq[48][38] = 64; GBFreq[23][19] = 63; GBFreq[40][42] = 62;
	 * GBFreq[31][63] = 61; GBFreq[16][23] = 60; GBFreq[26][21] = 59;
	 * GBFreq[32][76] = 58; GBFreq[23][58] = 57; GBFreq[41][37] = 56;
	 * GBFreq[30][43] = 55; GBFreq[47][38] = 54; GBFreq[21][46] = 53;
	 * GBFreq[18][33] = 52; GBFreq[52][37] = 51; GBFreq[36][8] = 50;
	 * GBFreq[49][24] = 49; GBFreq[15][66] = 48; GBFreq[35][77] = 47;
	 * GBFreq[27][58] = 46; GBFreq[35][51] = 45; GBFreq[24][69] = 44;
	 * GBFreq[20][54] = 43; GBFreq[24][41] = 42; GBFreq[41][0] = 41;
	 * GBFreq[33][71] = 40; GBFreq[23][52] = 39; GBFreq[29][67] = 38;
	 * GBFreq[46][51] = 37; GBFreq[46][90] = 36; GBFreq[49][33] = 35;
	 * GBFreq[33][28] = 34; GBFreq[37][86] = 33; GBFreq[39][22] = 32;
	 * GBFreq[37][37] = 31; GBFreq[29][62] = 30; GBFreq[29][50] = 29;
	 * GBFreq[36][89] = 28; GBFreq[42][44] = 27; GBFreq[51][82] = 26;
	 * GBFreq[28][83] = 25; GBFreq[15][78] = 24; GBFreq[46][62] = 23;
	 * GBFreq[19][69] = 22; GBFreq[51][23] = 21; GBFreq[37][69] = 20;
	 * GBFreq[25][5] = 19; GBFreq[51][85] = 18; GBFreq[48][77] = 17;
	 * GBFreq[32][46] = 16; GBFreq[53][60] = 15; GBFreq[28][57] = 14;
	 * GBFreq[54][82] = 13; GBFreq[54][15] = 12; GBFreq[49][54] = 11;
	 * GBFreq[53][87] = 10; GBFreq[27][16] = 9; GBFreq[29][34] = 8;
	 * GBFreq[20][44] = 7; GBFreq[42][73] = 6; GBFreq[47][71] = 5;
	 * GBFreq[29][37] = 4; GBFreq[25][50] = 3; GBFreq[18][84] = 2;
	 * GBFreq[50][45] = 1; GBFreq[48][46] = 0;
	 */
	// GBFreq[43][89] = -1; GBFreq[54][68] = -2;
	detect.Big5Freq[9][89] = 600
	detect.Big5Freq[11][15] = 599
	detect.Big5Freq[3][66] = 598
	detect.Big5Freq[6][121] = 597
	detect.Big5Freq[3][0] = 596
	detect.Big5Freq[5][82] = 595
	detect.Big5Freq[3][42] = 594
	detect.Big5Freq[5][34] = 593
	detect.Big5Freq[3][8] = 592
	detect.Big5Freq[3][6] = 591
	detect.Big5Freq[3][67] = 590
	detect.Big5Freq[7][139] = 589
	detect.Big5Freq[23][137] = 588
	detect.Big5Freq[12][46] = 587
	detect.Big5Freq[4][8] = 586
	detect.Big5Freq[4][41] = 585
	detect.Big5Freq[18][47] = 584
	detect.Big5Freq[12][114] = 583
	detect.Big5Freq[6][1] = 582
	detect.Big5Freq[22][60] = 581
	detect.Big5Freq[5][46] = 580
	detect.Big5Freq[11][79] = 579
	detect.Big5Freq[3][23] = 578
	detect.Big5Freq[7][114] = 577
	detect.Big5Freq[29][102] = 576
	detect.Big5Freq[19][14] = 575
	detect.Big5Freq[4][133] = 574
	detect.Big5Freq[3][29] = 573
	detect.Big5Freq[4][109] = 572
	detect.Big5Freq[14][127] = 571
	detect.Big5Freq[5][48] = 570
	detect.Big5Freq[13][104] = 569
	detect.Big5Freq[3][132] = 568
	detect.Big5Freq[26][64] = 567
	detect.Big5Freq[7][19] = 566
	detect.Big5Freq[4][12] = 565
	detect.Big5Freq[11][124] = 564
	detect.Big5Freq[7][89] = 563
	detect.Big5Freq[15][124] = 562
	detect.Big5Freq[4][108] = 561
	detect.Big5Freq[19][66] = 560
	detect.Big5Freq[3][21] = 559
	detect.Big5Freq[24][12] = 558
	detect.Big5Freq[28][111] = 557
	detect.Big5Freq[12][107] = 556
	detect.Big5Freq[3][112] = 555
	detect.Big5Freq[8][113] = 554
	detect.Big5Freq[5][40] = 553
	detect.Big5Freq[26][145] = 552
	detect.Big5Freq[3][48] = 551
	detect.Big5Freq[3][70] = 550
	detect.Big5Freq[22][17] = 549
	detect.Big5Freq[16][47] = 548
	detect.Big5Freq[3][53] = 547
	detect.Big5Freq[4][24] = 546
	detect.Big5Freq[32][120] = 545
	detect.Big5Freq[24][49] = 544
	detect.Big5Freq[24][142] = 543
	detect.Big5Freq[18][66] = 542
	detect.Big5Freq[29][150] = 541
	detect.Big5Freq[5][122] = 540
	detect.Big5Freq[5][114] = 539
	detect.Big5Freq[3][44] = 538
	detect.Big5Freq[10][128] = 537
	detect.Big5Freq[15][20] = 536
	detect.Big5Freq[13][33] = 535
	detect.Big5Freq[14][87] = 534
	detect.Big5Freq[3][126] = 533
	detect.Big5Freq[4][53] = 532
	detect.Big5Freq[4][40] = 531
	detect.Big5Freq[9][93] = 530
	detect.Big5Freq[15][137] = 529
	detect.Big5Freq[10][123] = 528
	detect.Big5Freq[4][56] = 527
	detect.Big5Freq[5][71] = 526
	detect.Big5Freq[10][8] = 525
	detect.Big5Freq[5][16] = 524
	detect.Big5Freq[5][146] = 523
	detect.Big5Freq[18][88] = 522
	detect.Big5Freq[24][4] = 521
	detect.Big5Freq[20][47] = 520
	detect.Big5Freq[5][33] = 519
	detect.Big5Freq[9][43] = 518
	detect.Big5Freq[20][12] = 517
	detect.Big5Freq[20][13] = 516
	detect.Big5Freq[5][156] = 515
	detect.Big5Freq[22][140] = 514
	detect.Big5Freq[8][146] = 513
	detect.Big5Freq[21][123] = 512
	detect.Big5Freq[4][90] = 511
	detect.Big5Freq[5][62] = 510
	detect.Big5Freq[17][59] = 509
	detect.Big5Freq[10][37] = 508
	detect.Big5Freq[18][107] = 507
	detect.Big5Freq[14][53] = 506
	detect.Big5Freq[22][51] = 505
	detect.Big5Freq[8][13] = 504
	detect.Big5Freq[5][29] = 503
	detect.Big5Freq[9][7] = 502
	detect.Big5Freq[22][14] = 501
	detect.Big5Freq[8][55] = 500
	detect.Big5Freq[33][9] = 499
	detect.Big5Freq[16][64] = 498
	detect.Big5Freq[7][131] = 497
	detect.Big5Freq[34][4] = 496
	detect.Big5Freq[7][101] = 495
	detect.Big5Freq[11][139] = 494
	detect.Big5Freq[3][135] = 493
	detect.Big5Freq[7][102] = 492
	detect.Big5Freq[17][13] = 491
	detect.Big5Freq[3][20] = 490
	detect.Big5Freq[27][106] = 489
	detect.Big5Freq[5][88] = 488
	detect.Big5Freq[6][33] = 487
	detect.Big5Freq[5][139] = 486
	detect.Big5Freq[6][0] = 485
	detect.Big5Freq[17][58] = 484
	detect.Big5Freq[5][133] = 483
	detect.Big5Freq[9][107] = 482
	detect.Big5Freq[23][39] = 481
	detect.Big5Freq[5][23] = 480
	detect.Big5Freq[3][79] = 479
	detect.Big5Freq[32][97] = 478
	detect.Big5Freq[3][136] = 477
	detect.Big5Freq[4][94] = 476
	detect.Big5Freq[21][61] = 475
	detect.Big5Freq[23][123] = 474
	detect.Big5Freq[26][16] = 473
	detect.Big5Freq[24][137] = 472
	detect.Big5Freq[22][18] = 471
	detect.Big5Freq[5][1] = 470
	detect.Big5Freq[20][119] = 469
	detect.Big5Freq[3][7] = 468
	detect.Big5Freq[10][79] = 467
	detect.Big5Freq[15][105] = 466
	detect.Big5Freq[3][144] = 465
	detect.Big5Freq[12][80] = 464
	detect.Big5Freq[15][73] = 463
	detect.Big5Freq[3][19] = 462
	detect.Big5Freq[8][109] = 461
	detect.Big5Freq[3][15] = 460
	detect.Big5Freq[31][82] = 459
	detect.Big5Freq[3][43] = 458
	detect.Big5Freq[25][119] = 457
	detect.Big5Freq[16][111] = 456
	detect.Big5Freq[7][77] = 455
	detect.Big5Freq[3][95] = 454
	detect.Big5Freq[24][82] = 453
	detect.Big5Freq[7][52] = 452
	detect.Big5Freq[9][151] = 451
	detect.Big5Freq[3][129] = 450
	detect.Big5Freq[5][87] = 449
	detect.Big5Freq[3][55] = 448
	detect.Big5Freq[8][153] = 447
	detect.Big5Freq[4][83] = 446
	detect.Big5Freq[3][114] = 445
	detect.Big5Freq[23][147] = 444
	detect.Big5Freq[15][31] = 443
	detect.Big5Freq[3][54] = 442
	detect.Big5Freq[11][122] = 441
	detect.Big5Freq[4][4] = 440
	detect.Big5Freq[34][149] = 439
	detect.Big5Freq[3][17] = 438
	detect.Big5Freq[21][64] = 437
	detect.Big5Freq[26][144] = 436
	detect.Big5Freq[4][62] = 435
	detect.Big5Freq[8][15] = 434
	detect.Big5Freq[35][80] = 433
	detect.Big5Freq[7][110] = 432
	detect.Big5Freq[23][114] = 431
	detect.Big5Freq[3][108] = 430
	detect.Big5Freq[3][62] = 429
	detect.Big5Freq[21][41] = 428
	detect.Big5Freq[15][99] = 427
	detect.Big5Freq[5][47] = 426
	detect.Big5Freq[4][96] = 425
	detect.Big5Freq[20][122] = 424
	detect.Big5Freq[5][21] = 423
	detect.Big5Freq[4][157] = 422
	detect.Big5Freq[16][14] = 421
	detect.Big5Freq[3][117] = 420
	detect.Big5Freq[7][129] = 419
	detect.Big5Freq[4][27] = 418
	detect.Big5Freq[5][30] = 417
	detect.Big5Freq[22][16] = 416
	detect.Big5Freq[5][64] = 415
	detect.Big5Freq[17][99] = 414
	detect.Big5Freq[17][57] = 413
	detect.Big5Freq[8][105] = 412
	detect.Big5Freq[5][112] = 411
	detect.Big5Freq[20][59] = 410
	detect.Big5Freq[6][129] = 409
	detect.Big5Freq[18][17] = 408
	detect.Big5Freq[3][92] = 407
	detect.Big5Freq[28][118] = 406
	detect.Big5Freq[3][109] = 405
	detect.Big5Freq[31][51] = 404
	detect.Big5Freq[13][116] = 403
	detect.Big5Freq[6][15] = 402
	detect.Big5Freq[36][136] = 401
	detect.Big5Freq[12][74] = 400
	detect.Big5Freq[20][88] = 399
	detect.Big5Freq[36][68] = 398
	detect.Big5Freq[3][147] = 397
	detect.Big5Freq[15][84] = 396
	detect.Big5Freq[16][32] = 395
	detect.Big5Freq[16][58] = 394
	detect.Big5Freq[7][66] = 393
	detect.Big5Freq[23][107] = 392
	detect.Big5Freq[9][6] = 391
	detect.Big5Freq[12][86] = 390
	detect.Big5Freq[23][112] = 389
	detect.Big5Freq[37][23] = 388
	detect.Big5Freq[3][138] = 387
	detect.Big5Freq[20][68] = 386
	detect.Big5Freq[15][116] = 385
	detect.Big5Freq[18][64] = 384
	detect.Big5Freq[12][139] = 383
	detect.Big5Freq[11][155] = 382
	detect.Big5Freq[4][156] = 381
	detect.Big5Freq[12][84] = 380
	detect.Big5Freq[18][49] = 379
	detect.Big5Freq[25][125] = 378
	detect.Big5Freq[25][147] = 377
	detect.Big5Freq[15][110] = 376
	detect.Big5Freq[19][96] = 375
	detect.Big5Freq[30][152] = 374
	detect.Big5Freq[6][31] = 373
	detect.Big5Freq[27][117] = 372
	detect.Big5Freq[3][10] = 371
	detect.Big5Freq[6][131] = 370
	detect.Big5Freq[13][112] = 369
	detect.Big5Freq[36][156] = 368
	detect.Big5Freq[4][60] = 367
	detect.Big5Freq[15][121] = 366
	detect.Big5Freq[4][112] = 365
	detect.Big5Freq[30][142] = 364
	detect.Big5Freq[23][154] = 363
	detect.Big5Freq[27][101] = 362
	detect.Big5Freq[9][140] = 361
	detect.Big5Freq[3][89] = 360
	detect.Big5Freq[18][148] = 359
	detect.Big5Freq[4][69] = 358
	detect.Big5Freq[16][49] = 357
	detect.Big5Freq[6][117] = 356
	detect.Big5Freq[36][55] = 355
	detect.Big5Freq[5][123] = 354
	detect.Big5Freq[4][126] = 353
	detect.Big5Freq[4][119] = 352
	detect.Big5Freq[9][95] = 351
	detect.Big5Freq[5][24] = 350
	detect.Big5Freq[16][133] = 349
	detect.Big5Freq[10][134] = 348
	detect.Big5Freq[26][59] = 347
	detect.Big5Freq[6][41] = 346
	detect.Big5Freq[6][146] = 345
	detect.Big5Freq[19][24] = 344
	detect.Big5Freq[5][113] = 343
	detect.Big5Freq[10][118] = 342
	detect.Big5Freq[34][151] = 341
	detect.Big5Freq[9][72] = 340
	detect.Big5Freq[31][25] = 339
	detect.Big5Freq[18][126] = 338
	detect.Big5Freq[18][28] = 337
	detect.Big5Freq[4][153] = 336
	detect.Big5Freq[3][84] = 335
	detect.Big5Freq[21][18] = 334
	detect.Big5Freq[25][129] = 333
	detect.Big5Freq[6][107] = 332
	detect.Big5Freq[12][25] = 331
	detect.Big5Freq[17][109] = 330
	detect.Big5Freq[7][76] = 329
	detect.Big5Freq[15][15] = 328
	detect.Big5Freq[4][14] = 327
	detect.Big5Freq[23][88] = 326
	detect.Big5Freq[18][2] = 325
	detect.Big5Freq[6][88] = 324
	detect.Big5Freq[16][84] = 323
	detect.Big5Freq[12][48] = 322
	detect.Big5Freq[7][68] = 321
	detect.Big5Freq[5][50] = 320
	detect.Big5Freq[13][54] = 319
	detect.Big5Freq[7][98] = 318
	detect.Big5Freq[11][6] = 317
	detect.Big5Freq[9][80] = 316
	detect.Big5Freq[16][41] = 315
	detect.Big5Freq[7][43] = 314
	detect.Big5Freq[28][117] = 313
	detect.Big5Freq[3][51] = 312
	detect.Big5Freq[7][3] = 311
	detect.Big5Freq[20][81] = 310
	detect.Big5Freq[4][2] = 309
	detect.Big5Freq[11][16] = 308
	detect.Big5Freq[10][4] = 307
	detect.Big5Freq[10][119] = 306
	detect.Big5Freq[6][142] = 305
	detect.Big5Freq[18][51] = 304
	detect.Big5Freq[8][144] = 303
	detect.Big5Freq[10][65] = 302
	detect.Big5Freq[11][64] = 301
	detect.Big5Freq[11][130] = 300
	detect.Big5Freq[9][92] = 299
	detect.Big5Freq[18][29] = 298
	detect.Big5Freq[18][78] = 297
	detect.Big5Freq[18][151] = 296
	detect.Big5Freq[33][127] = 295
	detect.Big5Freq[35][113] = 294
	detect.Big5Freq[10][155] = 293
	detect.Big5Freq[3][76] = 292
	detect.Big5Freq[36][123] = 291
	detect.Big5Freq[13][143] = 290
	detect.Big5Freq[5][135] = 289
	detect.Big5Freq[23][116] = 288
	detect.Big5Freq[6][101] = 287
	detect.Big5Freq[14][74] = 286
	detect.Big5Freq[7][153] = 285
	detect.Big5Freq[3][101] = 284
	detect.Big5Freq[9][74] = 283
	detect.Big5Freq[3][156] = 282
	detect.Big5Freq[4][147] = 281
	detect.Big5Freq[9][12] = 280
	detect.Big5Freq[18][133] = 279
	detect.Big5Freq[4][0] = 278
	detect.Big5Freq[7][155] = 277
	detect.Big5Freq[9][144] = 276
	detect.Big5Freq[23][49] = 275
	detect.Big5Freq[5][89] = 274
	detect.Big5Freq[10][11] = 273
	detect.Big5Freq[3][110] = 272
	detect.Big5Freq[3][40] = 271
	detect.Big5Freq[29][115] = 270
	detect.Big5Freq[9][100] = 269
	detect.Big5Freq[21][67] = 268
	detect.Big5Freq[23][145] = 267
	detect.Big5Freq[10][47] = 266
	detect.Big5Freq[4][31] = 265
	detect.Big5Freq[4][81] = 264
	detect.Big5Freq[22][62] = 263
	detect.Big5Freq[4][28] = 262
	detect.Big5Freq[27][39] = 261
	detect.Big5Freq[27][54] = 260
	detect.Big5Freq[32][46] = 259
	detect.Big5Freq[4][76] = 258
	detect.Big5Freq[26][15] = 257
	detect.Big5Freq[12][154] = 256
	detect.Big5Freq[9][150] = 255
	detect.Big5Freq[15][17] = 254
	detect.Big5Freq[5][129] = 253
	detect.Big5Freq[10][40] = 252
	detect.Big5Freq[13][37] = 251
	detect.Big5Freq[31][104] = 250
	detect.Big5Freq[3][152] = 249
	detect.Big5Freq[5][22] = 248
	detect.Big5Freq[8][48] = 247
	detect.Big5Freq[4][74] = 246
	detect.Big5Freq[6][17] = 245
	detect.Big5Freq[30][82] = 244
	detect.Big5Freq[4][116] = 243
	detect.Big5Freq[16][42] = 242
	detect.Big5Freq[5][55] = 241
	detect.Big5Freq[4][64] = 240
	detect.Big5Freq[14][19] = 239
	detect.Big5Freq[35][82] = 238
	detect.Big5Freq[30][139] = 237
	detect.Big5Freq[26][152] = 236
	detect.Big5Freq[32][32] = 235
	detect.Big5Freq[21][102] = 234
	detect.Big5Freq[10][131] = 233
	detect.Big5Freq[9][128] = 232
	detect.Big5Freq[3][87] = 231
	detect.Big5Freq[4][51] = 230
	detect.Big5Freq[10][15] = 229
	detect.Big5Freq[4][150] = 228
	detect.Big5Freq[7][4] = 227
	detect.Big5Freq[7][51] = 226
	detect.Big5Freq[7][157] = 225
	detect.Big5Freq[4][146] = 224
	detect.Big5Freq[4][91] = 223
	detect.Big5Freq[7][13] = 222
	detect.Big5Freq[17][116] = 221
	detect.Big5Freq[23][21] = 220
	detect.Big5Freq[5][106] = 219
	detect.Big5Freq[14][100] = 218
	detect.Big5Freq[10][152] = 217
	detect.Big5Freq[14][89] = 216
	detect.Big5Freq[6][138] = 215
	detect.Big5Freq[12][157] = 214
	detect.Big5Freq[10][102] = 213
	detect.Big5Freq[19][94] = 212
	detect.Big5Freq[7][74] = 211
	detect.Big5Freq[18][128] = 210
	detect.Big5Freq[27][111] = 209
	detect.Big5Freq[11][57] = 208
	detect.Big5Freq[3][131] = 207
	detect.Big5Freq[30][23] = 206
	detect.Big5Freq[30][126] = 205
	detect.Big5Freq[4][36] = 204
	detect.Big5Freq[26][124] = 203
	detect.Big5Freq[4][19] = 202
	detect.Big5Freq[9][152] = 201
	/*
	 * Big5Freq[5][0] = 200; Big5Freq[26][57] = 199; Big5Freq[13][155] =
	 * 198; Big5Freq[3][38] = 197; Big5Freq[9][155] = 196; Big5Freq[28][53]
	 * = 195; Big5Freq[15][71] = 194; Big5Freq[21][95] = 193;
	 * Big5Freq[15][112] = 192; Big5Freq[14][138] = 191; Big5Freq[8][18] =
	 * 190; Big5Freq[20][151] = 189; Big5Freq[37][27] = 188;
	 * Big5Freq[32][48] = 187; Big5Freq[23][66] = 186; Big5Freq[9][2] = 185;
	 * Big5Freq[13][133] = 184; Big5Freq[7][127] = 183; Big5Freq[3][11] =
	 * 182; Big5Freq[12][118] = 181; Big5Freq[13][101] = 180;
	 * Big5Freq[30][153] = 179; Big5Freq[4][65] = 178; Big5Freq[5][25] =
	 * 177; Big5Freq[5][140] = 176; Big5Freq[6][25] = 175; Big5Freq[4][52] =
	 * 174; Big5Freq[30][156] = 173; Big5Freq[16][13] = 172; Big5Freq[21][8]
	 * = 171; Big5Freq[19][74] = 170; Big5Freq[15][145] = 169;
	 * Big5Freq[9][15] = 168; Big5Freq[13][82] = 167; Big5Freq[26][86] =
	 * 166; Big5Freq[18][52] = 165; Big5Freq[6][109] = 164; Big5Freq[10][99]
	 * = 163; Big5Freq[18][101] = 162; Big5Freq[25][49] = 161;
	 * Big5Freq[31][79] = 160; Big5Freq[28][20] = 159; Big5Freq[12][115] =
	 * 158; Big5Freq[15][66] = 157; Big5Freq[11][104] = 156;
	 * Big5Freq[23][106] = 155; Big5Freq[34][157] = 154; Big5Freq[32][94] =
	 * 153; Big5Freq[29][88] = 152; Big5Freq[10][46] = 151;
	 * Big5Freq[13][118] = 150; Big5Freq[20][37] = 149; Big5Freq[12][30] =
	 * 148; Big5Freq[21][4] = 147; Big5Freq[16][33] = 146; Big5Freq[13][52]
	 * = 145; Big5Freq[4][7] = 144; Big5Freq[21][49] = 143; Big5Freq[3][27]
	 * = 142; Big5Freq[16][91] = 141; Big5Freq[5][155] = 140;
	 * Big5Freq[29][130] = 139; Big5Freq[3][125] = 138; Big5Freq[14][26] =
	 * 137; Big5Freq[15][39] = 136; Big5Freq[24][110] = 135;
	 * Big5Freq[7][141] = 134; Big5Freq[21][15] = 133; Big5Freq[32][104] =
	 * 132; Big5Freq[8][31] = 131; Big5Freq[34][112] = 130; Big5Freq[10][75]
	 * = 129; Big5Freq[21][23] = 128; Big5Freq[34][131] = 127;
	 * Big5Freq[12][3] = 126; Big5Freq[10][62] = 125; Big5Freq[9][120] =
	 * 124; Big5Freq[32][149] = 123; Big5Freq[8][44] = 122; Big5Freq[24][2]
	 * = 121; Big5Freq[6][148] = 120; Big5Freq[15][103] = 119;
	 * Big5Freq[36][54] = 118; Big5Freq[36][134] = 117; Big5Freq[11][7] =
	 * 116; Big5Freq[3][90] = 115; Big5Freq[36][73] = 114; Big5Freq[8][102]
	 * = 113; Big5Freq[12][87] = 112; Big5Freq[25][64] = 111; Big5Freq[9][1]
	 * = 110; Big5Freq[24][121] = 109; Big5Freq[5][75] = 108;
	 * Big5Freq[17][83] = 107; Big5Freq[18][57] = 106; Big5Freq[8][95] =
	 * 105; Big5Freq[14][36] = 104; Big5Freq[28][113] = 103;
	 * Big5Freq[12][56] = 102; Big5Freq[14][61] = 101; Big5Freq[25][138] =
	 * 100; Big5Freq[4][34] = 99; Big5Freq[11][152] = 98; Big5Freq[35][0] =
	 * 97; Big5Freq[4][15] = 96; Big5Freq[8][82] = 95; Big5Freq[20][73] =
	 * 94; Big5Freq[25][52] = 93; Big5Freq[24][6] = 92; Big5Freq[21][78] =
	 * 91; Big5Freq[17][32] = 90; Big5Freq[17][91] = 89; Big5Freq[5][76] =
	 * 88; Big5Freq[15][60] = 87; Big5Freq[15][150] = 86; Big5Freq[5][80] =
	 * 85; Big5Freq[15][81] = 84; Big5Freq[28][108] = 83; Big5Freq[18][14] =
	 * 82; Big5Freq[19][109] = 81; Big5Freq[28][133] = 80; Big5Freq[21][97]
	 * = 79; Big5Freq[5][105] = 78; Big5Freq[18][114] = 77; Big5Freq[16][95]
	 * = 76; Big5Freq[5][51] = 75; Big5Freq[3][148] = 74; Big5Freq[22][102]
	 * = 73; Big5Freq[4][123] = 72; Big5Freq[8][88] = 71; Big5Freq[25][111]
	 * = 70; Big5Freq[8][149] = 69; Big5Freq[9][48] = 68; Big5Freq[16][126]
	 * = 67; Big5Freq[33][150] = 66; Big5Freq[9][54] = 65; Big5Freq[29][104]
	 * = 64; Big5Freq[3][3] = 63; Big5Freq[11][49] = 62; Big5Freq[24][109] =
	 * 61; Big5Freq[28][116] = 60; Big5Freq[34][113] = 59; Big5Freq[5][3] =
	 * 58; Big5Freq[21][106] = 57; Big5Freq[4][98] = 56; Big5Freq[12][135] =
	 * 55; Big5Freq[16][101] = 54; Big5Freq[12][147] = 53; Big5Freq[27][55]
	 * = 52; Big5Freq[3][5] = 51; Big5Freq[11][101] = 50; Big5Freq[16][157]
	 * = 49; Big5Freq[22][114] = 48; Big5Freq[18][46] = 47; Big5Freq[4][29]
	 * = 46; Big5Freq[8][103] = 45; Big5Freq[16][151] = 44; Big5Freq[8][29]
	 * = 43; Big5Freq[15][114] = 42; Big5Freq[22][70] = 41;
	 * Big5Freq[13][121] = 40; Big5Freq[7][112] = 39; Big5Freq[20][83] = 38;
	 * Big5Freq[3][36] = 37; Big5Freq[10][103] = 36; Big5Freq[3][96] = 35;
	 * Big5Freq[21][79] = 34; Big5Freq[25][120] = 33; Big5Freq[29][121] =
	 * 32; Big5Freq[23][71] = 31; Big5Freq[21][22] = 30; Big5Freq[18][89] =
	 * 29; Big5Freq[25][104] = 28; Big5Freq[10][124] = 27; Big5Freq[26][4] =
	 * 26; Big5Freq[21][136] = 25; Big5Freq[6][112] = 24; Big5Freq[12][103]
	 * = 23; Big5Freq[17][66] = 22; Big5Freq[13][151] = 21;
	 * Big5Freq[33][152] = 20; Big5Freq[11][148] = 19; Big5Freq[13][57] =
	 * 18; Big5Freq[13][41] = 17; Big5Freq[7][60] = 16; Big5Freq[21][29] =
	 * 15; Big5Freq[9][157] = 14; Big5Freq[24][95] = 13; Big5Freq[15][148] =
	 * 12; Big5Freq[15][122] = 11; Big5Freq[6][125] = 10; Big5Freq[11][25] =
	 * 9; Big5Freq[20][55] = 8; Big5Freq[19][84] = 7; Big5Freq[21][82] = 6;
	 * Big5Freq[24][3] = 5; Big5Freq[13][70] = 4; Big5Freq[6][21] = 3;
	 * Big5Freq[21][86] = 2; Big5Freq[12][23] = 1; Big5Freq[3][85] = 0;
	 * EucTwfreq[45][90] = 600;
	 */
	detect.Big5PFreq[41][122] = 600
	detect.Big5PFreq[35][0] = 599
	detect.Big5PFreq[43][15] = 598
	detect.Big5PFreq[35][99] = 597
	detect.Big5PFreq[35][6] = 596
	detect.Big5PFreq[35][8] = 595
	detect.Big5PFreq[38][154] = 594
	detect.Big5PFreq[37][34] = 593
	detect.Big5PFreq[37][115] = 592
	detect.Big5PFreq[36][12] = 591
	detect.Big5PFreq[18][77] = 590
	detect.Big5PFreq[35][100] = 589
	detect.Big5PFreq[35][42] = 588
	detect.Big5PFreq[120][75] = 587
	detect.Big5PFreq[35][23] = 586
	detect.Big5PFreq[13][72] = 585
	detect.Big5PFreq[0][67] = 584
	detect.Big5PFreq[39][172] = 583
	detect.Big5PFreq[22][182] = 582
	detect.Big5PFreq[15][186] = 581
	detect.Big5PFreq[15][165] = 580
	detect.Big5PFreq[35][44] = 579
	detect.Big5PFreq[40][13] = 578
	detect.Big5PFreq[38][1] = 577
	detect.Big5PFreq[37][33] = 576
	detect.Big5PFreq[36][24] = 575
	detect.Big5PFreq[56][4] = 574
	detect.Big5PFreq[35][29] = 573
	detect.Big5PFreq[9][96] = 572
	detect.Big5PFreq[37][62] = 571
	detect.Big5PFreq[48][47] = 570
	detect.Big5PFreq[51][14] = 569
	detect.Big5PFreq[39][122] = 568
	detect.Big5PFreq[44][46] = 567
	detect.Big5PFreq[35][21] = 566
	detect.Big5PFreq[36][8] = 565
	detect.Big5PFreq[36][141] = 564
	detect.Big5PFreq[3][81] = 563
	detect.Big5PFreq[37][155] = 562
	detect.Big5PFreq[42][84] = 561
	detect.Big5PFreq[36][40] = 560
	detect.Big5PFreq[35][103] = 559
	detect.Big5PFreq[11][84] = 558
	detect.Big5PFreq[45][33] = 557
	detect.Big5PFreq[121][79] = 556
	detect.Big5PFreq[2][77] = 555
	detect.Big5PFreq[36][41] = 554
	detect.Big5PFreq[37][47] = 553
	detect.Big5PFreq[39][125] = 552
	detect.Big5PFreq[37][26] = 551
	detect.Big5PFreq[35][48] = 550
	detect.Big5PFreq[35][28] = 549
	detect.Big5PFreq[35][159] = 548
	detect.Big5PFreq[37][40] = 547
	detect.Big5PFreq[35][145] = 546
	detect.Big5PFreq[37][147] = 545
	detect.Big5PFreq[46][160] = 544
	detect.Big5PFreq[37][46] = 543
	detect.Big5PFreq[50][99] = 542
	detect.Big5PFreq[52][13] = 541
	detect.Big5PFreq[10][82] = 540
	detect.Big5PFreq[35][169] = 539
	detect.Big5PFreq[35][31] = 538
	detect.Big5PFreq[47][31] = 537
	detect.Big5PFreq[18][79] = 536
	detect.Big5PFreq[16][113] = 535
	detect.Big5PFreq[37][104] = 534
	detect.Big5PFreq[39][134] = 533
	detect.Big5PFreq[36][53] = 532
	detect.Big5PFreq[38][0] = 531
	detect.Big5PFreq[4][86] = 530
	detect.Big5PFreq[54][17] = 529
	detect.Big5PFreq[43][157] = 528
	detect.Big5PFreq[35][165] = 527
	detect.Big5PFreq[69][147] = 526
	detect.Big5PFreq[117][95] = 525
	detect.Big5PFreq[35][162] = 524
	detect.Big5PFreq[35][17] = 523
	detect.Big5PFreq[36][142] = 522
	detect.Big5PFreq[36][4] = 521
	detect.Big5PFreq[37][166] = 520
	detect.Big5PFreq[35][168] = 519
	detect.Big5PFreq[35][19] = 518
	detect.Big5PFreq[37][48] = 517
	detect.Big5PFreq[42][37] = 516
	detect.Big5PFreq[40][146] = 515
	detect.Big5PFreq[36][123] = 514
	detect.Big5PFreq[22][41] = 513
	detect.Big5PFreq[20][119] = 512
	detect.Big5PFreq[2][74] = 511
	detect.Big5PFreq[44][113] = 510
	detect.Big5PFreq[35][125] = 509
	detect.Big5PFreq[37][16] = 508
	detect.Big5PFreq[35][20] = 507
	detect.Big5PFreq[35][55] = 506
	detect.Big5PFreq[37][145] = 505
	detect.Big5PFreq[0][88] = 504
	detect.Big5PFreq[3][94] = 503
	detect.Big5PFreq[6][65] = 502
	detect.Big5PFreq[26][15] = 501
	detect.Big5PFreq[41][126] = 500
	detect.Big5PFreq[36][129] = 499
	detect.Big5PFreq[31][75] = 498
	detect.Big5PFreq[19][61] = 497
	detect.Big5PFreq[35][128] = 496
	detect.Big5PFreq[29][79] = 495
	detect.Big5PFreq[36][62] = 494
	detect.Big5PFreq[37][189] = 493
	detect.Big5PFreq[39][109] = 492
	detect.Big5PFreq[39][135] = 491
	detect.Big5PFreq[72][15] = 490
	detect.Big5PFreq[47][106] = 489
	detect.Big5PFreq[54][14] = 488
	detect.Big5PFreq[24][52] = 487
	detect.Big5PFreq[38][162] = 486
	detect.Big5PFreq[41][43] = 485
	detect.Big5PFreq[37][121] = 484
	detect.Big5PFreq[14][66] = 483
	detect.Big5PFreq[37][30] = 482
	detect.Big5PFreq[35][7] = 481
	detect.Big5PFreq[49][58] = 480
	detect.Big5PFreq[43][188] = 479
	detect.Big5PFreq[24][66] = 478
	detect.Big5PFreq[35][171] = 477
	detect.Big5PFreq[40][186] = 476
	detect.Big5PFreq[39][164] = 475
	detect.Big5PFreq[78][186] = 474
	detect.Big5PFreq[8][72] = 473
	detect.Big5PFreq[36][190] = 472
	detect.Big5PFreq[35][53] = 471
	detect.Big5PFreq[35][54] = 470
	detect.Big5PFreq[22][159] = 469
	detect.Big5PFreq[35][9] = 468
	detect.Big5PFreq[41][140] = 467
	detect.Big5PFreq[37][22] = 466
	detect.Big5PFreq[48][97] = 465
	detect.Big5PFreq[50][97] = 464
	detect.Big5PFreq[36][127] = 463
	detect.Big5PFreq[37][23] = 462
	detect.Big5PFreq[40][55] = 461
	detect.Big5PFreq[35][43] = 460
	detect.Big5PFreq[26][22] = 459
	detect.Big5PFreq[35][15] = 458
	detect.Big5PFreq[72][179] = 457
	detect.Big5PFreq[20][129] = 456
	detect.Big5PFreq[52][101] = 455
	detect.Big5PFreq[35][12] = 454
	detect.Big5PFreq[42][156] = 453
	detect.Big5PFreq[15][157] = 452
	detect.Big5PFreq[50][140] = 451
	detect.Big5PFreq[26][28] = 450
	detect.Big5PFreq[54][51] = 449
	detect.Big5PFreq[35][112] = 448
	detect.Big5PFreq[36][116] = 447
	detect.Big5PFreq[42][11] = 446
	detect.Big5PFreq[37][172] = 445
	detect.Big5PFreq[37][29] = 444
	detect.Big5PFreq[44][107] = 443
	detect.Big5PFreq[50][17] = 442
	detect.Big5PFreq[39][107] = 441
	detect.Big5PFreq[19][109] = 440
	detect.Big5PFreq[36][60] = 439
	detect.Big5PFreq[49][132] = 438
	detect.Big5PFreq[26][16] = 437
	detect.Big5PFreq[43][155] = 436
	detect.Big5PFreq[37][120] = 435
	detect.Big5PFreq[15][159] = 434
	detect.Big5PFreq[43][6] = 433
	detect.Big5PFreq[45][188] = 432
	detect.Big5PFreq[35][38] = 431
	detect.Big5PFreq[39][143] = 430
	detect.Big5PFreq[48][144] = 429
	detect.Big5PFreq[37][168] = 428
	detect.Big5PFreq[37][1] = 427
	detect.Big5PFreq[36][109] = 426
	detect.Big5PFreq[46][53] = 425
	detect.Big5PFreq[38][54] = 424
	detect.Big5PFreq[36][0] = 423
	detect.Big5PFreq[72][33] = 422
	detect.Big5PFreq[42][8] = 421
	detect.Big5PFreq[36][31] = 420
	detect.Big5PFreq[35][150] = 419
	detect.Big5PFreq[118][93] = 418
	detect.Big5PFreq[37][61] = 417
	detect.Big5PFreq[0][85] = 416
	detect.Big5PFreq[36][27] = 415
	detect.Big5PFreq[35][134] = 414
	detect.Big5PFreq[36][145] = 413
	detect.Big5PFreq[6][96] = 412
	detect.Big5PFreq[36][14] = 411
	detect.Big5PFreq[16][36] = 410
	detect.Big5PFreq[15][175] = 409
	detect.Big5PFreq[35][10] = 408
	detect.Big5PFreq[36][189] = 407
	detect.Big5PFreq[35][51] = 406
	detect.Big5PFreq[35][109] = 405
	detect.Big5PFreq[35][147] = 404
	detect.Big5PFreq[35][180] = 403
	detect.Big5PFreq[72][5] = 402
	detect.Big5PFreq[36][107] = 401
	detect.Big5PFreq[49][116] = 400
	detect.Big5PFreq[73][30] = 399
	detect.Big5PFreq[6][90] = 398
	detect.Big5PFreq[2][70] = 397
	detect.Big5PFreq[17][141] = 396
	detect.Big5PFreq[35][62] = 395
	detect.Big5PFreq[16][180] = 394
	detect.Big5PFreq[4][91] = 393
	detect.Big5PFreq[15][171] = 392
	detect.Big5PFreq[35][177] = 391
	detect.Big5PFreq[37][173] = 390
	detect.Big5PFreq[16][121] = 389
	detect.Big5PFreq[35][5] = 388
	detect.Big5PFreq[46][122] = 387
	detect.Big5PFreq[40][138] = 386
	detect.Big5PFreq[50][49] = 385
	detect.Big5PFreq[36][152] = 384
	detect.Big5PFreq[13][43] = 383
	detect.Big5PFreq[9][88] = 382
	detect.Big5PFreq[36][159] = 381
	detect.Big5PFreq[27][62] = 380
	detect.Big5PFreq[40][18] = 379
	detect.Big5PFreq[17][129] = 378
	detect.Big5PFreq[43][97] = 377
	detect.Big5PFreq[13][131] = 376
	detect.Big5PFreq[46][107] = 375
	detect.Big5PFreq[60][64] = 374
	detect.Big5PFreq[36][179] = 373
	detect.Big5PFreq[37][55] = 372
	detect.Big5PFreq[41][173] = 371
	detect.Big5PFreq[44][172] = 370
	detect.Big5PFreq[23][187] = 369
	detect.Big5PFreq[36][149] = 368
	detect.Big5PFreq[17][125] = 367
	detect.Big5PFreq[55][180] = 366
	detect.Big5PFreq[51][129] = 365
	detect.Big5PFreq[36][51] = 364
	detect.Big5PFreq[37][122] = 363
	detect.Big5PFreq[48][32] = 362
	detect.Big5PFreq[51][99] = 361
	detect.Big5PFreq[54][16] = 360
	detect.Big5PFreq[41][183] = 359
	detect.Big5PFreq[37][179] = 358
	detect.Big5PFreq[38][179] = 357
	detect.Big5PFreq[35][143] = 356
	detect.Big5PFreq[37][24] = 355
	detect.Big5PFreq[40][177] = 354
	detect.Big5PFreq[47][117] = 353
	detect.Big5PFreq[39][52] = 352
	detect.Big5PFreq[22][99] = 351
	detect.Big5PFreq[40][142] = 350
	detect.Big5PFreq[36][49] = 349
	detect.Big5PFreq[38][17] = 348
	detect.Big5PFreq[39][188] = 347
	detect.Big5PFreq[36][186] = 346
	detect.Big5PFreq[35][189] = 345
	detect.Big5PFreq[41][7] = 344
	detect.Big5PFreq[18][91] = 343
	detect.Big5PFreq[43][137] = 342
	detect.Big5PFreq[35][142] = 341
	detect.Big5PFreq[35][117] = 340
	detect.Big5PFreq[39][138] = 339
	detect.Big5PFreq[16][59] = 338
	detect.Big5PFreq[39][174] = 337
	detect.Big5PFreq[55][145] = 336
	detect.Big5PFreq[37][21] = 335
	detect.Big5PFreq[36][180] = 334
	detect.Big5PFreq[37][156] = 333
	detect.Big5PFreq[49][13] = 332
	detect.Big5PFreq[41][107] = 331
	detect.Big5PFreq[36][56] = 330
	detect.Big5PFreq[53][8] = 329
	detect.Big5PFreq[22][114] = 328
	detect.Big5PFreq[5][95] = 327
	detect.Big5PFreq[37][0] = 326
	detect.Big5PFreq[26][183] = 325
	detect.Big5PFreq[22][66] = 324
	detect.Big5PFreq[35][58] = 323
	detect.Big5PFreq[48][117] = 322
	detect.Big5PFreq[36][102] = 321
	detect.Big5PFreq[22][122] = 320
	detect.Big5PFreq[35][11] = 319
	detect.Big5PFreq[46][19] = 318
	detect.Big5PFreq[22][49] = 317
	detect.Big5PFreq[48][166] = 316
	detect.Big5PFreq[41][125] = 315
	detect.Big5PFreq[41][1] = 314
	detect.Big5PFreq[35][178] = 313
	detect.Big5PFreq[41][12] = 312
	detect.Big5PFreq[26][167] = 311
	detect.Big5PFreq[42][152] = 310
	detect.Big5PFreq[42][46] = 309
	detect.Big5PFreq[42][151] = 308
	detect.Big5PFreq[20][135] = 307
	detect.Big5PFreq[37][162] = 306
	detect.Big5PFreq[37][50] = 305
	detect.Big5PFreq[22][185] = 304
	detect.Big5PFreq[36][166] = 303
	detect.Big5PFreq[19][40] = 302
	detect.Big5PFreq[22][107] = 301
	detect.Big5PFreq[22][102] = 300
	detect.Big5PFreq[57][162] = 299
	detect.Big5PFreq[22][124] = 298
	detect.Big5PFreq[37][138] = 297
	detect.Big5PFreq[37][25] = 296
	detect.Big5PFreq[0][69] = 295
	detect.Big5PFreq[43][172] = 294
	detect.Big5PFreq[42][167] = 293
	detect.Big5PFreq[35][120] = 292
	detect.Big5PFreq[41][128] = 291
	detect.Big5PFreq[2][88] = 290
	detect.Big5PFreq[20][123] = 289
	detect.Big5PFreq[35][123] = 288
	detect.Big5PFreq[36][28] = 287
	detect.Big5PFreq[42][188] = 286
	detect.Big5PFreq[42][164] = 285
	detect.Big5PFreq[42][4] = 284
	detect.Big5PFreq[43][57] = 283
	detect.Big5PFreq[39][3] = 282
	detect.Big5PFreq[42][3] = 281
	detect.Big5PFreq[57][158] = 280
	detect.Big5PFreq[35][146] = 279
	detect.Big5PFreq[24][54] = 278
	detect.Big5PFreq[13][110] = 277
	detect.Big5PFreq[23][132] = 276
	detect.Big5PFreq[26][102] = 275
	detect.Big5PFreq[55][178] = 274
	detect.Big5PFreq[17][117] = 273
	detect.Big5PFreq[41][161] = 272
	detect.Big5PFreq[38][150] = 271
	detect.Big5PFreq[10][71] = 270
	detect.Big5PFreq[47][60] = 269
	detect.Big5PFreq[16][114] = 268
	detect.Big5PFreq[21][47] = 267
	detect.Big5PFreq[39][101] = 266
	detect.Big5PFreq[18][45] = 265
	detect.Big5PFreq[40][121] = 264
	detect.Big5PFreq[45][41] = 263
	detect.Big5PFreq[22][167] = 262
	detect.Big5PFreq[26][149] = 261
	detect.Big5PFreq[15][189] = 260
	detect.Big5PFreq[41][177] = 259
	detect.Big5PFreq[46][36] = 258
	detect.Big5PFreq[20][40] = 257
	detect.Big5PFreq[41][54] = 256
	detect.Big5PFreq[3][87] = 255
	detect.Big5PFreq[40][16] = 254
	detect.Big5PFreq[42][15] = 253
	detect.Big5PFreq[11][83] = 252
	detect.Big5PFreq[0][94] = 251
	detect.Big5PFreq[122][81] = 250
	detect.Big5PFreq[41][26] = 249
	detect.Big5PFreq[36][34] = 248
	detect.Big5PFreq[44][148] = 247
	detect.Big5PFreq[35][3] = 246
	detect.Big5PFreq[36][114] = 245
	detect.Big5PFreq[42][112] = 244
	detect.Big5PFreq[35][183] = 243
	detect.Big5PFreq[49][73] = 242
	detect.Big5PFreq[39][2] = 241
	detect.Big5PFreq[38][121] = 240
	detect.Big5PFreq[44][114] = 239
	detect.Big5PFreq[49][32] = 238
	detect.Big5PFreq[1][65] = 237
	detect.Big5PFreq[38][25] = 236
	detect.Big5PFreq[39][4] = 235
	detect.Big5PFreq[42][62] = 234
	detect.Big5PFreq[35][40] = 233
	detect.Big5PFreq[24][2] = 232
	detect.Big5PFreq[53][49] = 231
	detect.Big5PFreq[41][133] = 230
	detect.Big5PFreq[43][134] = 229
	detect.Big5PFreq[3][83] = 228
	detect.Big5PFreq[38][158] = 227
	detect.Big5PFreq[24][17] = 226
	detect.Big5PFreq[52][59] = 225
	detect.Big5PFreq[38][41] = 224
	detect.Big5PFreq[37][127] = 223
	detect.Big5PFreq[22][175] = 222
	detect.Big5PFreq[44][30] = 221
	detect.Big5PFreq[47][178] = 220
	detect.Big5PFreq[43][99] = 219
	detect.Big5PFreq[19][4] = 218
	detect.Big5PFreq[37][97] = 217
	detect.Big5PFreq[38][181] = 216
	detect.Big5PFreq[45][103] = 215
	detect.Big5PFreq[1][86] = 214
	detect.Big5PFreq[40][15] = 213
	detect.Big5PFreq[22][136] = 212
	detect.Big5PFreq[75][165] = 211
	detect.Big5PFreq[36][15] = 210
	detect.Big5PFreq[46][80] = 209
	detect.Big5PFreq[59][55] = 208
	detect.Big5PFreq[37][108] = 207
	detect.Big5PFreq[21][109] = 206
	detect.Big5PFreq[24][165] = 205
	detect.Big5PFreq[79][158] = 204
	detect.Big5PFreq[44][139] = 203
	detect.Big5PFreq[36][124] = 202
	detect.Big5PFreq[42][185] = 201
	detect.Big5PFreq[39][186] = 200
	detect.Big5PFreq[22][128] = 199
	detect.Big5PFreq[40][44] = 198
	detect.Big5PFreq[41][105] = 197
	detect.Big5PFreq[1][70] = 196
	detect.Big5PFreq[1][68] = 195
	detect.Big5PFreq[53][22] = 194
	detect.Big5PFreq[36][54] = 193
	detect.Big5PFreq[47][147] = 192
	detect.Big5PFreq[35][36] = 191
	detect.Big5PFreq[35][185] = 190
	detect.Big5PFreq[45][37] = 189
	detect.Big5PFreq[43][163] = 188
	detect.Big5PFreq[56][115] = 187
	detect.Big5PFreq[38][164] = 186
	detect.Big5PFreq[35][141] = 185
	detect.Big5PFreq[42][132] = 184
	detect.Big5PFreq[46][120] = 183
	detect.Big5PFreq[69][142] = 182
	detect.Big5PFreq[38][175] = 181
	detect.Big5PFreq[22][112] = 180
	detect.Big5PFreq[38][142] = 179
	detect.Big5PFreq[40][37] = 178
	detect.Big5PFreq[37][109] = 177
	detect.Big5PFreq[40][144] = 176
	detect.Big5PFreq[44][117] = 175
	detect.Big5PFreq[35][181] = 174
	detect.Big5PFreq[26][105] = 173
	detect.Big5PFreq[16][48] = 172
	detect.Big5PFreq[44][122] = 171
	detect.Big5PFreq[12][86] = 170
	detect.Big5PFreq[84][53] = 169
	detect.Big5PFreq[17][44] = 168
	detect.Big5PFreq[59][54] = 167
	detect.Big5PFreq[36][98] = 166
	detect.Big5PFreq[45][115] = 165
	detect.Big5PFreq[73][9] = 164
	detect.Big5PFreq[44][123] = 163
	detect.Big5PFreq[37][188] = 162
	detect.Big5PFreq[51][117] = 161
	detect.Big5PFreq[15][156] = 160
	detect.Big5PFreq[36][155] = 159
	detect.Big5PFreq[44][25] = 158
	detect.Big5PFreq[38][12] = 157
	detect.Big5PFreq[38][140] = 156
	detect.Big5PFreq[23][4] = 155
	detect.Big5PFreq[45][149] = 154
	detect.Big5PFreq[22][189] = 153
	detect.Big5PFreq[38][147] = 152
	detect.Big5PFreq[27][5] = 151
	detect.Big5PFreq[22][42] = 150
	detect.Big5PFreq[3][68] = 149
	detect.Big5PFreq[39][51] = 148
	detect.Big5PFreq[36][29] = 147
	detect.Big5PFreq[20][108] = 146
	detect.Big5PFreq[50][57] = 145
	detect.Big5PFreq[55][104] = 144
	detect.Big5PFreq[22][46] = 143
	detect.Big5PFreq[18][164] = 142
	detect.Big5PFreq[50][159] = 141
	detect.Big5PFreq[85][131] = 140
	detect.Big5PFreq[26][79] = 139
	detect.Big5PFreq[38][100] = 138
	detect.Big5PFreq[53][112] = 137
	detect.Big5PFreq[20][190] = 136
	detect.Big5PFreq[14][69] = 135
	detect.Big5PFreq[23][11] = 134
	detect.Big5PFreq[40][114] = 133
	detect.Big5PFreq[40][148] = 132
	detect.Big5PFreq[53][130] = 131
	detect.Big5PFreq[36][2] = 130
	detect.Big5PFreq[66][82] = 129
	detect.Big5PFreq[45][166] = 128
	detect.Big5PFreq[4][88] = 127
	detect.Big5PFreq[16][57] = 126
	detect.Big5PFreq[22][116] = 125
	detect.Big5PFreq[36][108] = 124
	detect.Big5PFreq[13][48] = 123
	detect.Big5PFreq[54][12] = 122
	detect.Big5PFreq[40][136] = 121
	detect.Big5PFreq[36][128] = 120
	detect.Big5PFreq[23][6] = 119
	detect.Big5PFreq[38][125] = 118
	detect.Big5PFreq[45][154] = 117
	detect.Big5PFreq[51][127] = 116
	detect.Big5PFreq[44][163] = 115
	detect.Big5PFreq[16][173] = 114
	detect.Big5PFreq[43][49] = 113
	detect.Big5PFreq[20][112] = 112
	detect.Big5PFreq[15][168] = 111
	detect.Big5PFreq[35][129] = 110
	detect.Big5PFreq[20][45] = 109
	detect.Big5PFreq[38][10] = 108
	detect.Big5PFreq[57][171] = 107
	detect.Big5PFreq[44][190] = 106
	detect.Big5PFreq[40][56] = 105
	detect.Big5PFreq[36][156] = 104
	detect.Big5PFreq[3][88] = 103
	detect.Big5PFreq[50][122] = 102
	detect.Big5PFreq[36][7] = 101
	detect.Big5PFreq[39][43] = 100
	detect.Big5PFreq[15][166] = 99
	detect.Big5PFreq[42][136] = 98
	detect.Big5PFreq[22][131] = 97
	detect.Big5PFreq[44][23] = 96
	detect.Big5PFreq[54][147] = 95
	detect.Big5PFreq[41][32] = 94
	detect.Big5PFreq[23][121] = 93
	detect.Big5PFreq[39][108] = 92
	detect.Big5PFreq[2][78] = 91
	detect.Big5PFreq[40][155] = 90
	detect.Big5PFreq[55][51] = 89
	detect.Big5PFreq[19][34] = 88
	detect.Big5PFreq[48][128] = 87
	detect.Big5PFreq[48][159] = 86
	detect.Big5PFreq[20][70] = 85
	detect.Big5PFreq[34][71] = 84
	detect.Big5PFreq[16][31] = 83
	detect.Big5PFreq[42][157] = 82
	detect.Big5PFreq[20][44] = 81
	detect.Big5PFreq[11][92] = 80
	detect.Big5PFreq[44][180] = 79
	detect.Big5PFreq[84][33] = 78
	detect.Big5PFreq[16][116] = 77
	detect.Big5PFreq[61][163] = 76
	detect.Big5PFreq[35][164] = 75
	detect.Big5PFreq[36][42] = 74
	detect.Big5PFreq[13][40] = 73
	detect.Big5PFreq[43][176] = 72
	detect.Big5PFreq[2][66] = 71
	detect.Big5PFreq[20][133] = 70
	detect.Big5PFreq[36][65] = 69
	detect.Big5PFreq[38][33] = 68
	detect.Big5PFreq[12][91] = 67
	detect.Big5PFreq[36][26] = 66
	detect.Big5PFreq[15][174] = 65
	detect.Big5PFreq[77][32] = 64
	detect.Big5PFreq[16][1] = 63
	detect.Big5PFreq[25][86] = 62
	detect.Big5PFreq[17][13] = 61
	detect.Big5PFreq[5][75] = 60
	detect.Big5PFreq[36][52] = 59
	detect.Big5PFreq[51][164] = 58
	detect.Big5PFreq[12][85] = 57
	detect.Big5PFreq[39][168] = 56
	detect.Big5PFreq[43][16] = 55
	detect.Big5PFreq[40][69] = 54
	detect.Big5PFreq[26][108] = 53
	detect.Big5PFreq[51][56] = 52
	detect.Big5PFreq[16][37] = 51
	detect.Big5PFreq[40][29] = 50
	detect.Big5PFreq[46][171] = 49
	detect.Big5PFreq[40][128] = 48
	detect.Big5PFreq[72][114] = 47
	detect.Big5PFreq[21][103] = 46
	detect.Big5PFreq[22][44] = 45
	detect.Big5PFreq[40][115] = 44
	detect.Big5PFreq[43][7] = 43
	detect.Big5PFreq[43][153] = 42
	detect.Big5PFreq[17][20] = 41
	detect.Big5PFreq[16][49] = 40
	detect.Big5PFreq[36][57] = 39
	detect.Big5PFreq[18][38] = 38
	detect.Big5PFreq[45][184] = 37
	detect.Big5PFreq[37][167] = 36
	detect.Big5PFreq[26][106] = 35
	detect.Big5PFreq[61][121] = 34
	detect.Big5PFreq[89][140] = 33
	detect.Big5PFreq[46][61] = 32
	detect.Big5PFreq[39][163] = 31
	detect.Big5PFreq[40][62] = 30
	detect.Big5PFreq[38][165] = 29
	detect.Big5PFreq[47][37] = 28
	detect.Big5PFreq[18][155] = 27
	detect.Big5PFreq[20][33] = 26
	detect.Big5PFreq[29][90] = 25
	detect.Big5PFreq[20][103] = 24
	detect.Big5PFreq[37][51] = 23
	detect.Big5PFreq[57][0] = 22
	detect.Big5PFreq[40][31] = 21
	detect.Big5PFreq[45][32] = 20
	detect.Big5PFreq[59][23] = 19
	detect.Big5PFreq[18][47] = 18
	detect.Big5PFreq[45][134] = 17
	detect.Big5PFreq[37][59] = 16
	detect.Big5PFreq[21][128] = 15
	detect.Big5PFreq[36][106] = 14
	detect.Big5PFreq[31][39] = 13
	detect.Big5PFreq[40][182] = 12
	detect.Big5PFreq[52][155] = 11
	detect.Big5PFreq[42][166] = 10
	detect.Big5PFreq[35][27] = 9
	detect.Big5PFreq[38][3] = 8
	detect.Big5PFreq[13][44] = 7
	detect.Big5PFreq[58][157] = 6
	detect.Big5PFreq[47][51] = 5
	detect.Big5PFreq[41][37] = 4
	detect.Big5PFreq[41][172] = 3
	detect.Big5PFreq[51][165] = 2
	detect.Big5PFreq[15][161] = 1
	detect.Big5PFreq[24][181] = 0
	detect.EucTwfreq[48][49] = 599
	detect.EucTwfreq[35][65] = 598
	detect.EucTwfreq[41][27] = 597
	detect.EucTwfreq[35][0] = 596
	detect.EucTwfreq[39][19] = 595
	detect.EucTwfreq[35][42] = 594
	detect.EucTwfreq[38][66] = 593
	detect.EucTwfreq[35][8] = 592
	detect.EucTwfreq[35][6] = 591
	detect.EucTwfreq[35][66] = 590
	detect.EucTwfreq[43][14] = 589
	detect.EucTwfreq[69][80] = 588
	detect.EucTwfreq[50][48] = 587
	detect.EucTwfreq[36][71] = 586
	detect.EucTwfreq[37][10] = 585
	detect.EucTwfreq[60][52] = 584
	detect.EucTwfreq[51][21] = 583
	detect.EucTwfreq[40][2] = 582
	detect.EucTwfreq[67][35] = 581
	detect.EucTwfreq[38][78] = 580
	detect.EucTwfreq[49][18] = 579
	detect.EucTwfreq[35][23] = 578
	detect.EucTwfreq[42][83] = 577
	detect.EucTwfreq[79][47] = 576
	detect.EucTwfreq[61][82] = 575
	detect.EucTwfreq[38][7] = 574
	detect.EucTwfreq[35][29] = 573
	detect.EucTwfreq[37][77] = 572
	detect.EucTwfreq[54][67] = 571
	detect.EucTwfreq[38][80] = 570
	detect.EucTwfreq[52][74] = 569
	detect.EucTwfreq[36][37] = 568
	detect.EucTwfreq[74][8] = 567
	detect.EucTwfreq[41][83] = 566
	detect.EucTwfreq[36][75] = 565
	detect.EucTwfreq[49][63] = 564
	detect.EucTwfreq[42][58] = 563
	detect.EucTwfreq[56][33] = 562
	detect.EucTwfreq[37][76] = 561
	detect.EucTwfreq[62][39] = 560
	detect.EucTwfreq[35][21] = 559
	detect.EucTwfreq[70][19] = 558
	detect.EucTwfreq[77][88] = 557
	detect.EucTwfreq[51][14] = 556
	detect.EucTwfreq[36][17] = 555
	detect.EucTwfreq[44][51] = 554
	detect.EucTwfreq[38][72] = 553
	detect.EucTwfreq[74][90] = 552
	detect.EucTwfreq[35][48] = 551
	detect.EucTwfreq[35][69] = 550
	detect.EucTwfreq[66][86] = 549
	detect.EucTwfreq[57][20] = 548
	detect.EucTwfreq[35][53] = 547
	detect.EucTwfreq[36][87] = 546
	detect.EucTwfreq[84][67] = 545
	detect.EucTwfreq[70][56] = 544
	detect.EucTwfreq[71][54] = 543
	detect.EucTwfreq[60][70] = 542
	detect.EucTwfreq[80][1] = 541
	detect.EucTwfreq[39][59] = 540
	detect.EucTwfreq[39][51] = 539
	detect.EucTwfreq[35][44] = 538
	detect.EucTwfreq[48][4] = 537
	detect.EucTwfreq[55][24] = 536
	detect.EucTwfreq[52][4] = 535
	detect.EucTwfreq[54][26] = 534
	detect.EucTwfreq[36][31] = 533
	detect.EucTwfreq[37][22] = 532
	detect.EucTwfreq[37][9] = 531
	detect.EucTwfreq[46][0] = 530
	detect.EucTwfreq[56][46] = 529
	detect.EucTwfreq[47][93] = 528
	detect.EucTwfreq[37][25] = 527
	detect.EucTwfreq[39][8] = 526
	detect.EucTwfreq[46][73] = 525
	detect.EucTwfreq[38][48] = 524
	detect.EucTwfreq[39][83] = 523
	detect.EucTwfreq[60][92] = 522
	detect.EucTwfreq[70][11] = 521
	detect.EucTwfreq[63][84] = 520
	detect.EucTwfreq[38][65] = 519
	detect.EucTwfreq[45][45] = 518
	detect.EucTwfreq[63][49] = 517
	detect.EucTwfreq[63][50] = 516
	detect.EucTwfreq[39][93] = 515
	detect.EucTwfreq[68][20] = 514
	detect.EucTwfreq[44][84] = 513
	detect.EucTwfreq[66][34] = 512
	detect.EucTwfreq[37][58] = 511
	detect.EucTwfreq[39][0] = 510
	detect.EucTwfreq[59][1] = 509
	detect.EucTwfreq[47][8] = 508
	detect.EucTwfreq[61][17] = 507
	detect.EucTwfreq[53][87] = 506
	detect.EucTwfreq[67][26] = 505
	detect.EucTwfreq[43][46] = 504
	detect.EucTwfreq[38][61] = 503
	detect.EucTwfreq[45][9] = 502
	detect.EucTwfreq[66][83] = 501
	detect.EucTwfreq[43][88] = 500
	detect.EucTwfreq[85][20] = 499
	detect.EucTwfreq[57][36] = 498
	detect.EucTwfreq[43][6] = 497
	detect.EucTwfreq[86][77] = 496
	detect.EucTwfreq[42][70] = 495
	detect.EucTwfreq[49][78] = 494
	detect.EucTwfreq[36][40] = 493
	detect.EucTwfreq[42][71] = 492
	detect.EucTwfreq[58][49] = 491
	detect.EucTwfreq[35][20] = 490
	detect.EucTwfreq[76][20] = 489
	detect.EucTwfreq[39][25] = 488
	detect.EucTwfreq[40][34] = 487
	detect.EucTwfreq[39][76] = 486
	detect.EucTwfreq[40][1] = 485
	detect.EucTwfreq[59][0] = 484
	detect.EucTwfreq[39][70] = 483
	detect.EucTwfreq[46][14] = 482
	detect.EucTwfreq[68][77] = 481
	detect.EucTwfreq[38][55] = 480
	detect.EucTwfreq[35][78] = 479
	detect.EucTwfreq[84][44] = 478
	detect.EucTwfreq[36][41] = 477
	detect.EucTwfreq[37][62] = 476
	detect.EucTwfreq[65][67] = 475
	detect.EucTwfreq[69][66] = 474
	detect.EucTwfreq[73][55] = 473
	detect.EucTwfreq[71][49] = 472
	detect.EucTwfreq[66][87] = 471
	detect.EucTwfreq[38][33] = 470
	detect.EucTwfreq[64][61] = 469
	detect.EucTwfreq[35][7] = 468
	detect.EucTwfreq[47][49] = 467
	detect.EucTwfreq[56][14] = 466
	detect.EucTwfreq[36][49] = 465
	detect.EucTwfreq[50][81] = 464
	detect.EucTwfreq[55][76] = 463
	detect.EucTwfreq[35][19] = 462
	detect.EucTwfreq[44][47] = 461
	detect.EucTwfreq[35][15] = 460
	detect.EucTwfreq[82][59] = 459
	detect.EucTwfreq[35][43] = 458
	detect.EucTwfreq[73][0] = 457
	detect.EucTwfreq[57][83] = 456
	detect.EucTwfreq[42][46] = 455
	detect.EucTwfreq[36][0] = 454
	detect.EucTwfreq[70][88] = 453
	detect.EucTwfreq[42][22] = 452
	detect.EucTwfreq[46][58] = 451
	detect.EucTwfreq[36][34] = 450
	detect.EucTwfreq[39][24] = 449
	detect.EucTwfreq[35][55] = 448
	detect.EucTwfreq[44][91] = 447
	detect.EucTwfreq[37][51] = 446
	detect.EucTwfreq[36][19] = 445
	detect.EucTwfreq[69][90] = 444
	detect.EucTwfreq[55][35] = 443
	detect.EucTwfreq[35][54] = 442
	detect.EucTwfreq[49][61] = 441
	detect.EucTwfreq[36][67] = 440
	detect.EucTwfreq[88][34] = 439
	detect.EucTwfreq[35][17] = 438
	detect.EucTwfreq[65][69] = 437
	detect.EucTwfreq[74][89] = 436
	detect.EucTwfreq[37][31] = 435
	detect.EucTwfreq[43][48] = 434
	detect.EucTwfreq[89][27] = 433
	detect.EucTwfreq[42][79] = 432
	detect.EucTwfreq[69][57] = 431
	detect.EucTwfreq[36][13] = 430
	detect.EucTwfreq[35][62] = 429
	detect.EucTwfreq[65][47] = 428
	detect.EucTwfreq[56][8] = 427
	detect.EucTwfreq[38][79] = 426
	detect.EucTwfreq[37][64] = 425
	detect.EucTwfreq[64][64] = 424
	detect.EucTwfreq[38][53] = 423
	detect.EucTwfreq[38][31] = 422
	detect.EucTwfreq[56][81] = 421
	detect.EucTwfreq[36][22] = 420
	detect.EucTwfreq[43][4] = 419
	detect.EucTwfreq[36][90] = 418
	detect.EucTwfreq[38][62] = 417
	detect.EucTwfreq[66][85] = 416
	detect.EucTwfreq[39][1] = 415
	detect.EucTwfreq[59][40] = 414
	detect.EucTwfreq[58][93] = 413
	detect.EucTwfreq[44][43] = 412
	detect.EucTwfreq[39][49] = 411
	detect.EucTwfreq[64][2] = 410
	detect.EucTwfreq[41][35] = 409
	detect.EucTwfreq[60][22] = 408
	detect.EucTwfreq[35][91] = 407
	detect.EucTwfreq[78][1] = 406
	detect.EucTwfreq[36][14] = 405
	detect.EucTwfreq[82][29] = 404
	detect.EucTwfreq[52][86] = 403
	detect.EucTwfreq[40][16] = 402
	detect.EucTwfreq[91][52] = 401
	detect.EucTwfreq[50][75] = 400
	detect.EucTwfreq[64][30] = 399
	detect.EucTwfreq[90][78] = 398
	detect.EucTwfreq[36][52] = 397
	detect.EucTwfreq[55][87] = 396
	detect.EucTwfreq[57][5] = 395
	detect.EucTwfreq[57][31] = 394
	detect.EucTwfreq[42][35] = 393
	detect.EucTwfreq[69][50] = 392
	detect.EucTwfreq[45][8] = 391
	detect.EucTwfreq[50][87] = 390
	detect.EucTwfreq[69][55] = 389
	detect.EucTwfreq[92][3] = 388
	detect.EucTwfreq[36][43] = 387
	detect.EucTwfreq[64][10] = 386
	detect.EucTwfreq[56][25] = 385
	detect.EucTwfreq[60][68] = 384
	detect.EucTwfreq[51][46] = 383
	detect.EucTwfreq[50][0] = 382
	detect.EucTwfreq[38][30] = 381
	detect.EucTwfreq[50][85] = 380
	detect.EucTwfreq[60][54] = 379
	detect.EucTwfreq[73][6] = 378
	detect.EucTwfreq[73][28] = 377
	detect.EucTwfreq[56][19] = 376
	detect.EucTwfreq[62][69] = 375
	detect.EucTwfreq[81][66] = 374
	detect.EucTwfreq[40][32] = 373
	detect.EucTwfreq[76][31] = 372
	detect.EucTwfreq[35][10] = 371
	detect.EucTwfreq[41][37] = 370
	detect.EucTwfreq[52][82] = 369
	detect.EucTwfreq[91][72] = 368
	detect.EucTwfreq[37][29] = 367
	detect.EucTwfreq[56][30] = 366
	detect.EucTwfreq[37][80] = 365
	detect.EucTwfreq[81][56] = 364
	detect.EucTwfreq[70][3] = 363
	detect.EucTwfreq[76][15] = 362
	detect.EucTwfreq[46][47] = 361
	detect.EucTwfreq[35][88] = 360
	detect.EucTwfreq[61][58] = 359
	detect.EucTwfreq[37][37] = 358
	detect.EucTwfreq[57][22] = 357
	detect.EucTwfreq[41][23] = 356
	detect.EucTwfreq[90][66] = 355
	detect.EucTwfreq[39][60] = 354
	detect.EucTwfreq[38][0] = 353
	detect.EucTwfreq[37][87] = 352
	detect.EucTwfreq[46][2] = 351
	detect.EucTwfreq[38][56] = 350
	detect.EucTwfreq[58][11] = 349
	detect.EucTwfreq[48][10] = 348
	detect.EucTwfreq[74][4] = 347
	detect.EucTwfreq[40][42] = 346
	detect.EucTwfreq[41][52] = 345
	detect.EucTwfreq[61][92] = 344
	detect.EucTwfreq[39][50] = 343
	detect.EucTwfreq[47][88] = 342
	detect.EucTwfreq[88][36] = 341
	detect.EucTwfreq[45][73] = 340
	detect.EucTwfreq[82][3] = 339
	detect.EucTwfreq[61][36] = 338
	detect.EucTwfreq[60][33] = 337
	detect.EucTwfreq[38][27] = 336
	detect.EucTwfreq[35][83] = 335
	detect.EucTwfreq[65][24] = 334
	detect.EucTwfreq[73][10] = 333
	detect.EucTwfreq[41][13] = 332
	detect.EucTwfreq[50][27] = 331
	detect.EucTwfreq[59][50] = 330
	detect.EucTwfreq[42][45] = 329
	detect.EucTwfreq[55][19] = 328
	detect.EucTwfreq[36][77] = 327
	detect.EucTwfreq[69][31] = 326
	detect.EucTwfreq[60][7] = 325
	detect.EucTwfreq[40][88] = 324
	detect.EucTwfreq[57][56] = 323
	detect.EucTwfreq[50][50] = 322
	detect.EucTwfreq[42][37] = 321
	detect.EucTwfreq[38][82] = 320
	detect.EucTwfreq[52][25] = 319
	detect.EucTwfreq[42][67] = 318
	detect.EucTwfreq[48][40] = 317
	detect.EucTwfreq[45][81] = 316
	detect.EucTwfreq[57][14] = 315
	detect.EucTwfreq[42][13] = 314
	detect.EucTwfreq[78][0] = 313
	detect.EucTwfreq[35][51] = 312
	detect.EucTwfreq[41][67] = 311
	detect.EucTwfreq[64][23] = 310
	detect.EucTwfreq[36][65] = 309
	detect.EucTwfreq[48][50] = 308
	detect.EucTwfreq[46][69] = 307
	detect.EucTwfreq[47][89] = 306
	detect.EucTwfreq[41][48] = 305
	detect.EucTwfreq[60][56] = 304
	detect.EucTwfreq[44][82] = 303
	detect.EucTwfreq[47][35] = 302
	detect.EucTwfreq[49][3] = 301
	detect.EucTwfreq[49][69] = 300
	detect.EucTwfreq[45][93] = 299
	detect.EucTwfreq[60][34] = 298
	detect.EucTwfreq[60][82] = 297
	detect.EucTwfreq[61][61] = 296
	detect.EucTwfreq[86][42] = 295
	detect.EucTwfreq[89][60] = 294
	detect.EucTwfreq[48][31] = 293
	detect.EucTwfreq[35][75] = 292
	detect.EucTwfreq[91][39] = 291
	detect.EucTwfreq[53][19] = 290
	detect.EucTwfreq[39][72] = 289
	detect.EucTwfreq[69][59] = 288
	detect.EucTwfreq[41][7] = 287
	detect.EucTwfreq[54][13] = 286
	detect.EucTwfreq[43][28] = 285
	detect.EucTwfreq[36][6] = 284
	detect.EucTwfreq[45][75] = 283
	detect.EucTwfreq[36][61] = 282
	detect.EucTwfreq[38][21] = 281
	detect.EucTwfreq[45][14] = 280
	detect.EucTwfreq[61][43] = 279
	detect.EucTwfreq[36][63] = 278
	detect.EucTwfreq[43][30] = 277
	detect.EucTwfreq[46][51] = 276
	detect.EucTwfreq[68][87] = 275
	detect.EucTwfreq[39][26] = 274
	detect.EucTwfreq[46][76] = 273
	detect.EucTwfreq[36][15] = 272
	detect.EucTwfreq[35][40] = 271
	detect.EucTwfreq[79][60] = 270
	detect.EucTwfreq[46][7] = 269
	detect.EucTwfreq[65][72] = 268
	detect.EucTwfreq[69][88] = 267
	detect.EucTwfreq[47][18] = 266
	detect.EucTwfreq[37][0] = 265
	detect.EucTwfreq[37][49] = 264
	detect.EucTwfreq[67][37] = 263
	detect.EucTwfreq[36][91] = 262
	detect.EucTwfreq[75][48] = 261
	detect.EucTwfreq[75][63] = 260
	detect.EucTwfreq[83][87] = 259
	detect.EucTwfreq[37][44] = 258
	detect.EucTwfreq[73][54] = 257
	detect.EucTwfreq[51][61] = 256
	detect.EucTwfreq[46][57] = 255
	detect.EucTwfreq[55][21] = 254
	detect.EucTwfreq[39][66] = 253
	detect.EucTwfreq[47][11] = 252
	detect.EucTwfreq[52][8] = 251
	detect.EucTwfreq[82][81] = 250
	detect.EucTwfreq[36][57] = 249
	detect.EucTwfreq[38][54] = 248
	detect.EucTwfreq[43][81] = 247
	detect.EucTwfreq[37][42] = 246
	detect.EucTwfreq[40][18] = 245
	detect.EucTwfreq[80][90] = 244
	detect.EucTwfreq[37][84] = 243
	detect.EucTwfreq[57][15] = 242
	detect.EucTwfreq[38][87] = 241
	detect.EucTwfreq[37][32] = 240
	detect.EucTwfreq[53][53] = 239
	detect.EucTwfreq[89][29] = 238
	detect.EucTwfreq[81][53] = 237
	detect.EucTwfreq[75][3] = 236
	detect.EucTwfreq[83][73] = 235
	detect.EucTwfreq[66][13] = 234
	detect.EucTwfreq[48][7] = 233
	detect.EucTwfreq[46][35] = 232
	detect.EucTwfreq[35][86] = 231
	detect.EucTwfreq[37][20] = 230
	detect.EucTwfreq[46][80] = 229
	detect.EucTwfreq[38][24] = 228
	detect.EucTwfreq[41][68] = 227
	detect.EucTwfreq[42][21] = 226
	detect.EucTwfreq[43][32] = 225
	detect.EucTwfreq[38][20] = 224
	detect.EucTwfreq[37][59] = 223
	detect.EucTwfreq[41][77] = 222
	detect.EucTwfreq[59][57] = 221
	detect.EucTwfreq[68][59] = 220
	detect.EucTwfreq[39][43] = 219
	detect.EucTwfreq[54][39] = 218
	detect.EucTwfreq[48][28] = 217
	detect.EucTwfreq[54][28] = 216
	detect.EucTwfreq[41][44] = 215
	detect.EucTwfreq[51][64] = 214
	detect.EucTwfreq[47][72] = 213
	detect.EucTwfreq[62][67] = 212
	detect.EucTwfreq[42][43] = 211
	detect.EucTwfreq[61][38] = 210
	detect.EucTwfreq[76][25] = 209
	detect.EucTwfreq[48][91] = 208
	detect.EucTwfreq[36][36] = 207
	detect.EucTwfreq[80][32] = 206
	detect.EucTwfreq[81][40] = 205
	detect.EucTwfreq[37][5] = 204
	detect.EucTwfreq[74][69] = 203
	detect.EucTwfreq[36][82] = 202
	detect.EucTwfreq[46][59] = 201
	/*
	 * EucTwfreq[38][32] = 200; EucTwfreq[74][2] = 199; EucTwfreq[53][31]
	 * = 198; EucTwfreq[35][38] = 197; EucTwfreq[46][62] = 196;
	 * EucTwfreq[77][31] = 195; EucTwfreq[55][74] = 194; EucTwfreq[66][6]
	 * = 193; EucTwfreq[56][21] = 192; EucTwfreq[54][78] = 191;
	 * EucTwfreq[43][51] = 190; EucTwfreq[64][93] = 189; EucTwfreq[92][7]
	 * = 188; EucTwfreq[83][89] = 187; EucTwfreq[69][9] = 186;
	 * EucTwfreq[45][4] = 185; EucTwfreq[53][9] = 184; EucTwfreq[43][2] =
	 * 183; EucTwfreq[35][11] = 182; EucTwfreq[51][25] = 181;
	 * EucTwfreq[52][71] = 180; EucTwfreq[81][67] = 179;
	 * EucTwfreq[37][33] = 178; EucTwfreq[38][57] = 177;
	 * EucTwfreq[39][77] = 176; EucTwfreq[40][26] = 175;
	 * EucTwfreq[37][21] = 174; EucTwfreq[81][70] = 173;
	 * EucTwfreq[56][80] = 172; EucTwfreq[65][14] = 171;
	 * EucTwfreq[62][47] = 170; EucTwfreq[56][54] = 169;
	 * EucTwfreq[45][17] = 168; EucTwfreq[52][52] = 167;
	 * EucTwfreq[74][30] = 166; EucTwfreq[60][57] = 165;
	 * EucTwfreq[41][15] = 164; EucTwfreq[47][69] = 163;
	 * EucTwfreq[61][11] = 162; EucTwfreq[72][25] = 161;
	 * EucTwfreq[82][56] = 160; EucTwfreq[76][92] = 159;
	 * EucTwfreq[51][22] = 158; EucTwfreq[55][69] = 157;
	 * EucTwfreq[49][43] = 156; EucTwfreq[69][49] = 155;
	 * EucTwfreq[88][42] = 154; EucTwfreq[84][41] = 153;
	 * EucTwfreq[79][33] = 152; EucTwfreq[47][17] = 151;
	 * EucTwfreq[52][88] = 150; EucTwfreq[63][74] = 149;
	 * EucTwfreq[50][32] = 148; EucTwfreq[65][10] = 147; EucTwfreq[57][6]
	 * = 146; EucTwfreq[52][23] = 145; EucTwfreq[36][70] = 144;
	 * EucTwfreq[65][55] = 143; EucTwfreq[35][27] = 142;
	 * EucTwfreq[57][63] = 141; EucTwfreq[39][92] = 140;
	 * EucTwfreq[79][75] = 139; EucTwfreq[36][30] = 138;
	 * EucTwfreq[53][60] = 137; EucTwfreq[55][43] = 136;
	 * EucTwfreq[71][22] = 135; EucTwfreq[43][16] = 134;
	 * EucTwfreq[65][21] = 133; EucTwfreq[84][51] = 132;
	 * EucTwfreq[43][64] = 131; EucTwfreq[87][91] = 130;
	 * EucTwfreq[47][45] = 129; EucTwfreq[65][29] = 128;
	 * EucTwfreq[88][16] = 127; EucTwfreq[50][5] = 126; EucTwfreq[47][33]
	 * = 125; EucTwfreq[46][27] = 124; EucTwfreq[85][2] = 123;
	 * EucTwfreq[43][77] = 122; EucTwfreq[70][9] = 121; EucTwfreq[41][54]
	 * = 120; EucTwfreq[56][12] = 119; EucTwfreq[90][65] = 118;
	 * EucTwfreq[91][50] = 117; EucTwfreq[48][41] = 116;
	 * EucTwfreq[35][89] = 115; EucTwfreq[90][83] = 114;
	 * EucTwfreq[44][40] = 113; EucTwfreq[50][88] = 112;
	 * EucTwfreq[72][39] = 111; EucTwfreq[45][3] = 110; EucTwfreq[71][33]
	 * = 109; EucTwfreq[39][12] = 108; EucTwfreq[59][24] = 107;
	 * EucTwfreq[60][62] = 106; EucTwfreq[44][33] = 105;
	 * EucTwfreq[53][70] = 104; EucTwfreq[77][90] = 103;
	 * EucTwfreq[50][58] = 102; EucTwfreq[54][1] = 101; EucTwfreq[73][19]
	 * = 100; EucTwfreq[37][3] = 99; EucTwfreq[49][91] = 98;
	 * EucTwfreq[88][43] = 97; EucTwfreq[36][78] = 96; EucTwfreq[44][20]
	 * = 95; EucTwfreq[64][15] = 94; EucTwfreq[72][28] = 93;
	 * EucTwfreq[70][13] = 92; EucTwfreq[65][83] = 91; EucTwfreq[58][68]
	 * = 90; EucTwfreq[59][32] = 89; EucTwfreq[39][13] = 88;
	 * EucTwfreq[55][64] = 87; EucTwfreq[56][59] = 86; EucTwfreq[39][17]
	 * = 85; EucTwfreq[55][84] = 84; EucTwfreq[77][85] = 83;
	 * EucTwfreq[60][19] = 82; EucTwfreq[62][82] = 81; EucTwfreq[78][16]
	 * = 80; EucTwfreq[66][8] = 79; EucTwfreq[39][42] = 78;
	 * EucTwfreq[61][24] = 77; EucTwfreq[57][67] = 76; EucTwfreq[38][83]
	 * = 75; EucTwfreq[36][53] = 74; EucTwfreq[67][76] = 73;
	 * EucTwfreq[37][91] = 72; EucTwfreq[44][26] = 71; EucTwfreq[72][86]
	 * = 70; EucTwfreq[44][87] = 69; EucTwfreq[45][50] = 68;
	 * EucTwfreq[58][4] = 67; EucTwfreq[86][65] = 66; EucTwfreq[45][56] =
	 * 65; EucTwfreq[79][49] = 64; EucTwfreq[35][3] = 63;
	 * EucTwfreq[48][83] = 62; EucTwfreq[71][21] = 61; EucTwfreq[77][93]
	 * = 60; EucTwfreq[87][92] = 59; EucTwfreq[38][35] = 58;
	 * EucTwfreq[66][17] = 57; EucTwfreq[37][66] = 56; EucTwfreq[51][42]
	 * = 55; EucTwfreq[57][73] = 54; EucTwfreq[51][54] = 53;
	 * EucTwfreq[75][64] = 52; EucTwfreq[35][5] = 51; EucTwfreq[49][40] =
	 * 50; EucTwfreq[58][35] = 49; EucTwfreq[67][88] = 48;
	 * EucTwfreq[60][51] = 47; EucTwfreq[36][92] = 46; EucTwfreq[44][41]
	 * = 45; EucTwfreq[58][29] = 44; EucTwfreq[43][62] = 43;
	 * EucTwfreq[56][23] = 42; EucTwfreq[67][44] = 41; EucTwfreq[52][91]
	 * = 40; EucTwfreq[42][81] = 39; EucTwfreq[64][25] = 38;
	 * EucTwfreq[35][36] = 37; EucTwfreq[47][73] = 36; EucTwfreq[36][1] =
	 * 35; EucTwfreq[65][84] = 34; EucTwfreq[73][1] = 33;
	 * EucTwfreq[79][66] = 32; EucTwfreq[69][14] = 31; EucTwfreq[65][28]
	 * = 30; EucTwfreq[60][93] = 29; EucTwfreq[72][79] = 28;
	 * EucTwfreq[48][0] = 27; EucTwfreq[73][43] = 26; EucTwfreq[66][47] =
	 * 25; EucTwfreq[41][18] = 24; EucTwfreq[51][10] = 23;
	 * EucTwfreq[59][7] = 22; EucTwfreq[53][27] = 21; EucTwfreq[86][67] =
	 * 20; EucTwfreq[49][87] = 19; EucTwfreq[52][28] = 18;
	 * EucTwfreq[52][12] = 17; EucTwfreq[42][30] = 16; EucTwfreq[65][35]
	 * = 15; EucTwfreq[46][64] = 14; EucTwfreq[71][7] = 13;
	 * EucTwfreq[56][57] = 12; EucTwfreq[56][31] = 11; EucTwfreq[41][31]
	 * = 10; EucTwfreq[48][59] = 9; EucTwfreq[63][92] = 8;
	 * EucTwfreq[62][57] = 7; EucTwfreq[65][87] = 6; EucTwfreq[70][10] =
	 * 5; EucTwfreq[52][40] = 4; EucTwfreq[40][22] = 3; EucTwfreq[65][91]
	 * = 2; EucTwfreq[50][25] = 1; EucTwfreq[35][84] = 0;
	 */
	detect.GBKFreq[52][132] = 600
	detect.GBKFreq[73][135] = 599
	detect.GBKFreq[49][123] = 598
	detect.GBKFreq[77][146] = 597
	detect.GBKFreq[81][123] = 596
	detect.GBKFreq[82][144] = 595
	detect.GBKFreq[51][179] = 594
	detect.GBKFreq[83][154] = 593
	detect.GBKFreq[71][139] = 592
	detect.GBKFreq[64][139] = 591
	detect.GBKFreq[85][144] = 590
	detect.GBKFreq[52][125] = 589
	detect.GBKFreq[88][25] = 588
	detect.GBKFreq[81][106] = 587
	detect.GBKFreq[81][148] = 586
	detect.GBKFreq[62][137] = 585
	detect.GBKFreq[94][0] = 584
	detect.GBKFreq[1][64] = 583
	detect.GBKFreq[67][163] = 582
	detect.GBKFreq[20][190] = 581
	detect.GBKFreq[57][131] = 580
	detect.GBKFreq[29][169] = 579
	detect.GBKFreq[72][143] = 578
	detect.GBKFreq[0][173] = 577
	detect.GBKFreq[11][23] = 576
	detect.GBKFreq[61][141] = 575
	detect.GBKFreq[60][123] = 574
	detect.GBKFreq[81][114] = 573
	detect.GBKFreq[82][131] = 572
	detect.GBKFreq[67][156] = 571
	detect.GBKFreq[71][167] = 570
	detect.GBKFreq[20][50] = 569
	detect.GBKFreq[77][132] = 568
	detect.GBKFreq[84][38] = 567
	detect.GBKFreq[26][29] = 566
	detect.GBKFreq[74][187] = 565
	detect.GBKFreq[62][116] = 564
	detect.GBKFreq[67][135] = 563
	detect.GBKFreq[5][86] = 562
	detect.GBKFreq[72][186] = 561
	detect.GBKFreq[75][161] = 560
	detect.GBKFreq[78][130] = 559
	detect.GBKFreq[94][30] = 558
	detect.GBKFreq[84][72] = 557
	detect.GBKFreq[1][67] = 556
	detect.GBKFreq[75][172] = 555
	detect.GBKFreq[74][185] = 554
	detect.GBKFreq[53][160] = 553
	detect.GBKFreq[123][14] = 552
	detect.GBKFreq[79][97] = 551
	detect.GBKFreq[85][110] = 550
	detect.GBKFreq[78][171] = 549
	detect.GBKFreq[52][131] = 548
	detect.GBKFreq[56][100] = 547
	detect.GBKFreq[50][182] = 546
	detect.GBKFreq[94][64] = 545
	detect.GBKFreq[106][74] = 544
	detect.GBKFreq[11][102] = 543
	detect.GBKFreq[53][124] = 542
	detect.GBKFreq[24][3] = 541
	detect.GBKFreq[86][148] = 540
	detect.GBKFreq[53][184] = 539
	detect.GBKFreq[86][147] = 538
	detect.GBKFreq[96][161] = 537
	detect.GBKFreq[82][77] = 536
	detect.GBKFreq[59][146] = 535
	detect.GBKFreq[84][126] = 534
	detect.GBKFreq[79][132] = 533
	detect.GBKFreq[85][123] = 532
	detect.GBKFreq[71][101] = 531
	detect.GBKFreq[85][106] = 530
	detect.GBKFreq[6][184] = 529
	detect.GBKFreq[57][156] = 528
	detect.GBKFreq[75][104] = 527
	detect.GBKFreq[50][137] = 526
	detect.GBKFreq[79][133] = 525
	detect.GBKFreq[76][108] = 524
	detect.GBKFreq[57][142] = 523
	detect.GBKFreq[84][130] = 522
	detect.GBKFreq[52][128] = 521
	detect.GBKFreq[47][44] = 520
	detect.GBKFreq[52][152] = 519
	detect.GBKFreq[54][104] = 518
	detect.GBKFreq[30][47] = 517
	detect.GBKFreq[71][123] = 516
	detect.GBKFreq[52][107] = 515
	detect.GBKFreq[45][84] = 514
	detect.GBKFreq[107][118] = 513
	detect.GBKFreq[5][161] = 512
	detect.GBKFreq[48][126] = 511
	detect.GBKFreq[67][170] = 510
	detect.GBKFreq[43][6] = 509
	detect.GBKFreq[70][112] = 508
	detect.GBKFreq[86][174] = 507
	detect.GBKFreq[84][166] = 506
	detect.GBKFreq[79][130] = 505
	detect.GBKFreq[57][141] = 504
	detect.GBKFreq[81][178] = 503
	detect.GBKFreq[56][187] = 502
	detect.GBKFreq[81][162] = 501
	detect.GBKFreq[53][104] = 500
	detect.GBKFreq[123][35] = 499
	detect.GBKFreq[70][169] = 498
	detect.GBKFreq[69][164] = 497
	detect.GBKFreq[109][61] = 496
	detect.GBKFreq[73][130] = 495
	detect.GBKFreq[62][134] = 494
	detect.GBKFreq[54][125] = 493
	detect.GBKFreq[79][105] = 492
	detect.GBKFreq[70][165] = 491
	detect.GBKFreq[71][189] = 490
	detect.GBKFreq[23][147] = 489
	detect.GBKFreq[51][139] = 488
	detect.GBKFreq[47][137] = 487
	detect.GBKFreq[77][123] = 486
	detect.GBKFreq[86][183] = 485
	detect.GBKFreq[63][173] = 484
	detect.GBKFreq[79][144] = 483
	detect.GBKFreq[84][159] = 482
	detect.GBKFreq[60][91] = 481
	detect.GBKFreq[66][187] = 480
	detect.GBKFreq[73][114] = 479
	detect.GBKFreq[85][56] = 478
	detect.GBKFreq[71][149] = 477
	detect.GBKFreq[84][189] = 476
	detect.GBKFreq[104][31] = 475
	detect.GBKFreq[83][82] = 474
	detect.GBKFreq[68][35] = 473
	detect.GBKFreq[11][77] = 472
	detect.GBKFreq[15][155] = 471
	detect.GBKFreq[83][153] = 470
	detect.GBKFreq[71][1] = 469
	detect.GBKFreq[53][190] = 468
	detect.GBKFreq[50][135] = 467
	detect.GBKFreq[3][147] = 466
	detect.GBKFreq[48][136] = 465
	detect.GBKFreq[66][166] = 464
	detect.GBKFreq[55][159] = 463
	detect.GBKFreq[82][150] = 462
	detect.GBKFreq[58][178] = 461
	detect.GBKFreq[64][102] = 460
	detect.GBKFreq[16][106] = 459
	detect.GBKFreq[68][110] = 458
	detect.GBKFreq[54][14] = 457
	detect.GBKFreq[60][140] = 456
	detect.GBKFreq[91][71] = 455
	detect.GBKFreq[54][150] = 454
	detect.GBKFreq[78][177] = 453
	detect.GBKFreq[78][117] = 452
	detect.GBKFreq[104][12] = 451
	detect.GBKFreq[73][150] = 450
	detect.GBKFreq[51][142] = 449
	detect.GBKFreq[81][145] = 448
	detect.GBKFreq[66][183] = 447
	detect.GBKFreq[51][178] = 446
	detect.GBKFreq[75][107] = 445
	detect.GBKFreq[65][119] = 444
	detect.GBKFreq[69][176] = 443
	detect.GBKFreq[59][122] = 442
	detect.GBKFreq[78][160] = 441
	detect.GBKFreq[85][183] = 440
	detect.GBKFreq[105][16] = 439
	detect.GBKFreq[73][110] = 438
	detect.GBKFreq[104][39] = 437
	detect.GBKFreq[119][16] = 436
	detect.GBKFreq[76][162] = 435
	detect.GBKFreq[67][152] = 434
	detect.GBKFreq[82][24] = 433
	detect.GBKFreq[73][121] = 432
	detect.GBKFreq[83][83] = 431
	detect.GBKFreq[82][145] = 430
	detect.GBKFreq[49][133] = 429
	detect.GBKFreq[94][13] = 428
	detect.GBKFreq[58][139] = 427
	detect.GBKFreq[74][189] = 426
	detect.GBKFreq[66][177] = 425
	detect.GBKFreq[85][184] = 424
	detect.GBKFreq[55][183] = 423
	detect.GBKFreq[71][107] = 422
	detect.GBKFreq[11][98] = 421
	detect.GBKFreq[72][153] = 420
	detect.GBKFreq[2][137] = 419
	detect.GBKFreq[59][147] = 418
	detect.GBKFreq[58][152] = 417
	detect.GBKFreq[55][144] = 416
	detect.GBKFreq[73][125] = 415
	detect.GBKFreq[52][154] = 414
	detect.GBKFreq[70][178] = 413
	detect.GBKFreq[79][148] = 412
	detect.GBKFreq[63][143] = 411
	detect.GBKFreq[50][140] = 410
	detect.GBKFreq[47][145] = 409
	detect.GBKFreq[48][123] = 408
	detect.GBKFreq[56][107] = 407
	detect.GBKFreq[84][83] = 406
	detect.GBKFreq[59][112] = 405
	detect.GBKFreq[124][72] = 404
	detect.GBKFreq[79][99] = 403
	detect.GBKFreq[3][37] = 402
	detect.GBKFreq[114][55] = 401
	detect.GBKFreq[85][152] = 400
	detect.GBKFreq[60][47] = 399
	detect.GBKFreq[65][96] = 398
	detect.GBKFreq[74][110] = 397
	detect.GBKFreq[86][182] = 396
	detect.GBKFreq[50][99] = 395
	detect.GBKFreq[67][186] = 394
	detect.GBKFreq[81][74] = 393
	detect.GBKFreq[80][37] = 392
	detect.GBKFreq[21][60] = 391
	detect.GBKFreq[110][12] = 390
	detect.GBKFreq[60][162] = 389
	detect.GBKFreq[29][115] = 388
	detect.GBKFreq[83][130] = 387
	detect.GBKFreq[52][136] = 386
	detect.GBKFreq[63][114] = 385
	detect.GBKFreq[49][127] = 384
	detect.GBKFreq[83][109] = 383
	detect.GBKFreq[66][128] = 382
	detect.GBKFreq[78][136] = 381
	detect.GBKFreq[81][180] = 380
	detect.GBKFreq[76][104] = 379
	detect.GBKFreq[56][156] = 378
	detect.GBKFreq[61][23] = 377
	detect.GBKFreq[4][30] = 376
	detect.GBKFreq[69][154] = 375
	detect.GBKFreq[100][37] = 374
	detect.GBKFreq[54][177] = 373
	detect.GBKFreq[23][119] = 372
	detect.GBKFreq[71][171] = 371
	detect.GBKFreq[84][146] = 370
	detect.GBKFreq[20][184] = 369
	detect.GBKFreq[86][76] = 368
	detect.GBKFreq[74][132] = 367
	detect.GBKFreq[47][97] = 366
	detect.GBKFreq[82][137] = 365
	detect.GBKFreq[94][56] = 364
	detect.GBKFreq[92][30] = 363
	detect.GBKFreq[19][117] = 362
	detect.GBKFreq[48][173] = 361
	detect.GBKFreq[2][136] = 360
	detect.GBKFreq[7][182] = 359
	detect.GBKFreq[74][188] = 358
	detect.GBKFreq[14][132] = 357
	detect.GBKFreq[62][172] = 356
	detect.GBKFreq[25][39] = 355
	detect.GBKFreq[85][129] = 354
	detect.GBKFreq[64][98] = 353
	detect.GBKFreq[67][127] = 352
	detect.GBKFreq[72][167] = 351
	detect.GBKFreq[57][143] = 350
	detect.GBKFreq[76][187] = 349
	detect.GBKFreq[83][181] = 348
	detect.GBKFreq[84][10] = 347
	detect.GBKFreq[55][166] = 346
	detect.GBKFreq[55][188] = 345
	detect.GBKFreq[13][151] = 344
	detect.GBKFreq[62][124] = 343
	detect.GBKFreq[53][136] = 342
	detect.GBKFreq[106][57] = 341
	detect.GBKFreq[47][166] = 340
	detect.GBKFreq[109][30] = 339
	detect.GBKFreq[78][114] = 338
	detect.GBKFreq[83][19] = 337
	detect.GBKFreq[56][162] = 336
	detect.GBKFreq[60][177] = 335
	detect.GBKFreq[88][9] = 334
	detect.GBKFreq[74][163] = 333
	detect.GBKFreq[52][156] = 332
	detect.GBKFreq[71][180] = 331
	detect.GBKFreq[60][57] = 330
	detect.GBKFreq[72][173] = 329
	detect.GBKFreq[82][91] = 328
	detect.GBKFreq[51][186] = 327
	detect.GBKFreq[75][86] = 326
	detect.GBKFreq[75][78] = 325
	detect.GBKFreq[76][170] = 324
	detect.GBKFreq[60][147] = 323
	detect.GBKFreq[82][75] = 322
	detect.GBKFreq[80][148] = 321
	detect.GBKFreq[86][150] = 320
	detect.GBKFreq[13][95] = 319
	detect.GBKFreq[0][11] = 318
	detect.GBKFreq[84][190] = 317
	detect.GBKFreq[76][166] = 316
	detect.GBKFreq[14][72] = 315
	detect.GBKFreq[67][144] = 314
	detect.GBKFreq[84][44] = 313
	detect.GBKFreq[72][125] = 312
	detect.GBKFreq[66][127] = 311
	detect.GBKFreq[60][25] = 310
	detect.GBKFreq[70][146] = 309
	detect.GBKFreq[79][135] = 308
	detect.GBKFreq[54][135] = 307
	detect.GBKFreq[60][104] = 306
	detect.GBKFreq[55][132] = 305
	detect.GBKFreq[94][2] = 304
	detect.GBKFreq[54][133] = 303
	detect.GBKFreq[56][190] = 302
	detect.GBKFreq[58][174] = 301
	detect.GBKFreq[80][144] = 300
	detect.GBKFreq[85][113] = 299
	/*
	 * GBKFreq[83][15] = 298; GBKFreq[105][80] = 297; GBKFreq[7][179] = 296;
	 * GBKFreq[93][4] = 295; GBKFreq[123][40] = 294; GBKFreq[85][120] = 293;
	 * GBKFreq[77][165] = 292; GBKFreq[86][67] = 291; GBKFreq[25][162] =
	 * 290; GBKFreq[77][183] = 289; GBKFreq[83][71] = 288; GBKFreq[78][99] =
	 * 287; GBKFreq[72][177] = 286; GBKFreq[71][97] = 285; GBKFreq[58][111]
	 * = 284; GBKFreq[77][175] = 283; GBKFreq[76][181] = 282;
	 * GBKFreq[71][142] = 281; GBKFreq[64][150] = 280; GBKFreq[5][142] =
	 * 279; GBKFreq[73][128] = 278; GBKFreq[73][156] = 277; GBKFreq[60][188]
	 * = 276; GBKFreq[64][56] = 275; GBKFreq[74][128] = 274;
	 * GBKFreq[48][163] = 273; GBKFreq[54][116] = 272; GBKFreq[73][127] =
	 * 271; GBKFreq[16][176] = 270; GBKFreq[62][149] = 269; GBKFreq[105][96]
	 * = 268; GBKFreq[55][186] = 267; GBKFreq[4][51] = 266; GBKFreq[48][113]
	 * = 265; GBKFreq[48][152] = 264; GBKFreq[23][9] = 263; GBKFreq[56][102]
	 * = 262; GBKFreq[11][81] = 261; GBKFreq[82][112] = 260; GBKFreq[65][85]
	 * = 259; GBKFreq[69][125] = 258; GBKFreq[68][31] = 257; GBKFreq[5][20]
	 * = 256; GBKFreq[60][176] = 255; GBKFreq[82][81] = 254;
	 * GBKFreq[72][107] = 253; GBKFreq[3][52] = 252; GBKFreq[71][157] = 251;
	 * GBKFreq[24][46] = 250; GBKFreq[69][108] = 249; GBKFreq[78][178] =
	 * 248; GBKFreq[9][69] = 247; GBKFreq[73][144] = 246; GBKFreq[63][187] =
	 * 245; GBKFreq[68][36] = 244; GBKFreq[47][151] = 243; GBKFreq[14][74] =
	 * 242; GBKFreq[47][114] = 241; GBKFreq[80][171] = 240; GBKFreq[75][152]
	 * = 239; GBKFreq[86][40] = 238; GBKFreq[93][43] = 237; GBKFreq[2][50] =
	 * 236; GBKFreq[62][66] = 235; GBKFreq[1][183] = 234; GBKFreq[74][124] =
	 * 233; GBKFreq[58][104] = 232; GBKFreq[83][106] = 231; GBKFreq[60][144]
	 * = 230; GBKFreq[48][99] = 229; GBKFreq[54][157] = 228;
	 * GBKFreq[70][179] = 227; GBKFreq[61][127] = 226; GBKFreq[57][135] =
	 * 225; GBKFreq[59][190] = 224; GBKFreq[77][116] = 223; GBKFreq[26][17]
	 * = 222; GBKFreq[60][13] = 221; GBKFreq[71][38] = 220; GBKFreq[85][177]
	 * = 219; GBKFreq[59][73] = 218; GBKFreq[50][150] = 217;
	 * GBKFreq[79][102] = 216; GBKFreq[76][118] = 215; GBKFreq[67][132] =
	 * 214; GBKFreq[73][146] = 213; GBKFreq[83][184] = 212; GBKFreq[86][159]
	 * = 211; GBKFreq[95][120] = 210; GBKFreq[23][139] = 209;
	 * GBKFreq[64][183] = 208; GBKFreq[85][103] = 207; GBKFreq[41][90] =
	 * 206; GBKFreq[87][72] = 205; GBKFreq[62][104] = 204; GBKFreq[79][168]
	 * = 203; GBKFreq[79][150] = 202; GBKFreq[104][20] = 201;
	 * GBKFreq[56][114] = 200; GBKFreq[84][26] = 199; GBKFreq[57][99] = 198;
	 * GBKFreq[62][154] = 197; GBKFreq[47][98] = 196; GBKFreq[61][64] = 195;
	 * GBKFreq[112][18] = 194; GBKFreq[123][19] = 193; GBKFreq[4][98] = 192;
	 * GBKFreq[47][163] = 191; GBKFreq[66][188] = 190; GBKFreq[81][85] =
	 * 189; GBKFreq[82][30] = 188; GBKFreq[65][83] = 187; GBKFreq[67][24] =
	 * 186; GBKFreq[68][179] = 185; GBKFreq[55][177] = 184; GBKFreq[2][122]
	 * = 183; GBKFreq[47][139] = 182; GBKFreq[79][158] = 181;
	 * GBKFreq[64][143] = 180; GBKFreq[100][24] = 179; GBKFreq[73][103] =
	 * 178; GBKFreq[50][148] = 177; GBKFreq[86][97] = 176; GBKFreq[59][116]
	 * = 175; GBKFreq[64][173] = 174; GBKFreq[99][91] = 173; GBKFreq[11][99]
	 * = 172; GBKFreq[78][179] = 171; GBKFreq[18][17] = 170;
	 * GBKFreq[58][185] = 169; GBKFreq[47][165] = 168; GBKFreq[67][131] =
	 * 167; GBKFreq[94][40] = 166; GBKFreq[74][153] = 165; GBKFreq[79][142]
	 * = 164; GBKFreq[57][98] = 163; GBKFreq[1][164] = 162; GBKFreq[55][168]
	 * = 161; GBKFreq[13][141] = 160; GBKFreq[51][31] = 159;
	 * GBKFreq[57][178] = 158; GBKFreq[50][189] = 157; GBKFreq[60][167] =
	 * 156; GBKFreq[80][34] = 155; GBKFreq[109][80] = 154; GBKFreq[85][54] =
	 * 153; GBKFreq[69][183] = 152; GBKFreq[67][143] = 151; GBKFreq[47][120]
	 * = 150; GBKFreq[45][75] = 149; GBKFreq[82][98] = 148; GBKFreq[83][22]
	 * = 147; GBKFreq[13][103] = 146; GBKFreq[49][174] = 145;
	 * GBKFreq[57][181] = 144; GBKFreq[64][127] = 143; GBKFreq[61][131] =
	 * 142; GBKFreq[52][180] = 141; GBKFreq[74][134] = 140; GBKFreq[84][187]
	 * = 139; GBKFreq[81][189] = 138; GBKFreq[47][160] = 137;
	 * GBKFreq[66][148] = 136; GBKFreq[7][4] = 135; GBKFreq[85][134] = 134;
	 * GBKFreq[88][13] = 133; GBKFreq[88][80] = 132; GBKFreq[69][166] = 131;
	 * GBKFreq[86][18] = 130; GBKFreq[79][141] = 129; GBKFreq[50][108] =
	 * 128; GBKFreq[94][69] = 127; GBKFreq[81][110] = 126; GBKFreq[69][119]
	 * = 125; GBKFreq[72][161] = 124; GBKFreq[106][45] = 123;
	 * GBKFreq[73][124] = 122; GBKFreq[94][28] = 121; GBKFreq[63][174] =
	 * 120; GBKFreq[3][149] = 119; GBKFreq[24][160] = 118; GBKFreq[113][94]
	 * = 117; GBKFreq[56][138] = 116; GBKFreq[64][185] = 115;
	 * GBKFreq[86][56] = 114; GBKFreq[56][150] = 113; GBKFreq[110][55] =
	 * 112; GBKFreq[28][13] = 111; GBKFreq[54][190] = 110; GBKFreq[8][180] =
	 * 109; GBKFreq[73][149] = 108; GBKFreq[80][155] = 107; GBKFreq[83][172]
	 * = 106; GBKFreq[67][174] = 105; GBKFreq[64][180] = 104;
	 * GBKFreq[84][46] = 103; GBKFreq[91][74] = 102; GBKFreq[69][134] = 101;
	 * GBKFreq[61][107] = 100; GBKFreq[47][171] = 99; GBKFreq[59][51] = 98;
	 * GBKFreq[109][74] = 97; GBKFreq[64][174] = 96; GBKFreq[52][151] = 95;
	 * GBKFreq[51][176] = 94; GBKFreq[80][157] = 93; GBKFreq[94][31] = 92;
	 * GBKFreq[79][155] = 91; GBKFreq[72][174] = 90; GBKFreq[69][113] = 89;
	 * GBKFreq[83][167] = 88; GBKFreq[83][122] = 87; GBKFreq[8][178] = 86;
	 * GBKFreq[70][186] = 85; GBKFreq[59][153] = 84; GBKFreq[84][68] = 83;
	 * GBKFreq[79][39] = 82; GBKFreq[47][180] = 81; GBKFreq[88][53] = 80;
	 * GBKFreq[57][154] = 79; GBKFreq[47][153] = 78; GBKFreq[3][153] = 77;
	 * GBKFreq[76][134] = 76; GBKFreq[51][166] = 75; GBKFreq[58][176] = 74;
	 * GBKFreq[27][138] = 73; GBKFreq[73][126] = 72; GBKFreq[76][185] = 71;
	 * GBKFreq[52][186] = 70; GBKFreq[81][151] = 69; GBKFreq[26][50] = 68;
	 * GBKFreq[76][173] = 67; GBKFreq[106][56] = 66; GBKFreq[85][142] = 65;
	 * GBKFreq[11][103] = 64; GBKFreq[69][159] = 63; GBKFreq[53][142] = 62;
	 * GBKFreq[7][6] = 61; GBKFreq[84][59] = 60; GBKFreq[86][3] = 59;
	 * GBKFreq[64][144] = 58; GBKFreq[1][187] = 57; GBKFreq[82][128] = 56;
	 * GBKFreq[3][66] = 55; GBKFreq[68][133] = 54; GBKFreq[55][167] = 53;
	 * GBKFreq[52][130] = 52; GBKFreq[61][133] = 51; GBKFreq[72][181] = 50;
	 * GBKFreq[25][98] = 49; GBKFreq[84][149] = 48; GBKFreq[91][91] = 47;
	 * GBKFreq[47][188] = 46; GBKFreq[68][130] = 45; GBKFreq[22][44] = 44;
	 * GBKFreq[81][121] = 43; GBKFreq[72][140] = 42; GBKFreq[55][133] = 41;
	 * GBKFreq[55][185] = 40; GBKFreq[56][105] = 39; GBKFreq[60][30] = 38;
	 * GBKFreq[70][103] = 37; GBKFreq[62][141] = 36; GBKFreq[70][144] = 35;
	 * GBKFreq[59][111] = 34; GBKFreq[54][17] = 33; GBKFreq[18][190] = 32;
	 * GBKFreq[65][164] = 31; GBKFreq[83][125] = 30; GBKFreq[61][121] = 29;
	 * GBKFreq[48][13] = 28; GBKFreq[51][189] = 27; GBKFreq[65][68] = 26;
	 * GBKFreq[7][0] = 25; GBKFreq[76][188] = 24; GBKFreq[85][117] = 23;
	 * GBKFreq[45][33] = 22; GBKFreq[78][187] = 21; GBKFreq[106][48] = 20;
	 * GBKFreq[59][52] = 19; GBKFreq[86][185] = 18; GBKFreq[84][121] = 17;
	 * GBKFreq[82][189] = 16; GBKFreq[68][156] = 15; GBKFreq[55][125] = 14;
	 * GBKFreq[65][175] = 13; GBKFreq[7][140] = 12; GBKFreq[50][106] = 11;
	 * GBKFreq[59][124] = 10; GBKFreq[67][115] = 9; GBKFreq[82][114] = 8;
	 * GBKFreq[74][121] = 7; GBKFreq[106][69] = 6; GBKFreq[94][27] = 5;
	 * GBKFreq[78][98] = 4; GBKFreq[85][186] = 3; GBKFreq[108][90] = 2;
	 * GBKFreq[62][160] = 1; GBKFreq[60][169] = 0;
	 */
	detect.KRFreq[31][43] = 600
	detect.KRFreq[19][56] = 599
	detect.KRFreq[38][46] = 598
	detect.KRFreq[3][3] = 597
	detect.KRFreq[29][77] = 596
	detect.KRFreq[19][33] = 595
	detect.KRFreq[30][0] = 594
	detect.KRFreq[29][89] = 593
	detect.KRFreq[31][26] = 592
	detect.KRFreq[31][38] = 591
	detect.KRFreq[32][85] = 590
	detect.KRFreq[15][0] = 589
	detect.KRFreq[16][54] = 588
	detect.KRFreq[15][76] = 587
	detect.KRFreq[31][25] = 586
	detect.KRFreq[23][13] = 585
	detect.KRFreq[28][34] = 584
	detect.KRFreq[18][9] = 583
	detect.KRFreq[29][37] = 582
	detect.KRFreq[22][45] = 581
	detect.KRFreq[19][46] = 580
	detect.KRFreq[16][65] = 579
	detect.KRFreq[23][5] = 578
	detect.KRFreq[26][70] = 577
	detect.KRFreq[31][53] = 576
	detect.KRFreq[27][12] = 575
	detect.KRFreq[30][67] = 574
	detect.KRFreq[31][57] = 573
	detect.KRFreq[20][20] = 572
	detect.KRFreq[30][31] = 571
	detect.KRFreq[20][72] = 570
	detect.KRFreq[15][51] = 569
	detect.KRFreq[3][8] = 568
	detect.KRFreq[32][53] = 567
	detect.KRFreq[27][85] = 566
	detect.KRFreq[25][23] = 565
	detect.KRFreq[15][44] = 564
	detect.KRFreq[32][3] = 563
	detect.KRFreq[31][68] = 562
	detect.KRFreq[30][24] = 561
	detect.KRFreq[29][49] = 560
	detect.KRFreq[27][49] = 559
	detect.KRFreq[23][23] = 558
	detect.KRFreq[31][91] = 557
	detect.KRFreq[31][46] = 556
	detect.KRFreq[19][74] = 555
	detect.KRFreq[27][27] = 554
	detect.KRFreq[3][17] = 553
	detect.KRFreq[20][38] = 552
	detect.KRFreq[21][82] = 551
	detect.KRFreq[28][25] = 550
	detect.KRFreq[32][5] = 549
	detect.KRFreq[31][23] = 548
	detect.KRFreq[25][45] = 547
	detect.KRFreq[32][87] = 546
	detect.KRFreq[18][26] = 545
	detect.KRFreq[24][10] = 544
	detect.KRFreq[26][82] = 543
	detect.KRFreq[15][89] = 542
	detect.KRFreq[28][36] = 541
	detect.KRFreq[28][31] = 540
	detect.KRFreq[16][23] = 539
	detect.KRFreq[16][77] = 538
	detect.KRFreq[19][84] = 537
	detect.KRFreq[23][72] = 536
	detect.KRFreq[38][48] = 535
	detect.KRFreq[23][2] = 534
	detect.KRFreq[30][20] = 533
	detect.KRFreq[38][47] = 532
	detect.KRFreq[39][12] = 531
	detect.KRFreq[23][21] = 530
	detect.KRFreq[18][17] = 529
	detect.KRFreq[30][87] = 528
	detect.KRFreq[29][62] = 527
	detect.KRFreq[29][87] = 526
	detect.KRFreq[34][53] = 525
	detect.KRFreq[32][29] = 524
	detect.KRFreq[35][0] = 523
	detect.KRFreq[24][43] = 522
	detect.KRFreq[36][44] = 521
	detect.KRFreq[20][30] = 520
	detect.KRFreq[39][86] = 519
	detect.KRFreq[22][14] = 518
	detect.KRFreq[29][39] = 517
	detect.KRFreq[28][38] = 516
	detect.KRFreq[23][79] = 515
	detect.KRFreq[24][56] = 514
	detect.KRFreq[29][63] = 513
	detect.KRFreq[31][45] = 512
	detect.KRFreq[23][26] = 511
	detect.KRFreq[15][87] = 510
	detect.KRFreq[30][74] = 509
	detect.KRFreq[24][69] = 508
	detect.KRFreq[20][4] = 507
	detect.KRFreq[27][50] = 506
	detect.KRFreq[30][75] = 505
	detect.KRFreq[24][13] = 504
	detect.KRFreq[30][8] = 503
	detect.KRFreq[31][6] = 502
	detect.KRFreq[25][80] = 501
	detect.KRFreq[36][8] = 500
	detect.KRFreq[15][18] = 499
	detect.KRFreq[39][23] = 498
	detect.KRFreq[16][24] = 497
	detect.KRFreq[31][89] = 496
	detect.KRFreq[15][71] = 495
	detect.KRFreq[15][57] = 494
	detect.KRFreq[30][11] = 493
	detect.KRFreq[15][36] = 492
	detect.KRFreq[16][60] = 491
	detect.KRFreq[24][45] = 490
	detect.KRFreq[37][35] = 489
	detect.KRFreq[24][87] = 488
	detect.KRFreq[20][45] = 487
	detect.KRFreq[31][90] = 486
	detect.KRFreq[32][21] = 485
	detect.KRFreq[19][70] = 484
	detect.KRFreq[24][15] = 483
	detect.KRFreq[26][92] = 482
	detect.KRFreq[37][13] = 481
	detect.KRFreq[39][2] = 480
	detect.KRFreq[23][70] = 479
	detect.KRFreq[27][25] = 478
	detect.KRFreq[15][69] = 477
	detect.KRFreq[19][61] = 476
	detect.KRFreq[31][58] = 475
	detect.KRFreq[24][57] = 474
	detect.KRFreq[36][74] = 473
	detect.KRFreq[21][6] = 472
	detect.KRFreq[30][44] = 471
	detect.KRFreq[15][91] = 470
	detect.KRFreq[27][16] = 469
	detect.KRFreq[29][42] = 468
	detect.KRFreq[33][86] = 467
	detect.KRFreq[29][41] = 466
	detect.KRFreq[20][68] = 465
	detect.KRFreq[25][47] = 464
	detect.KRFreq[22][0] = 463
	detect.KRFreq[18][14] = 462
	detect.KRFreq[31][28] = 461
	detect.KRFreq[15][2] = 460
	detect.KRFreq[23][76] = 459
	detect.KRFreq[38][32] = 458
	detect.KRFreq[29][82] = 457
	detect.KRFreq[21][86] = 456
	detect.KRFreq[24][62] = 455
	detect.KRFreq[31][64] = 454
	detect.KRFreq[38][26] = 453
	detect.KRFreq[32][86] = 452
	detect.KRFreq[22][32] = 451
	detect.KRFreq[19][59] = 450
	detect.KRFreq[34][18] = 449
	detect.KRFreq[18][54] = 448
	detect.KRFreq[38][63] = 447
	detect.KRFreq[36][23] = 446
	detect.KRFreq[35][35] = 445
	detect.KRFreq[32][62] = 444
	detect.KRFreq[28][35] = 443
	detect.KRFreq[27][13] = 442
	detect.KRFreq[31][59] = 441
	detect.KRFreq[29][29] = 440
	detect.KRFreq[15][64] = 439
	detect.KRFreq[26][84] = 438
	detect.KRFreq[21][90] = 437
	detect.KRFreq[20][24] = 436
	detect.KRFreq[16][18] = 435
	detect.KRFreq[22][23] = 434
	detect.KRFreq[31][14] = 433
	detect.KRFreq[15][1] = 432
	detect.KRFreq[18][63] = 431
	detect.KRFreq[19][10] = 430
	detect.KRFreq[25][49] = 429
	detect.KRFreq[36][57] = 428
	detect.KRFreq[20][22] = 427
	detect.KRFreq[15][15] = 426
	detect.KRFreq[31][51] = 425
	detect.KRFreq[24][60] = 424
	detect.KRFreq[31][70] = 423
	detect.KRFreq[15][7] = 422
	detect.KRFreq[28][40] = 421
	detect.KRFreq[18][41] = 420
	detect.KRFreq[15][38] = 419
	detect.KRFreq[32][0] = 418
	detect.KRFreq[19][51] = 417
	detect.KRFreq[34][62] = 416
	detect.KRFreq[16][27] = 415
	detect.KRFreq[20][70] = 414
	detect.KRFreq[22][33] = 413
	detect.KRFreq[26][73] = 412
	detect.KRFreq[20][79] = 411
	detect.KRFreq[23][6] = 410
	detect.KRFreq[24][85] = 409
	detect.KRFreq[38][51] = 408
	detect.KRFreq[29][88] = 407
	detect.KRFreq[38][55] = 406
	detect.KRFreq[32][32] = 405
	detect.KRFreq[27][18] = 404
	detect.KRFreq[23][87] = 403
	detect.KRFreq[35][6] = 402
	detect.KRFreq[34][27] = 401
	detect.KRFreq[39][35] = 400
	detect.KRFreq[30][88] = 399
	detect.KRFreq[32][92] = 398
	detect.KRFreq[32][49] = 397
	detect.KRFreq[24][61] = 396
	detect.KRFreq[18][74] = 395
	detect.KRFreq[23][77] = 394
	detect.KRFreq[23][50] = 393
	detect.KRFreq[23][32] = 392
	detect.KRFreq[23][36] = 391
	detect.KRFreq[38][38] = 390
	detect.KRFreq[29][86] = 389
	detect.KRFreq[36][15] = 388
	detect.KRFreq[31][50] = 387
	detect.KRFreq[15][86] = 386
	detect.KRFreq[39][13] = 385
	detect.KRFreq[34][26] = 384
	detect.KRFreq[19][34] = 383
	detect.KRFreq[16][3] = 382
	detect.KRFreq[26][93] = 381
	detect.KRFreq[19][67] = 380
	detect.KRFreq[24][72] = 379
	detect.KRFreq[29][17] = 378
	detect.KRFreq[23][24] = 377
	detect.KRFreq[25][19] = 376
	detect.KRFreq[18][65] = 375
	detect.KRFreq[30][78] = 374
	detect.KRFreq[27][52] = 373
	detect.KRFreq[22][18] = 372
	detect.KRFreq[16][38] = 371
	detect.KRFreq[21][26] = 370
	detect.KRFreq[34][20] = 369
	detect.KRFreq[15][42] = 368
	detect.KRFreq[16][71] = 367
	detect.KRFreq[17][17] = 366
	detect.KRFreq[24][71] = 365
	detect.KRFreq[18][84] = 364
	detect.KRFreq[15][40] = 363
	detect.KRFreq[31][62] = 362
	detect.KRFreq[15][8] = 361
	detect.KRFreq[16][69] = 360
	detect.KRFreq[29][79] = 359
	detect.KRFreq[38][91] = 358
	detect.KRFreq[31][92] = 357
	detect.KRFreq[20][77] = 356
	detect.KRFreq[3][16] = 355
	detect.KRFreq[27][87] = 354
	detect.KRFreq[16][25] = 353
	detect.KRFreq[36][33] = 352
	detect.KRFreq[37][76] = 351
	detect.KRFreq[30][12] = 350
	detect.KRFreq[26][75] = 349
	detect.KRFreq[25][14] = 348
	detect.KRFreq[32][26] = 347
	detect.KRFreq[23][22] = 346
	detect.KRFreq[20][90] = 345
	detect.KRFreq[19][8] = 344
	detect.KRFreq[38][41] = 343
	detect.KRFreq[34][2] = 342
	detect.KRFreq[39][4] = 341
	detect.KRFreq[27][89] = 340
	detect.KRFreq[28][41] = 339
	detect.KRFreq[28][44] = 338
	detect.KRFreq[24][92] = 337
	detect.KRFreq[34][65] = 336
	detect.KRFreq[39][14] = 335
	detect.KRFreq[21][38] = 334
	detect.KRFreq[19][31] = 333
	detect.KRFreq[37][39] = 332
	detect.KRFreq[33][41] = 331
	detect.KRFreq[38][4] = 330
	detect.KRFreq[23][80] = 329
	detect.KRFreq[25][24] = 328
	detect.KRFreq[37][17] = 327
	detect.KRFreq[22][16] = 326
	detect.KRFreq[22][46] = 325
	detect.KRFreq[33][91] = 324
	detect.KRFreq[24][89] = 323
	detect.KRFreq[30][52] = 322
	detect.KRFreq[29][38] = 321
	detect.KRFreq[38][85] = 320
	detect.KRFreq[15][12] = 319
	detect.KRFreq[27][58] = 318
	detect.KRFreq[29][52] = 317
	detect.KRFreq[37][38] = 316
	detect.KRFreq[34][41] = 315
	detect.KRFreq[31][65] = 314
	detect.KRFreq[29][53] = 313
	detect.KRFreq[22][47] = 312
	detect.KRFreq[22][19] = 311
	detect.KRFreq[26][0] = 310
	detect.KRFreq[37][86] = 309
	detect.KRFreq[35][4] = 308
	detect.KRFreq[36][54] = 307
	detect.KRFreq[20][76] = 306
	detect.KRFreq[30][9] = 305
	detect.KRFreq[30][33] = 304
	detect.KRFreq[23][17] = 303
	detect.KRFreq[23][33] = 302
	detect.KRFreq[38][52] = 301
	detect.KRFreq[15][19] = 300
	detect.KRFreq[28][45] = 299
	detect.KRFreq[29][78] = 298
	detect.KRFreq[23][15] = 297
	detect.KRFreq[33][5] = 296
	detect.KRFreq[17][40] = 295
	detect.KRFreq[30][83] = 294
	detect.KRFreq[18][1] = 293
	detect.KRFreq[30][81] = 292
	detect.KRFreq[19][40] = 291
	detect.KRFreq[24][47] = 290
	detect.KRFreq[17][56] = 289
	detect.KRFreq[39][80] = 288
	detect.KRFreq[30][46] = 287
	detect.KRFreq[16][61] = 286
	detect.KRFreq[26][78] = 285
	detect.KRFreq[26][57] = 284
	detect.KRFreq[20][46] = 283
	detect.KRFreq[25][15] = 282
	detect.KRFreq[25][91] = 281
	detect.KRFreq[21][83] = 280
	detect.KRFreq[30][77] = 279
	detect.KRFreq[35][30] = 278
	detect.KRFreq[30][34] = 277
	detect.KRFreq[20][69] = 276
	detect.KRFreq[35][10] = 275
	detect.KRFreq[29][70] = 274
	detect.KRFreq[22][50] = 273
	detect.KRFreq[18][0] = 272
	detect.KRFreq[22][64] = 271
	detect.KRFreq[38][65] = 270
	detect.KRFreq[22][70] = 269
	detect.KRFreq[24][58] = 268
	detect.KRFreq[19][66] = 267
	detect.KRFreq[30][59] = 266
	detect.KRFreq[37][14] = 265
	detect.KRFreq[16][56] = 264
	detect.KRFreq[29][85] = 263
	detect.KRFreq[31][15] = 262
	detect.KRFreq[36][84] = 261
	detect.KRFreq[39][15] = 260
	detect.KRFreq[39][90] = 259
	detect.KRFreq[18][12] = 258
	detect.KRFreq[21][93] = 257
	detect.KRFreq[24][66] = 256
	detect.KRFreq[27][90] = 255
	detect.KRFreq[25][90] = 254
	detect.KRFreq[22][24] = 253
	detect.KRFreq[36][67] = 252
	detect.KRFreq[33][90] = 251
	detect.KRFreq[15][60] = 250
	detect.KRFreq[23][85] = 249
	detect.KRFreq[34][1] = 248
	detect.KRFreq[39][37] = 247
	detect.KRFreq[21][18] = 246
	detect.KRFreq[34][4] = 245
	detect.KRFreq[28][33] = 244
	detect.KRFreq[15][13] = 243
	detect.KRFreq[32][22] = 242
	detect.KRFreq[30][76] = 241
	detect.KRFreq[20][21] = 240
	detect.KRFreq[38][66] = 239
	detect.KRFreq[32][55] = 238
	detect.KRFreq[32][89] = 237
	detect.KRFreq[25][26] = 236
	detect.KRFreq[16][80] = 235
	detect.KRFreq[15][43] = 234
	detect.KRFreq[38][54] = 233
	detect.KRFreq[39][68] = 232
	detect.KRFreq[22][88] = 231
	detect.KRFreq[21][84] = 230
	detect.KRFreq[21][17] = 229
	detect.KRFreq[20][28] = 228
	detect.KRFreq[32][1] = 227
	detect.KRFreq[33][87] = 226
	detect.KRFreq[38][71] = 225
	detect.KRFreq[37][47] = 224
	detect.KRFreq[18][77] = 223
	detect.KRFreq[37][58] = 222
	detect.KRFreq[34][74] = 221
	detect.KRFreq[32][54] = 220
	detect.KRFreq[27][33] = 219
	detect.KRFreq[32][93] = 218
	detect.KRFreq[23][51] = 217
	detect.KRFreq[20][57] = 216
	detect.KRFreq[22][37] = 215
	detect.KRFreq[39][10] = 214
	detect.KRFreq[39][17] = 213
	detect.KRFreq[33][4] = 212
	detect.KRFreq[32][84] = 211
	detect.KRFreq[34][3] = 210
	detect.KRFreq[28][27] = 209
	detect.KRFreq[15][79] = 208
	detect.KRFreq[34][21] = 207
	detect.KRFreq[34][69] = 206
	detect.KRFreq[21][62] = 205
	detect.KRFreq[36][24] = 204
	detect.KRFreq[16][89] = 203
	detect.KRFreq[18][48] = 202
	detect.KRFreq[38][15] = 201
	detect.KRFreq[36][58] = 200
	detect.KRFreq[21][56] = 199
	detect.KRFreq[34][48] = 198
	detect.KRFreq[21][15] = 197
	detect.KRFreq[39][3] = 196
	detect.KRFreq[16][44] = 195
	detect.KRFreq[18][79] = 194
	detect.KRFreq[25][13] = 193
	detect.KRFreq[29][47] = 192
	detect.KRFreq[38][88] = 191
	detect.KRFreq[20][71] = 190
	detect.KRFreq[16][58] = 189
	detect.KRFreq[35][57] = 188
	detect.KRFreq[29][30] = 187
	detect.KRFreq[29][23] = 186
	detect.KRFreq[34][93] = 185
	detect.KRFreq[30][85] = 184
	detect.KRFreq[15][80] = 183
	detect.KRFreq[32][78] = 182
	detect.KRFreq[37][82] = 181
	detect.KRFreq[22][40] = 180
	detect.KRFreq[21][69] = 179
	detect.KRFreq[26][85] = 178
	detect.KRFreq[31][31] = 177
	detect.KRFreq[28][64] = 176
	detect.KRFreq[38][13] = 175
	detect.KRFreq[25][2] = 174
	detect.KRFreq[22][34] = 173
	detect.KRFreq[28][28] = 172
	detect.KRFreq[24][91] = 171
	detect.KRFreq[33][74] = 170
	detect.KRFreq[29][40] = 169
	detect.KRFreq[15][77] = 168
	detect.KRFreq[32][80] = 167
	detect.KRFreq[30][41] = 166
	detect.KRFreq[23][30] = 165
	detect.KRFreq[24][63] = 164
	detect.KRFreq[30][53] = 163
	detect.KRFreq[39][70] = 162
	detect.KRFreq[23][61] = 161
	detect.KRFreq[37][27] = 160
	detect.KRFreq[16][55] = 159
	detect.KRFreq[22][74] = 158
	detect.KRFreq[26][50] = 157
	detect.KRFreq[16][10] = 156
	detect.KRFreq[34][63] = 155
	detect.KRFreq[35][14] = 154
	detect.KRFreq[17][7] = 153
	detect.KRFreq[15][59] = 152
	detect.KRFreq[27][23] = 151
	detect.KRFreq[18][70] = 150
	detect.KRFreq[32][56] = 149
	detect.KRFreq[37][87] = 148
	detect.KRFreq[17][61] = 147
	detect.KRFreq[18][83] = 146
	detect.KRFreq[23][86] = 145
	detect.KRFreq[17][31] = 144
	detect.KRFreq[23][83] = 143
	detect.KRFreq[35][2] = 142
	detect.KRFreq[18][64] = 141
	detect.KRFreq[27][43] = 140
	detect.KRFreq[32][42] = 139
	detect.KRFreq[25][76] = 138
	detect.KRFreq[19][85] = 137
	detect.KRFreq[37][81] = 136
	detect.KRFreq[38][83] = 135
	detect.KRFreq[35][7] = 134
	detect.KRFreq[16][51] = 133
	detect.KRFreq[27][22] = 132
	detect.KRFreq[16][76] = 131
	detect.KRFreq[22][4] = 130
	detect.KRFreq[38][84] = 129
	detect.KRFreq[17][83] = 128
	detect.KRFreq[24][46] = 127
	detect.KRFreq[33][15] = 126
	detect.KRFreq[20][48] = 125
	detect.KRFreq[17][30] = 124
	detect.KRFreq[30][93] = 123
	detect.KRFreq[28][11] = 122
	detect.KRFreq[28][30] = 121
	detect.KRFreq[15][62] = 120
	detect.KRFreq[17][87] = 119
	detect.KRFreq[32][81] = 118
	detect.KRFreq[23][37] = 117
	detect.KRFreq[30][22] = 116
	detect.KRFreq[32][66] = 115
	detect.KRFreq[33][78] = 114
	detect.KRFreq[21][4] = 113
	detect.KRFreq[31][17] = 112
	detect.KRFreq[39][61] = 111
	detect.KRFreq[18][76] = 110
	detect.KRFreq[15][85] = 109
	detect.KRFreq[31][47] = 108
	detect.KRFreq[19][57] = 107
	detect.KRFreq[23][55] = 106
	detect.KRFreq[27][29] = 105
	detect.KRFreq[29][46] = 104
	detect.KRFreq[33][0] = 103
	detect.KRFreq[16][83] = 102
	detect.KRFreq[39][78] = 101
	detect.KRFreq[32][77] = 100
	detect.KRFreq[36][25] = 99
	detect.KRFreq[34][19] = 98
	detect.KRFreq[38][49] = 97
	detect.KRFreq[19][25] = 96
	detect.KRFreq[23][53] = 95
	detect.KRFreq[28][43] = 94
	detect.KRFreq[31][44] = 93
	detect.KRFreq[36][34] = 92
	detect.KRFreq[16][34] = 91
	detect.KRFreq[35][1] = 90
	detect.KRFreq[19][87] = 89
	detect.KRFreq[18][53] = 88
	detect.KRFreq[29][54] = 87
	detect.KRFreq[22][41] = 86
	detect.KRFreq[38][18] = 85
	detect.KRFreq[22][2] = 84
	detect.KRFreq[20][3] = 83
	detect.KRFreq[39][69] = 82
	detect.KRFreq[30][29] = 81
	detect.KRFreq[28][19] = 80
	detect.KRFreq[29][90] = 79
	detect.KRFreq[17][86] = 78
	detect.KRFreq[15][9] = 77
	detect.KRFreq[39][73] = 76
	detect.KRFreq[15][37] = 75
	detect.KRFreq[35][40] = 74
	detect.KRFreq[33][77] = 73
	detect.KRFreq[27][86] = 72
	detect.KRFreq[36][79] = 71
	detect.KRFreq[23][18] = 70
	detect.KRFreq[34][87] = 69
	detect.KRFreq[39][24] = 68
	detect.KRFreq[26][8] = 67
	detect.KRFreq[33][48] = 66
	detect.KRFreq[39][30] = 65
	detect.KRFreq[33][28] = 64
	detect.KRFreq[16][67] = 63
	detect.KRFreq[31][78] = 62
	detect.KRFreq[32][23] = 61
	detect.KRFreq[24][55] = 60
	detect.KRFreq[30][68] = 59
	detect.KRFreq[18][60] = 58
	detect.KRFreq[15][17] = 57
	detect.KRFreq[23][34] = 56
	detect.KRFreq[20][49] = 55
	detect.KRFreq[15][78] = 54
	detect.KRFreq[24][14] = 53
	detect.KRFreq[19][41] = 52
	detect.KRFreq[31][55] = 51
	detect.KRFreq[21][39] = 50
	detect.KRFreq[35][9] = 49
	detect.KRFreq[30][15] = 48
	detect.KRFreq[20][52] = 47
	detect.KRFreq[35][71] = 46
	detect.KRFreq[20][7] = 45
	detect.KRFreq[29][72] = 44
	detect.KRFreq[37][77] = 43
	detect.KRFreq[22][35] = 42
	detect.KRFreq[20][61] = 41
	detect.KRFreq[31][60] = 40
	detect.KRFreq[20][93] = 39
	detect.KRFreq[27][92] = 38
	detect.KRFreq[28][16] = 37
	detect.KRFreq[36][26] = 36
	detect.KRFreq[18][89] = 35
	detect.KRFreq[21][63] = 34
	detect.KRFreq[22][52] = 33
	detect.KRFreq[24][65] = 32
	detect.KRFreq[31][8] = 31
	detect.KRFreq[31][49] = 30
	detect.KRFreq[33][30] = 29
	detect.KRFreq[37][15] = 28
	detect.KRFreq[18][18] = 27
	detect.KRFreq[25][50] = 26
	detect.KRFreq[29][20] = 25
	detect.KRFreq[35][48] = 24
	detect.KRFreq[38][75] = 23
	detect.KRFreq[26][83] = 22
	detect.KRFreq[21][87] = 21
	detect.KRFreq[27][71] = 20
	detect.KRFreq[32][91] = 19
	detect.KRFreq[25][73] = 18
	detect.KRFreq[16][84] = 17
	detect.KRFreq[25][31] = 16
	detect.KRFreq[17][90] = 15
	detect.KRFreq[18][40] = 14
	detect.KRFreq[17][77] = 13
	detect.KRFreq[17][35] = 12
	detect.KRFreq[23][52] = 11
	detect.KRFreq[23][35] = 10
	detect.KRFreq[16][5] = 9
	detect.KRFreq[23][58] = 8
	detect.KRFreq[19][60] = 7
	detect.KRFreq[30][32] = 6
	detect.KRFreq[38][34] = 5
	detect.KRFreq[23][4] = 4
	detect.KRFreq[23][1] = 3
	detect.KRFreq[27][57] = 2
	detect.KRFreq[39][38] = 1
	detect.KRFreq[32][33] = 0
	detect.JPFreq[3][74] = 600
	detect.JPFreq[3][45] = 599
	detect.JPFreq[3][3] = 598
	detect.JPFreq[3][24] = 597
	detect.JPFreq[3][30] = 596
	detect.JPFreq[3][42] = 595
	detect.JPFreq[3][46] = 594
	detect.JPFreq[3][39] = 593
	detect.JPFreq[3][11] = 592
	detect.JPFreq[3][37] = 591
	detect.JPFreq[3][38] = 590
	detect.JPFreq[3][31] = 589
	detect.JPFreq[3][41] = 588
	detect.JPFreq[3][5] = 587
	detect.JPFreq[3][10] = 586
	detect.JPFreq[3][75] = 585
	detect.JPFreq[3][65] = 584
	detect.JPFreq[3][72] = 583
	detect.JPFreq[37][91] = 582
	detect.JPFreq[0][27] = 581
	detect.JPFreq[3][18] = 580
	detect.JPFreq[3][22] = 579
	detect.JPFreq[3][61] = 578
	detect.JPFreq[3][14] = 577
	detect.JPFreq[24][80] = 576
	detect.JPFreq[4][82] = 575
	detect.JPFreq[17][80] = 574
	detect.JPFreq[30][44] = 573
	detect.JPFreq[3][73] = 572
	detect.JPFreq[3][64] = 571
	detect.JPFreq[38][14] = 570
	detect.JPFreq[33][70] = 569
	detect.JPFreq[3][1] = 568
	detect.JPFreq[3][16] = 567
	detect.JPFreq[3][35] = 566
	detect.JPFreq[3][40] = 565
	detect.JPFreq[4][74] = 564
	detect.JPFreq[4][24] = 563
	detect.JPFreq[42][59] = 562
	detect.JPFreq[3][7] = 561
	detect.JPFreq[3][71] = 560
	detect.JPFreq[3][12] = 559
	detect.JPFreq[15][75] = 558
	detect.JPFreq[3][20] = 557
	detect.JPFreq[4][39] = 556
	detect.JPFreq[34][69] = 555
	detect.JPFreq[3][28] = 554
	detect.JPFreq[35][24] = 553
	detect.JPFreq[3][82] = 552
	detect.JPFreq[28][47] = 551
	detect.JPFreq[3][67] = 550
	detect.JPFreq[37][16] = 549
	detect.JPFreq[26][93] = 548
	detect.JPFreq[4][1] = 547
	detect.JPFreq[26][85] = 546
	detect.JPFreq[31][14] = 545
	detect.JPFreq[4][3] = 544
	detect.JPFreq[4][72] = 543
	detect.JPFreq[24][51] = 542
	detect.JPFreq[27][51] = 541
	detect.JPFreq[27][49] = 540
	detect.JPFreq[22][77] = 539
	detect.JPFreq[27][10] = 538
	detect.JPFreq[29][68] = 537
	detect.JPFreq[20][35] = 536
	detect.JPFreq[41][11] = 535
	detect.JPFreq[24][70] = 534
	detect.JPFreq[36][61] = 533
	detect.JPFreq[31][23] = 532
	detect.JPFreq[43][16] = 531
	detect.JPFreq[23][68] = 530
	detect.JPFreq[32][15] = 529
	detect.JPFreq[3][32] = 528
	detect.JPFreq[19][53] = 527
	detect.JPFreq[40][83] = 526
	detect.JPFreq[4][14] = 525
	detect.JPFreq[36][9] = 524
	detect.JPFreq[4][73] = 523
	detect.JPFreq[23][10] = 522
	detect.JPFreq[3][63] = 521
	detect.JPFreq[39][14] = 520
	detect.JPFreq[3][78] = 519
	detect.JPFreq[33][47] = 518
	detect.JPFreq[21][39] = 517
	detect.JPFreq[34][46] = 516
	detect.JPFreq[36][75] = 515
	detect.JPFreq[41][92] = 514
	detect.JPFreq[37][93] = 513
	detect.JPFreq[4][34] = 512
	detect.JPFreq[15][86] = 511
	detect.JPFreq[46][1] = 510
	detect.JPFreq[37][65] = 509
	detect.JPFreq[3][62] = 508
	detect.JPFreq[32][73] = 507
	detect.JPFreq[21][65] = 506
	detect.JPFreq[29][75] = 505
	detect.JPFreq[26][51] = 504
	detect.JPFreq[3][34] = 503
	detect.JPFreq[4][10] = 502
	detect.JPFreq[30][22] = 501
	detect.JPFreq[35][73] = 500
	detect.JPFreq[17][82] = 499
	detect.JPFreq[45][8] = 498
	detect.JPFreq[27][73] = 497
	detect.JPFreq[18][55] = 496
	detect.JPFreq[25][2] = 495
	detect.JPFreq[3][26] = 494
	detect.JPFreq[45][46] = 493
	detect.JPFreq[4][22] = 492
	detect.JPFreq[4][40] = 491
	detect.JPFreq[18][10] = 490
	detect.JPFreq[32][9] = 489
	detect.JPFreq[26][49] = 488
	detect.JPFreq[3][47] = 487
	detect.JPFreq[24][65] = 486
	detect.JPFreq[4][76] = 485
	detect.JPFreq[43][67] = 484
	detect.JPFreq[3][9] = 483
	detect.JPFreq[41][37] = 482
	detect.JPFreq[33][68] = 481
	detect.JPFreq[43][31] = 480
	detect.JPFreq[19][55] = 479
	detect.JPFreq[4][30] = 478
	detect.JPFreq[27][33] = 477
	detect.JPFreq[16][62] = 476
	detect.JPFreq[36][35] = 475
	detect.JPFreq[37][15] = 474
	detect.JPFreq[27][70] = 473
	detect.JPFreq[22][71] = 472
	detect.JPFreq[33][45] = 471
	detect.JPFreq[31][78] = 470
	detect.JPFreq[43][59] = 469
	detect.JPFreq[32][19] = 468
	detect.JPFreq[17][28] = 467
	detect.JPFreq[40][28] = 466
	detect.JPFreq[20][93] = 465
	detect.JPFreq[18][15] = 464
	detect.JPFreq[4][23] = 463
	detect.JPFreq[3][23] = 462
	detect.JPFreq[26][64] = 461
	detect.JPFreq[44][92] = 460
	detect.JPFreq[17][27] = 459
	detect.JPFreq[3][56] = 458
	detect.JPFreq[25][38] = 457
	detect.JPFreq[23][31] = 456
	detect.JPFreq[35][43] = 455
	detect.JPFreq[4][54] = 454
	detect.JPFreq[35][19] = 453
	detect.JPFreq[22][47] = 452
	detect.JPFreq[42][0] = 451
	detect.JPFreq[23][28] = 450
	detect.JPFreq[46][33] = 449
	detect.JPFreq[36][85] = 448
	detect.JPFreq[31][12] = 447
	detect.JPFreq[3][76] = 446
	detect.JPFreq[4][75] = 445
	detect.JPFreq[36][56] = 444
	detect.JPFreq[4][64] = 443
	detect.JPFreq[25][77] = 442
	detect.JPFreq[15][52] = 441
	detect.JPFreq[33][73] = 440
	detect.JPFreq[3][55] = 439
	detect.JPFreq[43][82] = 438
	detect.JPFreq[27][82] = 437
	detect.JPFreq[20][3] = 436
	detect.JPFreq[40][51] = 435
	detect.JPFreq[3][17] = 434
	detect.JPFreq[27][71] = 433
	detect.JPFreq[4][52] = 432
	detect.JPFreq[44][48] = 431
	detect.JPFreq[27][2] = 430
	detect.JPFreq[17][39] = 429
	detect.JPFreq[31][8] = 428
	detect.JPFreq[44][54] = 427
	detect.JPFreq[43][18] = 426
	detect.JPFreq[43][77] = 425
	detect.JPFreq[4][61] = 424
	detect.JPFreq[19][91] = 423
	detect.JPFreq[31][13] = 422
	detect.JPFreq[44][71] = 421
	detect.JPFreq[20][0] = 420
	detect.JPFreq[23][87] = 419
	detect.JPFreq[21][14] = 418
	detect.JPFreq[29][13] = 417
	detect.JPFreq[3][58] = 416
	detect.JPFreq[26][18] = 415
	detect.JPFreq[4][47] = 414
	detect.JPFreq[4][18] = 413
	detect.JPFreq[3][53] = 412
	detect.JPFreq[26][92] = 411
	detect.JPFreq[21][7] = 410
	detect.JPFreq[4][37] = 409
	detect.JPFreq[4][63] = 408
	detect.JPFreq[36][51] = 407
	detect.JPFreq[4][32] = 406
	detect.JPFreq[28][73] = 405
	detect.JPFreq[4][50] = 404
	detect.JPFreq[41][60] = 403
	detect.JPFreq[23][1] = 402
	detect.JPFreq[36][92] = 401
	detect.JPFreq[15][41] = 400
	detect.JPFreq[21][71] = 399
	detect.JPFreq[41][30] = 398
	detect.JPFreq[32][76] = 397
	detect.JPFreq[17][34] = 396
	detect.JPFreq[26][15] = 395
	detect.JPFreq[26][25] = 394
	detect.JPFreq[31][77] = 393
	detect.JPFreq[31][3] = 392
	detect.JPFreq[46][34] = 391
	detect.JPFreq[27][84] = 390
	detect.JPFreq[23][8] = 389
	detect.JPFreq[16][0] = 388
	detect.JPFreq[28][80] = 387
	detect.JPFreq[26][54] = 386
	detect.JPFreq[33][18] = 385
	detect.JPFreq[31][20] = 384
	detect.JPFreq[31][62] = 383
	detect.JPFreq[30][41] = 382
	detect.JPFreq[33][30] = 381
	detect.JPFreq[45][45] = 380
	detect.JPFreq[37][82] = 379
	detect.JPFreq[15][33] = 378
	detect.JPFreq[20][12] = 377
	detect.JPFreq[18][5] = 376
	detect.JPFreq[28][86] = 375
	detect.JPFreq[30][19] = 374
	detect.JPFreq[42][43] = 373
	detect.JPFreq[36][31] = 372
	detect.JPFreq[17][93] = 371
	detect.JPFreq[4][15] = 370
	detect.JPFreq[21][20] = 369
	detect.JPFreq[23][21] = 368
	detect.JPFreq[28][72] = 367
	detect.JPFreq[4][20] = 366
	detect.JPFreq[26][55] = 365
	detect.JPFreq[21][5] = 364
	detect.JPFreq[19][16] = 363
	detect.JPFreq[23][64] = 362
	detect.JPFreq[40][59] = 361
	detect.JPFreq[37][26] = 360
	detect.JPFreq[26][56] = 359
	detect.JPFreq[4][12] = 358
	detect.JPFreq[33][71] = 357
	detect.JPFreq[32][39] = 356
	detect.JPFreq[38][40] = 355
	detect.JPFreq[22][74] = 354
	detect.JPFreq[3][25] = 353
	detect.JPFreq[15][48] = 352
	detect.JPFreq[41][82] = 351
	detect.JPFreq[41][9] = 350
	detect.JPFreq[25][48] = 349
	detect.JPFreq[31][71] = 348
	detect.JPFreq[43][29] = 347
	detect.JPFreq[26][80] = 346
	detect.JPFreq[4][5] = 345
	detect.JPFreq[18][71] = 344
	detect.JPFreq[29][0] = 343
	detect.JPFreq[43][43] = 342
	detect.JPFreq[23][81] = 341
	detect.JPFreq[4][42] = 340
	detect.JPFreq[44][28] = 339
	detect.JPFreq[23][93] = 338
	detect.JPFreq[17][81] = 337
	detect.JPFreq[25][25] = 336
	detect.JPFreq[41][23] = 335
	detect.JPFreq[34][35] = 334
	detect.JPFreq[4][53] = 333
	detect.JPFreq[28][36] = 332
	detect.JPFreq[4][41] = 331
	detect.JPFreq[25][60] = 330
	detect.JPFreq[23][20] = 329
	detect.JPFreq[3][43] = 328
	detect.JPFreq[24][79] = 327
	detect.JPFreq[29][41] = 326
	detect.JPFreq[30][83] = 325
	detect.JPFreq[3][50] = 324
	detect.JPFreq[22][18] = 323
	detect.JPFreq[18][3] = 322
	detect.JPFreq[39][30] = 321
	detect.JPFreq[4][28] = 320
	detect.JPFreq[21][64] = 319
	detect.JPFreq[4][68] = 318
	detect.JPFreq[17][71] = 317
	detect.JPFreq[27][0] = 316
	detect.JPFreq[39][28] = 315
	detect.JPFreq[30][13] = 314
	detect.JPFreq[36][70] = 313
	detect.JPFreq[20][82] = 312
	detect.JPFreq[33][38] = 311
	detect.JPFreq[44][87] = 310
	detect.JPFreq[34][45] = 309
	detect.JPFreq[4][26] = 308
	detect.JPFreq[24][44] = 307
	detect.JPFreq[38][67] = 306
	detect.JPFreq[38][6] = 305
	detect.JPFreq[30][68] = 304
	detect.JPFreq[15][89] = 303
	detect.JPFreq[24][93] = 302
	detect.JPFreq[40][41] = 301
	detect.JPFreq[38][3] = 300
	detect.JPFreq[28][23] = 299
	detect.JPFreq[26][17] = 298
	detect.JPFreq[4][38] = 297
	detect.JPFreq[22][78] = 296
	detect.JPFreq[15][37] = 295
	detect.JPFreq[25][85] = 294
	detect.JPFreq[4][9] = 293
	detect.JPFreq[4][7] = 292
	detect.JPFreq[27][53] = 291
	detect.JPFreq[39][29] = 290
	detect.JPFreq[41][43] = 289
	detect.JPFreq[25][62] = 288
	detect.JPFreq[4][48] = 287
	detect.JPFreq[28][28] = 286
	detect.JPFreq[21][40] = 285
	detect.JPFreq[36][73] = 284
	detect.JPFreq[26][39] = 283
	detect.JPFreq[22][54] = 282
	detect.JPFreq[33][5] = 281
	detect.JPFreq[19][21] = 280
	detect.JPFreq[46][31] = 279
	detect.JPFreq[20][64] = 278
	detect.JPFreq[26][63] = 277
	detect.JPFreq[22][23] = 276
	detect.JPFreq[25][81] = 275
	detect.JPFreq[4][62] = 274
	detect.JPFreq[37][31] = 273
	detect.JPFreq[40][52] = 272
	detect.JPFreq[29][79] = 271
	detect.JPFreq[41][48] = 270
	detect.JPFreq[31][57] = 269
	detect.JPFreq[32][92] = 268
	detect.JPFreq[36][36] = 267
	detect.JPFreq[27][7] = 266
	detect.JPFreq[35][29] = 265
	detect.JPFreq[37][34] = 264
	detect.JPFreq[34][42] = 263
	detect.JPFreq[27][15] = 262
	detect.JPFreq[33][27] = 261
	detect.JPFreq[31][38] = 260
	detect.JPFreq[19][79] = 259
	detect.JPFreq[4][31] = 258
	detect.JPFreq[4][66] = 257
	detect.JPFreq[17][32] = 256
	detect.JPFreq[26][67] = 255
	detect.JPFreq[16][30] = 254
	detect.JPFreq[26][46] = 253
	detect.JPFreq[24][26] = 252
	detect.JPFreq[35][10] = 251
	detect.JPFreq[18][37] = 250
	detect.JPFreq[3][19] = 249
	detect.JPFreq[33][69] = 248
	detect.JPFreq[31][9] = 247
	detect.JPFreq[45][29] = 246
	detect.JPFreq[3][15] = 245
	detect.JPFreq[18][54] = 244
	detect.JPFreq[3][44] = 243
	detect.JPFreq[31][29] = 242
	detect.JPFreq[18][45] = 241
	detect.JPFreq[38][28] = 240
	detect.JPFreq[24][12] = 239
	detect.JPFreq[35][82] = 238
	detect.JPFreq[17][43] = 237
	detect.JPFreq[28][9] = 236
	detect.JPFreq[23][25] = 235
	detect.JPFreq[44][37] = 234
	detect.JPFreq[23][75] = 233
	detect.JPFreq[23][92] = 232
	detect.JPFreq[0][24] = 231
	detect.JPFreq[19][74] = 230
	detect.JPFreq[45][32] = 229
	detect.JPFreq[16][72] = 228
	detect.JPFreq[16][93] = 227
	detect.JPFreq[45][13] = 226
	detect.JPFreq[24][8] = 225
	detect.JPFreq[25][47] = 224
	detect.JPFreq[28][26] = 223
	detect.JPFreq[43][81] = 222
	detect.JPFreq[32][71] = 221
	detect.JPFreq[18][41] = 220
	detect.JPFreq[26][62] = 219
	detect.JPFreq[41][24] = 218
	detect.JPFreq[40][11] = 217
	detect.JPFreq[43][57] = 216
	detect.JPFreq[34][53] = 215
	detect.JPFreq[20][32] = 214
	detect.JPFreq[34][43] = 213
	detect.JPFreq[41][91] = 212
	detect.JPFreq[29][57] = 211
	detect.JPFreq[15][43] = 210
	detect.JPFreq[22][89] = 209
	detect.JPFreq[33][83] = 208
	detect.JPFreq[43][20] = 207
	detect.JPFreq[25][58] = 206
	detect.JPFreq[30][30] = 205
	detect.JPFreq[4][56] = 204
	detect.JPFreq[17][64] = 203
	detect.JPFreq[23][0] = 202
	detect.JPFreq[44][12] = 201
	detect.JPFreq[25][37] = 200
	detect.JPFreq[35][13] = 199
	detect.JPFreq[20][30] = 198
	detect.JPFreq[21][84] = 197
	detect.JPFreq[29][14] = 196
	detect.JPFreq[30][5] = 195
	detect.JPFreq[37][2] = 194
	detect.JPFreq[4][78] = 193
	detect.JPFreq[29][78] = 192
	detect.JPFreq[29][84] = 191
	detect.JPFreq[32][86] = 190
	detect.JPFreq[20][68] = 189
	detect.JPFreq[30][39] = 188
	detect.JPFreq[15][69] = 187
	detect.JPFreq[4][60] = 186
	detect.JPFreq[20][61] = 185
	detect.JPFreq[41][67] = 184
	detect.JPFreq[16][35] = 183
	detect.JPFreq[36][57] = 182
	detect.JPFreq[39][80] = 181
	detect.JPFreq[4][59] = 180
	detect.JPFreq[4][44] = 179
	detect.JPFreq[40][54] = 178
	detect.JPFreq[30][8] = 177
	detect.JPFreq[44][30] = 176
	detect.JPFreq[31][93] = 175
	detect.JPFreq[31][47] = 174
	detect.JPFreq[16][70] = 173
	detect.JPFreq[21][0] = 172
	detect.JPFreq[17][35] = 171
	detect.JPFreq[21][67] = 170
	detect.JPFreq[44][18] = 169
	detect.JPFreq[36][29] = 168
	detect.JPFreq[18][67] = 167
	detect.JPFreq[24][28] = 166
	detect.JPFreq[36][24] = 165
	detect.JPFreq[23][5] = 164
	detect.JPFreq[31][65] = 163
	detect.JPFreq[26][59] = 162
	detect.JPFreq[28][2] = 161
	detect.JPFreq[39][69] = 160
	detect.JPFreq[42][40] = 159
	detect.JPFreq[37][80] = 158
	detect.JPFreq[15][66] = 157
	detect.JPFreq[34][38] = 156
	detect.JPFreq[28][48] = 155
	detect.JPFreq[37][77] = 154
	detect.JPFreq[29][34] = 153
	detect.JPFreq[33][12] = 152
	detect.JPFreq[4][65] = 151
	detect.JPFreq[30][31] = 150
	detect.JPFreq[27][92] = 149
	detect.JPFreq[4][2] = 148
	detect.JPFreq[4][51] = 147
	detect.JPFreq[23][77] = 146
	detect.JPFreq[4][35] = 145
	detect.JPFreq[3][13] = 144
	detect.JPFreq[26][26] = 143
	detect.JPFreq[44][4] = 142
	detect.JPFreq[39][53] = 141
	detect.JPFreq[20][11] = 140
	detect.JPFreq[40][33] = 139
	detect.JPFreq[45][7] = 138
	detect.JPFreq[4][70] = 137
	detect.JPFreq[3][49] = 136
	detect.JPFreq[20][59] = 135
	detect.JPFreq[21][12] = 134
	detect.JPFreq[33][53] = 133
	detect.JPFreq[20][14] = 132
	detect.JPFreq[37][18] = 131
	detect.JPFreq[18][17] = 130
	detect.JPFreq[36][23] = 129
	detect.JPFreq[18][57] = 128
	detect.JPFreq[26][74] = 127
	detect.JPFreq[35][2] = 126
	detect.JPFreq[38][58] = 125
	detect.JPFreq[34][68] = 124
	detect.JPFreq[29][81] = 123
	detect.JPFreq[20][69] = 122
	detect.JPFreq[39][86] = 121
	detect.JPFreq[4][16] = 120
	detect.JPFreq[16][49] = 119
	detect.JPFreq[15][72] = 118
	detect.JPFreq[26][35] = 117
	detect.JPFreq[32][14] = 116
	detect.JPFreq[40][90] = 115
	detect.JPFreq[33][79] = 114
	detect.JPFreq[35][4] = 113
	detect.JPFreq[23][33] = 112
	detect.JPFreq[19][19] = 111
	detect.JPFreq[31][41] = 110
	detect.JPFreq[44][1] = 109
	detect.JPFreq[22][56] = 108
	detect.JPFreq[31][27] = 107
	detect.JPFreq[32][18] = 106
	detect.JPFreq[27][32] = 105
	detect.JPFreq[37][39] = 104
	detect.JPFreq[42][11] = 103
	detect.JPFreq[29][71] = 102
	detect.JPFreq[32][58] = 101
	detect.JPFreq[46][10] = 100
	detect.JPFreq[17][30] = 99
	detect.JPFreq[38][15] = 98
	detect.JPFreq[29][60] = 97
	detect.JPFreq[4][11] = 96
	detect.JPFreq[38][31] = 95
	detect.JPFreq[40][79] = 94
	detect.JPFreq[28][49] = 93
	detect.JPFreq[28][84] = 92
	detect.JPFreq[26][77] = 91
	detect.JPFreq[22][32] = 90
	detect.JPFreq[33][17] = 89
	detect.JPFreq[23][18] = 88
	detect.JPFreq[32][64] = 87
	detect.JPFreq[4][6] = 86
	detect.JPFreq[33][51] = 85
	detect.JPFreq[44][77] = 84
	detect.JPFreq[29][5] = 83
	detect.JPFreq[46][25] = 82
	detect.JPFreq[19][58] = 81
	detect.JPFreq[4][46] = 80
	detect.JPFreq[15][71] = 79
	detect.JPFreq[18][58] = 78
	detect.JPFreq[26][45] = 77
	detect.JPFreq[45][66] = 76
	detect.JPFreq[34][10] = 75
	detect.JPFreq[19][37] = 74
	detect.JPFreq[33][65] = 73
	detect.JPFreq[44][52] = 72
	detect.JPFreq[16][38] = 71
	detect.JPFreq[36][46] = 70
	detect.JPFreq[20][26] = 69
	detect.JPFreq[30][37] = 68
	detect.JPFreq[4][58] = 67
	detect.JPFreq[43][2] = 66
	detect.JPFreq[30][18] = 65
	detect.JPFreq[19][35] = 64
	detect.JPFreq[15][68] = 63
	detect.JPFreq[3][36] = 62
	detect.JPFreq[35][40] = 61
	detect.JPFreq[36][32] = 60
	detect.JPFreq[37][14] = 59
	detect.JPFreq[17][11] = 58
	detect.JPFreq[19][78] = 57
	detect.JPFreq[37][11] = 56
	detect.JPFreq[28][63] = 55
	detect.JPFreq[29][61] = 54
	detect.JPFreq[33][3] = 53
	detect.JPFreq[41][52] = 52
	detect.JPFreq[33][63] = 51
	detect.JPFreq[22][41] = 50
	detect.JPFreq[4][19] = 49
	detect.JPFreq[32][41] = 48
	detect.JPFreq[24][4] = 47
	detect.JPFreq[31][28] = 46
	detect.JPFreq[43][30] = 45
	detect.JPFreq[17][3] = 44
	detect.JPFreq[43][70] = 43
	detect.JPFreq[34][19] = 42
	detect.JPFreq[20][77] = 41
	detect.JPFreq[18][83] = 40
	detect.JPFreq[17][15] = 39
	detect.JPFreq[23][61] = 38
	detect.JPFreq[40][27] = 37
	detect.JPFreq[16][48] = 36
	detect.JPFreq[39][78] = 35
	detect.JPFreq[41][53] = 34
	detect.JPFreq[40][91] = 33
	detect.JPFreq[40][72] = 32
	detect.JPFreq[18][52] = 31
	detect.JPFreq[35][66] = 30
	detect.JPFreq[39][93] = 29
	detect.JPFreq[19][48] = 28
	detect.JPFreq[26][36] = 27
	detect.JPFreq[27][25] = 26
	detect.JPFreq[42][71] = 25
	detect.JPFreq[42][85] = 24
	detect.JPFreq[26][48] = 23
	detect.JPFreq[28][15] = 22
	detect.JPFreq[3][66] = 21
	detect.JPFreq[25][24] = 20
	detect.JPFreq[27][43] = 19
	detect.JPFreq[27][78] = 18
	detect.JPFreq[45][43] = 17
	detect.JPFreq[27][72] = 16
	detect.JPFreq[40][29] = 15
	detect.JPFreq[41][0] = 14
	detect.JPFreq[19][57] = 13
	detect.JPFreq[15][59] = 12
	detect.JPFreq[29][29] = 11
	detect.JPFreq[4][25] = 10
	detect.JPFreq[21][42] = 9
	detect.JPFreq[23][35] = 8
	detect.JPFreq[33][1] = 7
	detect.JPFreq[4][57] = 6
	detect.JPFreq[17][60] = 5
	detect.JPFreq[25][19] = 4
	detect.JPFreq[22][65] = 3
	detect.JPFreq[42][29] = 2
	detect.JPFreq[27][66] = 1
	detect.JPFreq[26][89] = 0
}
