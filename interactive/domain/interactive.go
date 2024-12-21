package domain

type Interactive struct {
	Biz        string
	BizID      int64
	ReadCnt    int
	LikeCnt    int
	CollectCnt int
	Liked      bool
	Collected  bool
}
