package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	//以下のポートのサーバーとのコネクションを確立する
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //WithInsecureは本番環境では非推奨
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	//main関数終了時、コネクションがクローズ
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// callListFiles(client)
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownLoadRequest{Filename: "name.txt"}
	stream, err := client.DownLoad(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Response from Download(bytes): %v", res.GetDate())
		log.Printf("Response from Download(string): %v", string(res.GetDate()))
	}
}
