package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kushsharma/go-grpc-base/proto/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func initServer(conf Config) *cobra.Command {
	thisCmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			grpcAddr := fmt.Sprintf("localhost:%d", conf.ServerPort)
			httpAddr := fmt.Sprintf("localhost:%d", conf.ServerPort+1)

			//create a tcp listener for grpc
			lis, err := net.Listen("tcp", grpcAddr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcOpts := []grpc.ServerOption{}
			grpcServer := grpc.NewServer(grpcOpts...)
			// runtime service instance over gprc
			pb.RegisterRuntimeServiceServer(grpcServer, newRuntimeServiceServer())
			// start grpc server
			go func() {
				log.Fatal(grpcServer.Serve(lis))
			}()

			// prepare http proxy
			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			httpOpts := []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithBlock(),
			}
			mux := runtime.NewServeMux()
			if err := pb.RegisterRuntimeServiceHandlerFromEndpoint(ctx, mux, grpcAddr, httpOpts); err != nil {
				return err
			}

			log.Info("starting grpc server at ", grpcAddr)
			log.Info("starting http proxy at ", httpAddr)
			return http.ListenAndServe(httpAddr, mux)
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

func (sv *runtimeServiceServer) DeploySpecifications(stream pb.RuntimeService_DeploySpecificationsServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Debug("request received for dag: ", in.Dag)

		// process request
		// ..
		// push processing channel
		// ..

		// .. listen for processed channel
		// send ack
		if err := stream.Send(&pb.DeploySpecificationResponse{
			Succcess: true,
			Id:       in.Dag,
		}); err != nil {
			return err
		}
	}
}

func newRuntimeServiceServer() *runtimeServiceServer {
	return &runtimeServiceServer{}
}
