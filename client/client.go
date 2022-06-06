package main

import (
	"log"
	"time"
	"context"
	pb "mailinglist/proto"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponse(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Fatalf("	error: %v", err)
	}

	if res.EmailEntry == nil {
		log.Printf("	email not found")
	} else {
		log.Printf("	response: %v", res.EmailEntry)
	}
}

func createEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("create email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateEmail(ctx &pb.CreateEmailRequest{EmailEntry: addr})
	logResponse(res, err)

	return res.EmailEntry
}

func getEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("get email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmail(ctx &pb.GetEmailRequest{EmailEntry: addr})
	logResponse(res, err)

	return res.EmailEntry
}

func getEmailBatch(client pb.MailingListServiceClient, count int, page int) {
	log.Println("get email batch")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmailBatch(ctx &pb.GetEmailBatchRequest{Count: int32(count), Page: int32(page)})
	if err != nil {
		log.Fatalf("	error: %v", err)
	}
	log.Println("response:")
	for i := 0 ; i < len(res.EmailEntries) ; i++ {
		log.Printf("	item [%v of %v]: %s", i+1, len(res.EmailEntries), res.EmailEntries[i])
	}
}

func updateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) *pb.EmailEntry {
	log.Println("update email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateEmail(ctx &pb.UpdateEmailRequest{EmailEntry: &entry})
	logResponse(res, err)

	return res.EmailEntry
}

func deleteEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("delete email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteEmail(ctx &pb.DeleteEmailRequest{EmailEntry: addr})
	logResponse(res, err)

	return res.EmailEntry
}

var args struct {
	GrpcAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main() {
	arg.MustParse(&args)

	if args.GrpcAddr == "" {
		agrs.GrpcAddr = ":8081"
	}

	conn, err := grpc.Dial(args.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMailingListServiceClient(conn)

	newEmail := createEmail(client, "taylor@taylor.edu")
	newEmail.ConfirmedAt = 10000
	updateEmail(client, *newEmail)
	deleteEmail(client, newEmail.Email)
	getEmailBatch(client, 5, 1)
}