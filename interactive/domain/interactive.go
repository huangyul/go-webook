package domain

type Interactive struct {
	Id         int64
	BizId      int64
	ReadCnt    int64
	CollectCnt int64
	LikeCnt    int64
	Liked      bool
	Collectd   bool
}
