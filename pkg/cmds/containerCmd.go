package cmds

import (
	"context"
	"cri/pkg/utils"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
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

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		client := v1alpha2.NewRuntimeServiceClient(grpcClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		req := &v1alpha2.ListContainersRequest{}
		rsp, err := client.ListContainers(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}

		// 输出表格
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"CONTAINER ID", "IMAGE", "CREATED", "STATE", "NAME", "ATTEMPT", "POD ID"})
		for _, container := range rsp.GetContainers() {
			containerId := container.Id[:13]
			image := container.Image.GetImage()
			created := utils.NsToTime(container.CreatedAt)
			state := strings.Replace(container.State.String(), "CONTAINER_", "", -1)
			name := container.Metadata.Name
			attempt := container.Metadata.Attempt
			podId := container.PodSandboxId[:13]

			row := []string{containerId, image, created, state, name, strconv.Itoa(int(attempt)), podId}
			table.Append(row)
		}
		utils.SetTable(table)
		table.Render()
	},
}

var execCmd = &cobra.Command{
	Use:     "exec",
	Short:   "Run a command in a running container",
	Example: "exec containerId sh -t",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatalln("参数错误")
		}

		containerId, command := args[0], args[1:]

		client := v1alpha2.NewRuntimeServiceClient(grpcClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		req := &v1alpha2.ExecRequest{
			ContainerId: containerId,
			Cmd:         command,
			Tty:         TTY,
			Stdin:       true,
			Stdout:      true,
			Stderr:      !TTY,
		}
		rsp, err := client.Exec(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}

		URL, err := url.Parse(rsp.Url)
		if err != nil {
			log.Fatalln(err)
		}
		exec, err := remotecommand.NewSPDYExecutor(&rest.Config{
			TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		}, "POST", URL)
		if err != nil {
			log.Fatalln(err)
		}

		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Tty:    TTY,
		})
		if err != nil {
			log.Fatalln(err)
		}
	},
}
