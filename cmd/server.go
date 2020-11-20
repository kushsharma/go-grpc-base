package cmd

import (
	"context"
	"fmt"
	"net"

	pb "github.com/kushsharma/go-grpc-base/protos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func initServer(conf Config) *cobra.Command {
	thisCmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", conf.ServerPort))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			var opts []grpc.ServerOption
			grpcServer := grpc.NewServer(opts...)
			pb.RegisterRuntimeServiceServer(grpcServer, newRuntimeServiceServer())

			log.Info("starting server at ", conf.ServerPort)
			return grpcServer.Serve(lis)
		},
	}
	return thisCmd
}

type runtimeServiceServer struct {
	pb.UnimplementedRuntimeServiceServer // https://github.com/grpc/grpc-go/issues/3794
}

func (sv *runtimeServiceServer) Ping(ctx context.Context, version *pb.VersionRequest) (*pb.VersionResponse, error) {
	log.Debugf("client with version %s requested for ping", version.ClientVersion)
	response := &pb.VersionResponse{
		ServerVersion: "1.0.0",
	}
	return response, nil
}

func newRuntimeServiceServer() *runtimeServiceServer {
	return &runtimeServiceServer{}
}
