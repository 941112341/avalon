package inline

/*
	if seed > target return seed
*/
func FirstLargerInt(seed, target int) int {
	i := seed
	for i < target {
		i = i << 1
		if i < seed {
			return target
		}
	}
	return i
}

//(maxInt - length, maxInt]
func BuildIntList(maxInt, length int64) []int64 {
	list := make([]int64, 0)
	for i := maxInt - length + 1; i <= maxInt; i++ {
		list = append(list, i)
	}
	return list
}

func LastInt64(list []int64) *int64 {
	var p *int64
	l := len(list) - 1
	if l > -1 {
		p = &list[l]
	}
	return p
}
