package main

import (
	"fmt"
	"os"

	"github.com/closetool/mydocker/cgroups/subsystems"
	"github.com/closetool/mydocker/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name: "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:"cpushare",
			Usage:"cpushare limit",
		},
		cli.StringFlag{
			Name:"cpuset",
			Usage:"cpuset limit",
		},
		cli.StringFlag{
			Name:"v",
			Usage: "volume",
		},
		cli.StringFlag{
			Name: "name",
			Usage: "container name",
		},
		cli.StringSliceFlag{
			Name:"e",
			Usage: "set environment",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		imageName := cmdArray[0]
		cmdArray = cmdArray[1:]

		createTty := context.Bool("ti")
		detach := context.Bool("d")
		if createTty && detach {
			return fmt.Errorf("ti and d parameter can not both provided")
		}
		res := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CpuShare: context.String("cpushare"),
			CpuSet:context.String("cpuset"),
		}
		log.Infof("createTty %v",createTty)
		containerName := context.String("name")
		volume := context.String("v")

		envSlice := context.StringSlice("e")
		log.Infof("envSlice = %s",envSlice)
		Run(createTty, cmdArray,res,containerName,volume,imageName,envSlice)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}

var commitCommand = cli.Command {
	Name: "commit",
	Usage: "commit a container into image",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 0 {
			return fmt.Errorf("Missing container name and image name")
		}
		containerName := context.Args().Get(0)
		imageName := context.Args().Get(1)
		commitContainer(containerName,imageName)
		return nil
	},
}

var listCommand = cli.Command {
	Name: "ps",
	Usage: "list all the containers",
	Action: func(context *cli.Context) error {
		ListContainers()
		return nil
	},
}

var logCommand = cli.Command {
	Name: "logs",
	Usage: "print logs of a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Please input your container name")
		}
		containerName := context.Args().Get(0)
		logContainer(containerName)
		return nil
	},
}

var execCommand = cli.Command {
	Name: "exec",
	Usage: "exec a command into container",
	Action: func(context *cli.Context) error {
		if os.Getenv(ENV_EXEC_PID) != "" {
			log.Infof("pid callback pid %d",os.Getgid())
			return nil
		}
		if len(context.Args()) < 2 {
			return fmt.Errorf("Missing container name or command")
		}

		containerName := context.Args().Get(0)
		var commandArray []string
		for _, arg := range context.Args().Tail() {
			commandArray = append(commandArray, arg)
		}
		ExecContainer(containerName, commandArray)
		return nil
	},
}

var stopCommand = cli.Command {
	Name: "stop",
	Usage: "stop a container",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		stopContainer(containerName)
		return nil
	},
}

var removeCommand = cli.Command {
	Name: "rm",
	Usage: "remove unused containers",
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container name")
		}
		containerName := context.Args().Get(0)
		removeContainer(containerName)
		return nil
	},
}