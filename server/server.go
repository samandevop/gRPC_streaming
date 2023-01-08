package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "app/calculatorpb"
)

type server struct {
	*pb.UnimplementedCalculatorServiceServer
}

func (s *server) SquareRoot(ctx context.Context, req *pb.SquareRequest) (*pb.SquareResponse, error) {

	var sr float32 = float32(req.Number) / 2
	var temp float32
	for {
		temp = sr
		sr = (temp + (float32(req.Number) / temp)) / 2
		if (temp - sr) == 0 {
			break
		}
	}
	return &pb.SquareResponse{
		SqrtResult: sr,
	}, nil
}

func (s *server) PerfectNumber(req *pb.PerfectNumberRequest, stream pb.CalculatorService_PerfectNumberServer) error {

	fmt.Println("Req:", req)

	number := req.GetNumber()
	nums := int64(1)
	for nums <= number {
		if findPerfectNumber(nums) {
			stream.Send(&pb.PerfectNumberResponse{
				PerfectNumber: nums,
			})
			nums++
		} else {
			nums++
		}
	}

	return nil
}

func (s *server) TotalNumber(stream pb.CalculatorService_TotalNumberServer) error {

	var (
		total float64
	)

	for {

		req, err := stream.Recv()

		if err == io.EOF {

			err = stream.SendAndClose(&pb.TotalNumberResponse{
				TotalNumber: total,
			})

			if err != nil {
				log.Println("error while TotalNumber Recv:", err)
				return err
			}

			return nil
		}

		if err != nil {
			log.Println("error while TotalNumber Recv:", err)
			return err
		}

		total += req.Number
	}
}

func (s *server) FindMinimum(stream pb.CalculatorService_FindMinimumServer) error {
	boolean := true
	minimum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if boolean {
			minimum = req.Number
			boolean = false
			err = stream.Send(&pb.FindMinimumResponse{Minimum: minimum})
			if err != nil {
				log.Println("error while FindMinimum Send:", err)
				return err
			}
		}

		if err != nil {
			log.Println("error while FindMinimum Recv:", err)
			return err
		}

		fmt.Println(req)

		if req.Number < minimum {
			minimum = req.Number
			err = stream.Send(&pb.FindMinimumResponse{Minimum: minimum})
			if err != nil {
				log.Println("error while FindMinimum Send:", err)
				return err
			}
		}
	}
}

func findPerfectNumber(n int64) bool {
	var total int64
	for i := int64(1); i < n; i++ {
		if n%i == 0 {
			total += i
		}
	}

	if total == n {
		return true
	} else {
		return false
	}
}

func main() {

	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Listening :9002...")
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}
