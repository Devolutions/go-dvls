package dvls

const listPageSize int = 100

type pagedResponse struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalCount  int `json:"totalCount"`
	TotalPage   int `json:"totalPage"`
}

func fetchAllPages(fetchPage func(pageNumber int) (pagedResponse, int, error)) error {
	fetched := 0

	for pageNumber := 1; ; pageNumber++ {
		meta, count, err := fetchPage(pageNumber)
		if err != nil {
			return err
		}

		if count == 0 {
			return nil
		}
		fetched += count

		if meta.TotalCount > 0 && fetched >= meta.TotalCount {
			return nil
		}
		if meta.TotalPage > 0 && pageNumber >= meta.TotalPage {
			return nil
		}
	}
}
