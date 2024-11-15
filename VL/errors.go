package vl

const (
	ErrBottomRecordDoesNotExist = "bottom record does not exist"
	ErrUidIsEmpty               = "UID is empty" //TODO
	ErrDataIsEmpty              = "data is empty"
	ErrRecordIsNotFound         = `record is not found, uid=%v`
	ErrRecordIsOutdated         = `record is outdated, uid=%v`
	ErrRecordIsTooBig           = "record is too big"
	ErrTtlIsZero                = "zero TTL will totally disable the cache"
)
