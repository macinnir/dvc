package utils

func Paged(limit, page int64) (offset int64) {

	if page == 0 {
		return 0
	}

	return page * limit

}
