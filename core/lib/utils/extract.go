package utils

func BatchIDs(ids []int64, batchSize int) [][]int64 {

	batchCount := len(ids) / batchSize
	batchRemainder := len(ids) % batchSize
	if batchRemainder > 0 {
		batchCount++
	}

	batches := make([][]int64, batchCount)

	for k := range batches {

		batchLen := batchSize

		// Is Last?
		if k == batchCount-1 && batchRemainder > 0 {
			batchLen = batchRemainder
		}

		batches[k] = make([]int64, batchLen)

		l := 0
		indexOfID := k * batchSize
		finish := indexOfID + batchLen

		for {
			if indexOfID == finish {
				break
			}

			batches[k][l] = ids[indexOfID]

			l++
			indexOfID++
		}
	}

	return batches
}

func ExtractIDs(objLen int, fn func(index int) int64) []int64 {

	ids := []int64{}
	gate := map[int64]struct{}{}
	for k := 0; k < objLen; k++ {
		if _, ok := gate[fn(k)]; !ok {
			gate[fn(k)] = struct{}{}
			ids = append(ids, fn(k))
		}
	}

	return ids
}

func SumFloat(objLen int, fn func(index int) float64) float64 {
	var total float64 = 0
	for k := 0; k < objLen; k++ {
		total += fn(k)
	}
	return total
}

func SumInt64(objLen int, fn func(index int) int64) int64 {
	var total int64 = 0
	for k := 0; k < objLen; k++ {
		total += fn(k)
	}
	return total
}
