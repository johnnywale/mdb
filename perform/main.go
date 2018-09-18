package main

import (
	"fmt"
	"github.com/johnnywale/mdb/dao"
	"github.com/johnnywale/mdb/feed"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"sync"
	"time"
)

func test(c int64, dao *dao.FdbDao, wg *sync.WaitGroup) {
	var function = func(tr fdb.Transaction, directory directory.DirectorySubspace) {
		for i := 0; i < 10000; i++ {
			tx := &feed.Transaction{}
			tx.Id = int64(c*10000) + int64(i)
			tx.AccountId = 22121
			tx.AccountSeqIds = []int64{100, 200}
			tx.Amount = &feed.Money{Amount: 2000, CurrencyUnit: "usd"}
			tx.ApplicationId = 200
			tx.Balance = &feed.Money{Amount: 2000, CurrencyUnit: "usd"}
			tx.BalanceType = feed.BalanceType_CREDIT_BALANCE
			tx.Browser = feed.Browser_CHROME
			tx.BrowserVer = 12
			tx.CampaignId = 21323
			tx.Category = feed.TransactionCategory_BUY
			tx.City = "CN"
			tx.CountryCode = "us"
			tx.CountryIpCode = "us"
			tx.Created = &timestamp.Timestamp{Seconds: 12, Nanos: 212}
			tx.CreatedBy = 12
			tx.CurrencyUnit = "usd"
			tx.DeviceType = feed.DeviceType_MOBILE
			tx.DeviceTypeVersion = "af"
			tx.ExchangeRateBatchId = 120
			tx.ExternalRef = "uuabi=rcqfgaolplm"
			tx.Ip = "192.168.11.11"
			tx.Isp = "cn"
			tx.ItemId = 123
			tx.LoginContextId = 123
			tx.LoginTime = &timestamp.Timestamp{Seconds: 12, Nanos: 212}
			tx.MetaData = "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
			tx.Os = feed.OperatingSystem_BSD
			tx.OsVersion = "fab"
			tx.RedirectTime = &timestamp.Timestamp{Seconds: 12, Nanos: 212}
			tx.Session = "sessionid"
			tx.Test = false
			tx.TransactionTime = &timestamp.Timestamp{Seconds: 12, Nanos: 212}
			tx.TxId = "ac_Tx202192123v"
			tx.Type = feed.TransactionType_DEBIT
			tx.VendorId = 122
			tx.WalletCode = "fff"
			tx.WalletId = 12432
			bs, _ := proto.Marshal(tx)
			key := directory.Pack(tuple.Tuple{tx.Id})
			tr.Set(key, bs)
		}
	}
	dao.Execute([]string{"abc", "efg"}, function)
	wg.Done()
}

func main() {
	var start = time.Now()
	dao := dao.NewFdbDao()
	if false {
		dao.Clean([]string{"abc", "efg"})
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go test(int64(i), dao, &wg)
		}
		wg.Wait()
		fmt.Printf("write 100000 %s \n", time.Now().Sub(start))
		start = time.Now()
	} else {
		fmt.Printf("skip insert")
	}
	dao.Load([]string{"abc", "efg"})
	fmt.Printf("read 100000  in %s  \n", time.Now().Sub(start))
}
