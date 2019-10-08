package xarray

import "reflect"

func IsInArray(arr interface{}, item interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		a := reflect.ValueOf(arr)
		for i := 0; i < a.Len(); i++ {
			if reflect.DeepEqual(item, a.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

func IsInIntArray(arr []int, item int) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func IsInStrArray(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func IsInFloatArray(arr []float32, item float32) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func IsInFloat64Array(arr []float64, item float64) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return true
}

func IsInInt32Array(arr []int32, item int32) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func IsInInt64Array(arr []int64, item int64) bool  {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}


