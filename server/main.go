package main

import (
	"bytes"
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

// ListFilesメソッドの実装
func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {

	fmt.Println("ListFiles was invoked")

	dir := "/Users/fujiwarahideyuki/project/grpc-lesson/strage"

	//変数dirのパスを取得して代入
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	//ファイル名が格納されるスライス
	var filenames []string

	for _, path := range paths {
		//変数pathにあるのがファイルであれば
		if !path.IsDir() {
			//スライスにファイル名を追加
			filenames = append(filenames, path.Name())
		}
	}

	//戻り値のメッセージ
	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}

	return res, nil

}

func (*server) Download(req *pb.DownLoadRequest, stream pb.FileService_DownLoadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := "/Users/fujiwarahideyuki/project/grpc-lesson/strage/" + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		res := &pb.DownLoadResponse{Date: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data(bytes): %v", data)
		log.Printf("received data(string): %v", string(data))
		buf.Write(data)
	}

}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data: %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %vbytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

func myLogging() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("request data: %+v", req)

		resp, err = handler(ctx, req) // handler関数は実際はクライアントが呼んだ関数が実行される
		if err != nil {
			return nil, err
		}
		log.Printf("response data: %+v", resp)

		return resp, nil
	}
}

func main() {
	//50051ポートを開く
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//gRPCのサーバー構造体を取得
	//加えてmyLogging関数をインターセプターに追加
	// s := grpc.NewServer()
	s := grpc.NewServer(grpc.UnaryInterceptor(myLogging()))

	//gRPCサーバーに構造体の内容を登録する
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is running...")
	//変数lisのプロトコル、ポート番号でサーバーを起動
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
