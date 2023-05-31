package cmds

import (
	"context"
	"cri/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	v1alpha2 "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"log"
	"os"
	"time"
)

var imageListCmd = &cobra.Command{
	Use:   "images",
	Short: "List images",
	Run: func(cmd *cobra.Command, args []string) {
		client := v1alpha2.NewImageServiceClient(grpcClient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		req := &v1alpha2.ListImagesRequest{}
		rsp, err := client.ListImages(ctx, req)
		if err != nil {
			log.Fatalln(err)
		}

		// 输出表格
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"IMAGE", "TAG", "IMAGE ID", "SIZE"})
		for _, image := range rsp.GetImages() {
			imageName, _ := utils.ParseRepoDigest(image.RepoDigests)
			imageTag := utils.ParseRepoTag(image.RepoTags, imageName)[0][1]
			imageId := utils.ParseImageID(image.Id)
			imageSize := utils.ParseSize(image.Size_)

			row := []string{imageName, imageTag, imageId, imageSize}
			table.Append(row)
		}
		utils.SetTable(table)
		table.Render()
	},
}
