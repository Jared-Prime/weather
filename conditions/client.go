package conditions

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"google.golang.org/grpc"
)

type Client struct{}

func (c *Client) Start() error {
	conn, err := grpc.Dial("127.0.0.1:10000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := NewConditionsClient(conn)

	stream, err := client.Report(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		reply, err := stream.CloseAndRecv()
		log.Println(reply)
		log.Println(err)
		conn.Close()
	}()

	for {
		log.Println("sending random weather report")
		err := stream.Send(&Condition{
			Time:        &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Location:    "test",
			Temperature: rand.Float32() * 100,
			Humidity:    rand.Float32() * 100,
		})

		if err != nil {
			return err
		}
	}
}
