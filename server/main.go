package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

//ListFilesメソッドの実装
func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {

	fmt.Println("ListFiles was invoked")

	dir := "/Users/fujihara_hideyuki/projects/grpc-lesson/strage"

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

func main() {
	//50051ポートを開く
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//gRPCのサーバー構造体を取得
	s := grpc.NewServer()

	//gRPCサーバーに構造体の内容を登録する
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is running...")
	//変数lisのプロトコル、ポート番号でサーバーを起動
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
