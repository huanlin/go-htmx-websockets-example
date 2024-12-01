package main

import (
	"fmt"
	"time"

	"github.com/huanlin/go-htmx-websockets-example/internal/hardware"
)

func main() {
	fmt.Println("Starting system monitor...")
	go func() {
		for {
			systemSection, err := hardware.GetSystemSection()
			if err != nil {
				fmt.Println(err)
			}

			diskSection, err := hardware.GetDiskSection()
			if err != nil {
				fmt.Println(err)
			}
			cpuSection, err := hardware.GetCpuSection()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(systemSection)
			fmt.Println(diskSection)
			fmt.Println(cpuSection)

			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(5 * time.Minute)
}
