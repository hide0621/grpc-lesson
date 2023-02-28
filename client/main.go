package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"
	"os"
	"time"

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
	// callDownload(client)
	// CallUpload(client)
	CallUploadAndNotifyProgress(client)
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

func CallUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/fujiwarahideyuki/project/grpc-lesson/strage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}

func CallUploadAndNotifyProgress(client pb.FileServiceClient) {

	filename := "sports.txt"
	path := "/Users/fujiwarahideyuki/project/grpc-lesson/strage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// request
	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalln(sendErr)
			}
			time.Sleep(1 * time.Second)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("received message: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch

}
