package cmds

import (
	"context"
	"cri/pkg/utils"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"log"
	"time"
)

var runpCmd = &cobra.Command{
	Use:     "runp",
	Short:   "Run a new pod",
	Example: "runp ./sandbox.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("请指定pod-config.yaml")
		}

		client := v1alpha2.NewRuntimeServiceClient(grpcClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		config := &v1alpha2.PodSandboxConfig{}
		if err := utils.YamlFile2Struct(args[0], config); err != nil {
			log.Fatalln(err)
		}
		req := &v1alpha2.RunPodSandboxRequest{Config: config}

		rsp, err := client.RunPodSandbox(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(rsp.PodSandboxId)
	},
}
