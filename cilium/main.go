//go:build linux

// The eBPF program will be attached to the page allocation tracepoint and
// prints out the number of times it has been reached. The tracepoint fields
// are printed into /sys/kernel/debug/tracing/trace_pipe.
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go bpf ../bpf/tracepoint.c -- -I../bpf

func main() {
	// Subscribe to signals for terminating the program.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("error loading objects: %v", err)
	}
	defer objs.Close()

	// Open a tracepoint and attach the pre-compiled program.
	kp, err := link.Tracepoint("syscalls", "sys_enter_mkdir", objs.SysEnterMkdir)
	if err != nil {
		log.Fatalf("error opening tracepoint: %v", err)
	}
	defer kp.Close()

	// Read loop reporting the total amount of times the kernel
	// function was entered, once per second.
	t := time.NewTicker(time.Second)
	fmt.Println("Polling for mkdir events...")
	for {
		select {
		case <-t.C:
			var val uint64
			if err := objs.SyscallCount.Lookup(uint32(0), &val); err != nil {
				log.Fatalf("error reading ebpf map: %v", err)
			}
			fmt.Printf("\rNo. of syscalls: %d", val)
		case <-stop:
			return
		}
	}
}
