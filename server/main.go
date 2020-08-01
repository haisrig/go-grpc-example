package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/haisrig/chatapp/proto"
	"google.golang.org/grpc"
)

type chatbotServer struct{}

var wishes map[string]string

func (c *chatbotServer) AskGenie(context context.Context, q *proto.Question) (*proto.Answer, error) {
	fmt.Println("User: ", q.Question)
	return &proto.Answer{
		Answer: "Hello, How are you?",
	}, nil
}

func (c *chatbotServer) SendGifts(f *proto.Festival, stream proto.SpiritualService_SendGiftsServer) error {
	fmt.Printf("Giene: Got a reuqest for gifts for %s. Sending!!!\n", f.Name)
	for i := 1; i <= 5; i++ {
		gift := proto.Gift{
			Name: "Gift " + strconv.Itoa(i) + " for festival " + f.Name,
		}
		if err := stream.Send(&gift); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (c *chatbotServer) LetsPray(stream proto.SpiritualService_LetsPrayServer) error {
	for {
		if wish, err := stream.Recv(); err == nil {
			dw := strings.ToLower(wish.Name)
			gb := "Not in my list. Next please..."
			for k, v := range wishes {
				if strings.Contains(dw, k) {
					gb = v
				}
			}
			fmt.Println("User:", wish.Name)
			stream.Send(&proto.Blessing{Name: gb})
		} else {
			return err
		}
	}
}

func main() {
	wishes = make(map[string]string)
	wishes["money"] = "You will get it. Don't worry!!!"
	wishes["husband"] = "Most Unluckiest Person. Always lives in dreams!!!"
	wishes["wife"] = "Most Luckiest Person. Can't say anything more :)"
	wishes["happiness"] = "You have to create it on your own!!!"
	wishes["role model"] = "Look around. It can be anyone. It can your Boss, Colleague, Neighbour, Worker ..."
	listener, _ := net.Listen("tcp", "127.0.0.1:8000")
	grpcServer := grpc.NewServer()
	proto.RegisterSpiritualServiceServer(grpcServer, &chatbotServer{})
	fmt.Println("Listening on 8000")
	grpcServer.Serve(listener)
}
