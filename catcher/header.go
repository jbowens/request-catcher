package catcher

type Header struct {
	ID        int64  `json:"-" db:"id"`
	RequestID int64  `json:"-" db:"request_id"`
	Key       string `json:"key" db:"key"`
	Value     string `json:"value" db:"value"`
}
