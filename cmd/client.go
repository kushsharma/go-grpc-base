package cmd

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/kushsharma/go-grpc-base/protos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	client_version = "0.0.1a"
)

func initClient(conf Config) *cobra.Command {
	thisCmd := &cobra.Command{
		Use: "client",
	}

	thisCmd.AddCommand(&cobra.Command{
		Use: "ping",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Info("trying to ping at ", conf.ServerPort)

			var conn *grpc.ClientConn
			if conn, err = createConnection(conf.ServerPort); err != nil {
				return err
			}
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			runtimeServiceClient := pb.NewRuntimeServiceClient(conn)
			versionResponse, err := runtimeServiceClient.Ping(ctx, &pb.VersionRequest{
				ClientVersion: client_version,
			})
			if err != nil {
				return err
			}
			log.Infof("returned %s version from server after ping", versionResponse.ServerVersion)
			return nil
		},
	})
	thisCmd.AddCommand(&cobra.Command{
		Use: "deploy",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var conn *grpc.ClientConn
			if conn, err = createConnection(conf.ServerPort); err != nil {
				return err
			}
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			runtimeServiceClient := pb.NewRuntimeServiceClient(conn)
			stream, err := runtimeServiceClient.DeploySpecifications(ctx)
			if err != nil {
				return err
			}

			//emulate specs
			specs := []*pb.DeploySpecificationRequest{
				{
					Dag:   "dag-1",
					Table: "table-1",
					Task: map[string]string{
						"query.sql": "select * from 1",
					},
				},
				{
					Dag:   "dag-2",
					Table: "table-2",
					Task: map[string]string{
						"query.sql": "select * from 2",
					},
				},
			}

			waitc := make(chan struct{})
			go func() {
				for {
					in, err := stream.Recv()
					if err == io.EOF {
						// read done.
						close(waitc)
						return
					}
					if err != nil {
						log.Fatalf("Failed to receive ack : %v", err)
					}
					log.Printf("received ack for: %s", in.Id)
				}
			}()
			for _, spec := range specs {
				if err := stream.Send(spec); err != nil {
					log.Fatalf("Failed to send a spec: %v", err)
				}
			}
			stream.CloseSend()
			<-waitc

			return nil
		},
	})
	return thisCmd
}

func createConnection(port int) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
