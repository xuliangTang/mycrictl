package cmds

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1alpha2 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"log"
	"time"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display runtime version information",
	Run: func(cmd *cobra.Command, args []string) {
		_ = v1alpha2.PodSandboxConfig{}
		client := v1alpha2.NewRuntimeServiceClient(grpcClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		req := &v1alpha2.VersionRequest{}
		rsp, err := client.Version(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Version: ", rsp.Version)
		fmt.Println("RuntimeName: ", rsp.RuntimeName)
		fmt.Println("RuntimeVersion: ", rsp.RuntimeVersion)
		fmt.Println("RuntimeApiVersion: ", rsp.RuntimeApiVersion)
	},
}
