// The client demonstrates how to supply an OAuth2 token for every RPC.
package main

import (
	"crypto/tls"
	"github.com/johnnywale/mdb/feed"
	"github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"time"
)

func main() {
	perRPC := oauth.NewOauthAccess(fetchToken())
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	}
	conn, err := grpc.Dial(":8081", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := feed.NewProtoServiceClient(conn).ManyToOne(ctx)
	for i := 0; i < 100000; i++ {
		point := &feed.Request{}
		tx := &feed.Transaction{}
		tx.Id = 10000 + int64(i)
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
		point.Path = []string{"aaaa", "ddddd"}
		bs, _ := proto.Marshal(tx)
		point.Protobuf = bs
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
}

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "eyJhbGciOiJIUzI1NiJ9.eyJ1aWQiOjEsImF0IjozLCJ1c2VyX25hbWUiOiJ1c2VybmFtZSIsImN0eCI6LTEsInNjb3BlIjpbInR4OncsdHg6ciJdLCJwaWQiOi0xLCJhaWQiOjEsInVyIjozLCJhbiI6ImFjY291bnRfbmFtZSIsImp0aSI6IjRlMTg0MDE2LWIxZTktNGNiZC05YmRjLTE3OWI4NWQzODAzOCIsImNsaWVudF9pZCI6ImNsaWVudF9pZCIsImFwIjoiMSJ9.OQyOkrUQ5RwGgl_9-12QT5imLbtMRrlRs8rDZQAnqaA",
	}
}
