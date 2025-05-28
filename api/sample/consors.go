package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/DavidGamba/go-getoptions"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"os"
	"time"
)

func main() {
	// Parse options
	var exchange string
	var authtoken string
	opt := getoptions.New()
	opt.StringVar(&exchange, "exchange", "TRG")
	opt.StringVar(&authtoken, "token", "foobar")
	_, _ = opt.Parse(os.Args[1:])

	// Setup TLS parameters
	// The certificate has to be CN=localhost but we don't check the chain!
	tlsconfig := tls.Config{ServerName: "localhost", InsecureSkipVerify: true}
	creds := credentials.NewTLS(&tlsconfig)

	// Dial in blocking mode to check and check for errors
	conn, err := grpc.Dial(":40443", grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	// Authenticate to get token
	token := login(conn, authtoken)
	if token == `` {
		log.Fatal("Could not authenticate")
	}

	myaccounts := get_accounts(conn, token)
	var grandtotal float64
	for _, account := range myaccounts {
		grandtotal += get_balance(conn, token, account)
		grandtotal += get_depot_value(conn, token, account, exchange)
	}

	fmt.Println(int64(grandtotal))
	os.Exit(0)
}

func exchange_exists(conn *grpc.ClientConn, token string, id string) bool {
	client := consors_tapi.NewStockExchangeServiceClient(conn)
	atoken := consors_tapi.AccessTokenRequest{AccessToken: token}
	result, err := client.GetStockExchanges(context.Background(), &atoken)
	if err != nil {
		return false
	}

	spew.Dump(result.StockExchangeInfos)
	for _, exchange := range result.StockExchangeInfos {
		if exchange.StockExchange.Id == id {
			return true
		}
	}

	return false
}

func get_balance(conn *grpc.ClientConn, token string, account *consors_tapi.TradingAccount) float64 {
	client := consors_tapi.NewAccountServiceClient(conn)
	tar := consors_tapi.TradingAccountRequest{AccessToken: token,
		TradingAccount: account,
	}
	streamer, err := client.StreamTradingAccount(context.Background(), &tar)
	if err != nil {
		return -1
	}

	for {
		data, _ := streamer.Recv()
		if data.Balance > 0 {
			return data.Balance
			break
		}
	}

	return -1
}

func get_depot_value(conn *grpc.ClientConn, token string, account *consors_tapi.TradingAccount, exchange string) float64 {
	client := consors_tapi.NewDepotServiceClient(conn)
	atoken := consors_tapi.TradingAccountRequest{AccessToken: token,
		TradingAccount: account}
	streamer, err := client.StreamDepot(context.Background(), &atoken)
	if err != nil {
		return -1
	}

	var total float64
	data, err := streamer.Recv()
	for _, n := range data.Entries {
		position, err := get_last_price_parser(conn, token, n.SecurityCode.Code, exchange)
		if err == nil {
			total = total + (position * n.TotalAmount)
		}
		if err != nil {
			log.Printf("Failed to get_last_price for %s: %s\n", n.SecurityCode.Code, err)
		}
	}

	return total
}

func get_accounts(conn *grpc.ClientConn, token string) []*consors_tapi.TradingAccount {
	atoken := consors_tapi.AccessTokenRequest{AccessToken: token}
	client := consors_tapi.NewAccountServiceClient(conn)
	accresponse, err := client.GetTradingAccounts(context.Background(), &atoken)
	if err != nil {
		return nil
	}
	return accresponse.Accounts
}

func login(conn *grpc.ClientConn, secret string) string {
	lr := consors_tapi.LoginRequest{Secret: secret}
	client := consors_tapi.NewAccessServiceClient(conn)
	response, err := client.Login(context.Background(), &lr)
	if err != nil {
		return ``
	}
	return response.AccessToken
}
