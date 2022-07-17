package wiki

import "testing"

func TestFetchItemsPage(t *testing.T) {
	data := fetchItemsPage(map[string]struct{}{
		// "龙门币": {},
	})
	t.Logf("%+v", data)
}
