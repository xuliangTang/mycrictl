package cmds

import (
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const criAddr = "unix:///run/containerd/containerd.sock"

var grpcClient *grpc.ClientConn
var TTY bool

func initClient() {
	gopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, criAddr, gopts...)
	if err != nil {
		log.Fatalln(err)
	}

	grpcClient = conn
}

func RunCmd() {
	cmd := &cobra.Command{
		Use:          "mycrictl",
		Short:        "my crictl",
		SilenceUsage: true,
	}

	initClient()

	execCmd.Flags().BoolVarP(&TTY, "tty", "t", false, "")
	cmd.AddCommand(versionCmd, imageListCmd, runpCmd, runCmd, psCmd, execCmd)
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
