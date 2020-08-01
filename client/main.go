package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/haisrig/chatapp/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client started....")
	conn, _ := grpc.Dial("localhost:8000", grpc.WithInsecure())
	client := proto.NewSpiritualServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var input string
	fmt.Println("Please enter communication type[unary|one2many|binary]:")
	fmt.Scanln(&input)
	if strings.EqualFold(input, "unary") {
		fmt.Println("Say something to Genie...")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		response, _ := client.AskGenie(ctx, &proto.Question{Question: text})
		fmt.Println("Genie:", response.Answer)
	} else if strings.EqualFold(input, "one2many") {
		gClient, _ := client.SendGifts(ctx, &proto.Festival{Name: "Diwali"})
		for {
			gift, err := gClient.Recv()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Genie: ", gift)
		}
	} else if strings.EqualFold(input, "binary") {
		fmt.Println("Ask Genie if you have any questions???")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		wClient, _ := client.LetsPray(ctx)
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			wClient.Send(&proto.Wish{Name: text})
			wish, err := wClient.Recv()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("Genie: ", wish.Name)
		}
	}
}
