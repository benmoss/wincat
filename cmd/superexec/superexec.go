package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"
)

func main() {
	containerID := os.Args[1]
	if containerID == "" {
		log.Fatal("Usage: superexec.exe container-id")

	}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "gcr.io/cf-london-servces-k8s/windows-images/socat:latest",
		Cmd:          []string{"/wincat.exe", "80"},
		OpenStdin: true,
		AttachStdout: true,
		AttachStdin: true,
		//Tty: true,

	}, &container.HostConfig{
		NetworkMode: container.NetworkMode("container:"+containerID),
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}


	hijack, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stdin: true,
		Stdout: true,
		Stream: true,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		_, e := io.Copy(hijack.Conn, os.Stdin)
		if e != nil {
			panic(e)
		}
	}()
	_, e := stdcopy.StdCopy(os.Stdout, ioutil.Discard, hijack.Reader)
	if e != nil {
		fmt.Printf("wtf %#v\n", e)
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

}
