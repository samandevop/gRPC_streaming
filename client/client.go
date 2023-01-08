package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "app/calculatorpb"
)

func main() {

	conn, err := grpc.Dial("localhost:9002", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDirectionalStream(c)
}

func doUnary(c pb.CalculatorServiceClient) {

	fmt.Println("Starting to do SquareRoot Unary RPC...")

	req := pb.SquareRequest{
		Number: 9,
	}

	res, err := c.SquareRoot(context.Background(), &req)
	if err != nil {
		log.Println("error while calling SquareRoot RPC:", err)
	}

	fmt.Println("Reponse from SquareRoot:", res.GetSqrtResult())
}

func doServerStreaming(c pb.CalculatorServiceClient) {

	fmt.Println("Starting to do PerfectNumber Server Streamin RPC...")

	stream, err := c.PerfectNumber(context.Background(), &pb.PerfectNumberRequest{
		Number: 1000,
	})
	if err != nil {
		log.Println("error while calling PerfectNumber RPC:", err)
	}

	for {

		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("error while calling PerfectNumber Stream Recv RPC:", err)
		}

		fmt.Println(resp)
	}
}

func doClientStreaming(c pb.CalculatorServiceClient) {

	stream, err := c.TotalNumber(context.Background())
	if err != nil {
		log.Println("error while calling TotalNumber RPC:", err)
	}

	nums := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10.1}

	for _, num := range nums {
		err = stream.Send(&pb.TotalNumberRequest{
			Number: num,
		})
		if err != nil {
			log.Println("error while calling TotalNumber Client Stream Send RPC:", err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("error while calling TotalNumber Client Stream Send RPC:", err)
		return
	}

	fmt.Println(res)

}

func doBiDirectionalStream(c pb.CalculatorServiceClient) {

	stream, err := c.FindMinimum(context.Background())
	if err != nil {
		log.Println("error while calling FindMinimum Send RPC:", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{4, 6, 2, 23, 13, 7, 8, 53, 1}

		for _, num := range numbers {
			stream.Send(&pb.FindMinimumRequest{Number: num})
			time.Sleep(time.Second * 1)
		}

		stream.CloseSend()
	}()

	go func() {
		for {

			res, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("error while calling FindMinimum Recv RPC:", err)
				return
			}

			fmt.Println(res)
		}

		close(waitc)
	}()

	<-waitc
}
