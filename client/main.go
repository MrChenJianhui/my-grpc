package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github/my-grpc/pb/dog"
	"google.golang.org/grpc"
)

func main() {
	l, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("creating client error!")
	}

	/*
		传1出1
		client := dog.NewSearchServiceClient(l)
			res,err := client.Search(context.Background(),&dog.DogReq{Name: "i am client"})
			if err!=nil{
				fmt.Errorf("using search function error!")
			}else{
				fmt.Println(res)
			}
	*/

	/*
		服务端流式
			client := dog.NewSearchServiceClient(l)
			c,err := client.SearchIn(context.Background())
			i:=0
			for{
				if i > 10{
					res,err := c.CloseAndRecv()
					if err != nil{
						fmt.Println(err)
					}else{
						fmt.Println(res)
					}
					break
				}
				i++
				time.Sleep(time.Second)
				c.Send(&dog.DogReq{
					Name: fmt.Sprintf("this is number ",i),
				})
			}
	*/

	/*
		客户端流式
			client:=dog.NewSearchServiceClient(l)
			c,err:=client.SearchOut(context.Background(),&dog.DogReq{Name: "wangwang"})
			for  {
				req,err := c.Recv()
				if err!=nil{
					fmt.Println(err)
					break
				}
				fmt.Println(req)
			}
	*/

	/*
		客户端流式
		client := dog.NewSearchServiceClient(l)
		c, err := client.SearchIO(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		var i int32 = 0
		for {
			if i > 10 {
				break
			}
			c.Send(&dog.DogReq{Name: "wangwang", Age: i})
			time.Sleep(time.Second)
			answer, err := c.Recv()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(answer)
			i++
		}
	*/

	client := dog.NewSearchServiceClient(l)
	c, err := client.SearchIO(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			time.Sleep(time.Second)
			err = c.Send(&dog.DogReq{Name: "wangwang"})
			if err != nil {
				fmt.Println(err)
				wg.Done()
				break
			}
		}
	}()
	go func() {
		for {
			req, err := c.Recv()
			if err != nil {
				wg.Done()
				fmt.Println(err)
			} else {
				fmt.Println(req)
			}
		}
	}()
	wg.Wait()
}
