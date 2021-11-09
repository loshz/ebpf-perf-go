# eBPF Performance Evaluations for Go
[![Build Status](https://github.com/syscll/ebpf-perf-go/workflows/ci/badge.svg)](https://github.com/syscll/ebpf-perf-go/actions)

This repository contains evaluations for creating and running eBPF programs using 3rd party Go libraries.

Currently supporting the [Cilium](https://github.com/cilium/ebpf) library.

## Usage
In order to compile the eBPF bytecode you will need the required header files:
- `./bpf/headers/vmlinux.h`: auto-generated using vmlinux BTF 5.14.16
- [libbpf](https://github.com/libbpf/libbpf) must be installed.

## eBPF
- `./bpf/tracepoint.c`: a basic tracepoint for the `mkdir` syscall.
