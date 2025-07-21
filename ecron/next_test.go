package ecron

import (
	"fmt"
	"testing"
	"time"
)

func TestCronNext(t *testing.T) {

	next, err := GetNext("0 0 0/1 * * *", time.Now(), 10)

	fmt.Println(next, err)
}
