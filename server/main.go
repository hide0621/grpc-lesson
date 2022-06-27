package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io/ioutil"
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
