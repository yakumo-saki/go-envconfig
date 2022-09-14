package envconfig

import "fmt"

// buildSliceFromValueMap make slice[string] from config suffix by _00 ~ _99
func buildSliceFromValueMap(prefix string, valueMap map[string]string) []string {
	ret := make([]string, 0)
	for i := 0; i < 100; i++ {
		v, ok := valueMap[prefix+fmt.Sprint(i)]
		if ok {
			ret = append(ret, v)
		}

		// 0〜9 の場合 00~09 もチェック
		if i < 10 {
			v, ok := valueMap[fmt.Sprintf("%s%02d", prefix, i)]
			if ok {
				ret = append(ret, v)
			}
		}
	}
	return ret
}
