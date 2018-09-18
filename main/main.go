package main

import (
	"fmt"
	"github.com/johnnywale/mdb/dao"
	"github.com/johnnywale/mdb/feed"
	"github.com/johnnywale/mdb/rest"
)

func main() {
	fdbDao := dao.NewFdbDao()
	feed := feed.NewServer()
	go func() {
		feed.Register(fdbDao)
		feed.Start()
	}()
	fmt.Printf("start http \n")
	rest := rest.NewRestServer(fdbDao)
	rest.Start()

}
