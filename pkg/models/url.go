package models

import (
	"crypto/md5"
	"encoding/hex"
	"hash/adler32"
)

type Url struct {
	Id           int       `json:"id"`
	Raw          string    `json:"raw"`
	Shortened    string    `json:"shortened"`
	UserId       int       `json:"user_id"`
	RedirectsNum Redirects `json:"redirects"`
}

// NewUrl cretes new url object and shortens url
// it gets md5 hex string and choose every nth characters
// where n is len of hash modulo by adler32 hash sum (to make it shorter)
func NewUrl(raw string, userId int) *Url {
	h := md5.Sum([]byte(raw))
	gap := adler32.Checksum([]byte(raw))%4 + 1
	u := Url{
		Raw:    raw,
		UserId: userId,
	}
	res := make([]byte, 0)
	var i uint32
	sh := hex.EncodeToString(h[:])
	for i = 0; i < uint32(len(h)); i += gap {
		res = append(res, sh[i])
	}
	u.Shortened = string(res)
	return &u
}
