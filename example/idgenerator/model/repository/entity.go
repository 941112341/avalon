package repository

type IdGenerator struct {
	ID      int64
	MaxID   int64
	Length  int64
	BizID   string
	Version int
}
