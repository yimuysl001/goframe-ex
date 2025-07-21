package eid

import (
	"github.com/oklog/ulid/v2"
)

func GetUlid() string {
	//entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	//ms := ulid.Timestamp(time.Now())
	//return ulid.MustNew(ms, entropy).String()
	return ulid.Make().String()
}

func ParseUlid(ulidstr string) (id ulid.ULID, err error) {

	return ulid.Parse(ulidstr)

}
