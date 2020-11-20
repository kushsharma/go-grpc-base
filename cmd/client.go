package cmd

import (
	"context"
	"fmt"

	pb "github.com/kushsharma/go-grpc-base/protos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func initClient(conf Config) *cobra.Command {
	thisCmd := &cobra.Command{
		Use: "client",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Info("trying to connect at ", conf.ServerPort)

			var opts []grpc.DialOption
			opts = append(opts, grpc.WithInsecure())
			opts = append(opts, grpc.WithBlock())

			conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", conf.ServerPort), opts...)
			if err != nil {
				return err
			}
			defer conn.Close()

			runtimeServiceClient := pb.NewRuntimeServiceClient(conn)
			versionResponse, err := runtimeServiceClient.Ping(context.Background(), &pb.VersionRequest{
				ClientVersion: "0.0.1a",
			})
			if err != nil {
				return err
			}
			log.Infof("returned %s version from server after ping", versionResponse.ServerVersion)

			return nil
		},
	}
	return thisCmd
}
