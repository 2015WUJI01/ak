package wiki

// 龙门币   | item | Wiki | 未找到
// 龙门币   | item | WikiShort | 找到，但值为空
const tml = "| %-12s | %-8s | %-10s |"

const (
	ItemNotFound   = "item | %-10s | %-10s | 未找到" // ItemNotFound xxx 未找到
	ItemFoundEmpty = "item | %-10s | %-10s | 值为空" // ItemFoundEmpty xxx 找到了，但值为空
	ItemFounded    = "item | %-10s | %-10s | "    // ItemNotFound xxx 未找到

	OprNotFound   = " opr | %-10s | %-10s | 未找到" // ItemNotFound xxx 未找到
	OprFoundEmpty = " opr | %-10s | %-10s | 值为空" // ItemFoundEmpty xxx 找到了，但值为空
)
