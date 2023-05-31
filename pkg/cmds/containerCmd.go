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

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run a new container inside a running sandbox",
	Example: "run sandboxId ./container.yaml ./sandbox.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			log.Fatalln("参数错误")
		}

		podSandboxId, containerConfigPath, sandboxConfigPath := args[0], args[1], args[2]

		// 解析容器配置对象
		containerConfig := &v1alpha2.ContainerConfig{}
		if err := utils.YamlFile2Struct(containerConfigPath, containerConfig); err != nil {
			log.Fatalln(err)
		}

		// 解析pod sandbox配置对象
		podSandboxConfig := &v1alpha2.PodSandboxConfig{}
		if err := utils.YamlFile2Struct(sandboxConfigPath, podSandboxConfig); err != nil {
			log.Fatalln(err)
		}

		client := v1alpha2.NewRuntimeServiceClient(grpcClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// 创建容器
		createReq := &v1alpha2.CreateContainerRequest{
			PodSandboxId:  podSandboxId,
			Config:        containerConfig,
			SandboxConfig: podSandboxConfig,
		}
		createRsp, err := client.CreateContainer(ctx, createReq)
		if err != nil {
			log.Fatalln(err)
		}

		// 启动容器
		startReq := &v1alpha2.StartContainerRequest{
			ContainerId: createRsp.ContainerId,
		}
		_, err = client.StartContainer(ctx, startReq)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(createRsp.ContainerId)
	},
}
