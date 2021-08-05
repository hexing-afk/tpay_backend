package utils

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {

}

func TestGetUniqueId(t *testing.T) {
	t.Logf("id: %v", GetUniqueId())
}

func TestRandString(t *testing.T) {
	for i := 1; i <= 10; i++ {
		t.Logf("str: %v", RandString(32))
		t.Logf("len: %v", len(RandString(32)))
	}

}

func TestRemoveRepeatedElement(t *testing.T) {
	n := []int64{1, 3, 4, 4, 2, 2, 5, 21, 3, 12, 3, 87, 123, 123, 1, 3, 5, 1, 123, 123}
	//a := DistinctSliceOnIt64(n)
	//delRepeatElem(n)
	a := removeDuplication_map(n)
	t.Logf("a：%v", a)
}

func delRepeatElem(nums []int64) int {
	fmt.Println(nums, &nums[0])
	for i := 0; i < len(nums)-1; i++ {
		if nums[i]^nums[i+1] == 0 { //重复元素执行异或操作等于0.
			nums = append(nums[:i], nums[i+1:]...) //删除重复元素
		}
	}
	fmt.Println(nums, &nums[0])
	return len(nums)
}

func TestSort(t *testing.T) {
	var sli = []int64{4, 3, 3, 15, 131, 1, 5, 3, 6, 4, 57, 9, 31, 23, 25}

	length := len(sli)
	var minIdx int
	var temp int64
	for i := 0; i < length-1; i++ {
		minIdx = i
		for j := i + 1; j < length; j++ {
			if sli[j] < sli[minIdx] {
				minIdx = j
			}
		}

		temp = sli[i]
		sli[i] = sli[minIdx]
		sli[minIdx] = temp
	}

	t.Logf("sli: %v", sli)
}

func removeDuplication_map(arr []int64) []int64 {
	set := make(map[int64]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}
