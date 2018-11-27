package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	cmd := exec.Command("ping", "www.baidu.com")
	ppReader, err := cmd.StdoutPipe()
	defer ppReader.Close()
	var bufReader = bufio.NewReader(ppReader)
	if err != nil {
		fmt.Printf("create cmd stdoutpipe failed,error:%s\n", err)
		os.Exit(1)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Printf("cannot start cmd,error:%s\n", err)
		os.Exit(1)
	}
	go func() {
		var buffer []byte = make([]byte, 4096)
		for {
			n, err := bufReader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("pipe has Closed\n")
					break
				} else {
					fmt.Println("Read content failed, err: %v", err)
				}
			}
			fmt.Print(string(buffer[:n]))
		}
	}()
	time.Sleep(10 * time.Second)
	err = stopProcess(cmd)
	if err != nil {
		fmt.Printf("stop child process failed,error:%s", err)
		os.Exit(1)
	}
	cmd.Wait()
	time.Sleep(1 * time.Second)
}

func stopProcess(cmd *exec.Cmd) error {
	pro, err := os.FindProcess(cmd.Process.Pid)
	if err != nil {
		return err
	}
	err = pro.Signal(syscall.SIGINT)
	if err != nil {
		return err
	}
	fmt.Printf("结束子进程%s成功\n", cmd.Path)
	return nil
}
