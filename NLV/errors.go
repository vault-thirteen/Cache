package nvl

const (
	ErrBottomRecordDoesNotExist = "bottom record does not exist"
	ErrUidIsEmpty               = "UID is empty" //TODO
	ErrRecordIsNotFound         = `record is not found, uid=%v`
	ErrRecordIsOutdated         = `record is outdated, uid=%v`
	ErrTtlIsZero                = "zero TTL will totally disable the cache"
)
