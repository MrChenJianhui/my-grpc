package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github/my-grpc/pb/dog"
	"google.golang.org/grpc"
)

type dogServe struct {
	dog.UnimplementedSearchServiceServer
}

func (*dogServe) Search(ctx context.Context, req *dog.DogReq) (*dog.DogResp, error) {
	msg := req.Name
	resp := &dog.DogResp{Name: msg + "this is server"}
	return resp, nil
}
func (*dogServe) SearchIn(server dog.SearchService_SearchInServer) error {
	for {
		req, err := server.Recv()
		fmt.Println(req)
		if err != nil {
			server.SendAndClose(&dog.DogResp{
				Name: "done",
			})
			break
		}
	}
	return nil
}
func (*dogServe) SearchOut(req *dog.DogReq, server dog.SearchService_SearchOutServer) error {
	i := 0
	for {
		if i > 10 {
			break
		}
		time.Sleep(time.Second)
		server.Send(&dog.DogResp{
			Name: req.Name,
		})
		i++
	}
	return nil
}

func (*dogServe) SearchIO(server dog.SearchService_SearchIOServer) error {
	ch := make(chan string, 0)
	for {
		req, err := server.Recv()
		if err != nil {
			ch <- "over"
			break
		}
		ch <- req.Name
	}
	for {
		s := <-ch
		if s == "over" {
			break
		}
		server.Send(&dog.DogResp{
			Name: s,
		})
	}
	return nil
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Errorf("creating listener error!")
	}
	s := grpc.NewServer()
	dog.RegisterSearchServiceServer(s, &dogServe{})
	s.Serve(l)
}
