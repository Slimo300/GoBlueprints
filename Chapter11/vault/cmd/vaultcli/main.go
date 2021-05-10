package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/Slimo300/GoBlueprints/Chapter10/vault/pb"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081", "gRPC address")
	)
	flag.Parse()

	ctx := context.Background()

	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()

	//vaultService := grpcclient.New(conn)
	client := pb.NewVaultClient(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)

	switch cmd {
	case "hash":
		var password string
		password, args = pop(args)
		hash(ctx, client, password)
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		validate(ctx, client, password, hash)
	default:
		log.Fatalln("unknown command", cmd)
	}

}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}

	return s[0], s[1:]
}

func hash(ctx context.Context, client pb.VaultClient, password string) {
	h, err := client.Hash(ctx, &pb.HashRequest{Password: password})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(h)
}

func validate(ctx context.Context, client pb.VaultClient, password, hash string) {
	valid, err := client.Validate(ctx, &pb.ValidateRequest{Password: password, Hash: hash})
	if err != nil {
		log.Fatalln(err.Error())
	}

	if !valid.Valid {
		fmt.Println("invalid")
		os.Exit(1)
	}
	fmt.Println("valid")
}
