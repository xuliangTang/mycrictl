package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
	"log"
	"time"
)

func main() {
	gopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := "unix:///run/containerd/containerd.sock"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, gopts...)
	if err != nil {
		log.Fatalln(err)
	}

	// proto定义: https://github.com/kubernetes/cri-api/blob/master/pkg/apis/runtime/v1/api.proto
	req := &v1.VersionRequest{}
	rsp := &v1.VersionResponse{}
	err = conn.Invoke(ctx, "/runtime.v1.RuntimeService/Version", req, rsp)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rsp)
}
