package tests

import (
	"testing"
)

var hardcodeLink = "https://www.avito.ru/js/1/map/items?categoryId=24&locationId=653240&correctorMode=0&page=1&map=eyJzZWFyY2hBcmVhIjp7ImxhdEJvdHRvbSI6NTkuOTAwMTc2NjU1Mjc4MDI0LCJsYXRUb3AiOjU5LjkwODIyNTA1NDMyMzkxNCwibG9uTGVmdCI6MzAuNTE3OTQ0MzQ1NjgwNjQzLCJsb25SaWdodCI6MzAuNTYxMzc0Njc0MDQ5ODF9LCJ6b29tIjoxNX0%3D&params%5B201%5D=1060&params%5B504%5D=5256&params%5B550%5D%5B0%5D=5705&params%5B550%5D%5B1%5D=5704&params%5B550%5D%5B2%5D=5703&params%5B550%5D%5B3%5D=5702&verticalCategoryId=1&rootCategoryId=4&localPriority=0&disabledFilters%5Bids%5D%5B0%5D=byTitle&disabledFilters%5Bslugs%5D%5B0%5D=bt&subscription%5Bvisible%5D=true&subscription%5BisShowSavedTooltip%5D=false&subscription%5BisErrorSaved%5D=false&subscription%5BisAuthenticated%5D=true&searchArea%5BlatBottom%5D=59.900176655278024&searchArea%5BlonLeft%5D=30.517944345680643&searchArea%5BlatTop%5D=59.908225054323914&searchArea%5BlonRight%5D=30.56137467404981&viewPort%5Bwidth%5D=1012&viewPort%5Bheight%5D=374&limit=50&countAndItemsOnly=1"
var hardcodeLink2 = "https://www.avito.ru/js/1/map/items?a=b&c=d"

func TestSome(t *testing.T) {
	for i := range 5 {
		println(i)
	}
}
