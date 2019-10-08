package xstr

import "regexp"

// 过滤行政区划单位
func FilterRegion (s string) string  {
	reg := regexp.MustCompile(`省|特别行政区|市|自治州|少数民族自治州|自治县|县|镇|乡`)
	return reg.ReplaceAllString(s, "")
}
