package main

import (
	"fmt"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

func main() {
	newUlid := xulid.Instance().NewULID()
	fmt.Println(newUlid)
	fmt.Println(newUlid)
}
