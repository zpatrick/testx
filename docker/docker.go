package docker

import "context"

// https://stackoverflow.com/questions/45429276/how-to-run-docker-run-using-go-sdk-for-docker

// cli, err := client.NewEnvClient()
//     if err != nil {
//         panic(err)
//     }

//     ctx := context.Background()
//     resp, err := cli.ContainerCreate(ctx, &container.Config{
//         Image:        "mongo",
//         ExposedPorts: nat.PortSet{"8080": struct{}{}},
//     }, &container.HostConfig{
//         PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("8080"): {{HostIP: "127.0.0.1", HostPort: "8080"}}},
//     }, nil, "mongo-go-cli")
//     if err != nil {
//         panic(err)
//     }

//     if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
//         panic(err)
//     }

type Container struct {
}

type ContainerConfig struct {
	Image       ImageConfig
	Ports       map[int]int
	Environment map[string]string
}

type ImageConfig struct {
	Name string
	Tag  string
}

type PortConfig struct {
	Inside  int
	Outside int
}

func NewContainer(cfg ContainerConfig) *Container {

	return &Container{}
}

func (c *Container) Start(ctx context.Context) error {
	return nil
}

func (c *Container) Stop(ctx context.Context) error {
	return nil
}

func (c *Container) IsRunning(ctx context.Context) bool {
	return false
}
