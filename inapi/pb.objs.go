// Code generated by github.com/hooto/protobuf_slice
// DO NOT EDIT!

package inapi

func PbStringSliceEqual(ls, ls2 []string) bool {
	if len(ls) != len(ls2) {
		return false
	}
	for _, v := range ls {
		hit := false
		for _, v2 := range ls2 {
			if v == v2 {
				hit = true
				break
			}
		}
		if !hit {
			return false
		}
	}
	return true
}

func PbStringSliceSyncSlice(ls, ls2 []string) ([]string, bool) {
	if len(ls2) == 0 {
		return ls, false
	}
	hit := false
	changed := false
	for _, v2 := range ls2 {
		hit = false
		for _, v := range ls {
			if v == v2 {
				hit = true
				break
			}
		}
		if !hit {
			ls = append(ls, v2)
			changed = true
		}
	}
	return ls, changed
}

func PbInt32SliceEqual(ls, ls2 []int32) bool {
	if len(ls) != len(ls2) {
		return false
	}
	for _, v := range ls {
		hit := false
		for _, v2 := range ls2 {
			if v == v2 {
				hit = true
				break
			}
		}
		if !hit {
			return false
		}
	}
	return true
}

func PbInt32SliceSyncSlice(ls, ls2 []int32) ([]int32, bool) {
	if len(ls2) == 0 {
		return ls, false
	}
	hit := false
	changed := false
	for _, v2 := range ls2 {
		hit = false
		for _, v := range ls {
			if v == v2 {
				hit = true
				break
			}
		}
		if !hit {
			ls = append(ls, v2)
			changed = true
		}
	}
	return ls, changed
}
