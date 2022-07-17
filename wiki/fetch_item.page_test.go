package wiki

import (
	"fmt"
	"testing"
	"time"
)

func TestFetchItemInfo(t *testing.T) {
	res := FetchItemInfo("理智", "先锋芯片")
	for name, v := range res {
		fmt.Printf("%v:\n\t%v\n\t%v\n\t%v\n",
			name,
			v["image"].(string),
			v["wikishort"].(string),
			v["updatedat"].(time.Time).Format("2006-01-02 15:04:05"),
		)
	}
}
