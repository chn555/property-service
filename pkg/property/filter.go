package property

type SortOrder int8

const (
	Ascending SortOrder = iota
	Descending
)

type AmountType int8

const (
	All AmountType = iota
	Expense
	Income
)
