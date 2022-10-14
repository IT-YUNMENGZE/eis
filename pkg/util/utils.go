package util

import "strconv"

func StrToInt64(str string) int64 {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return int64(i)
	}
}

func Uint64ToInt64(intNum uint64) int64 {
	//uint64 ==> string ==> int64
	return StrToInt64(strconv.FormatUint(intNum, 10))
}

func Split(s string) string {
    var n []rune
    for _, r := range s {
		if r >= '0' && r <= '9' || r == '.' {
			n = append(n, r)
		} else {
			break
		}
    }
    return string(n)
}