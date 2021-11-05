#include "common.h"
#include "bpf_helpers.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct bpf_map_def SEC("maps") syscall_count = {
    .type = BPF_MAP_TYPE_ARRAY,
    .key_size = sizeof(u32),
    .value_size = sizeof(u64),
    .max_entries = 1,
};

// This tracepoint is defined in:
// /sys/kernel/debug/tracing/events/syscalls/sys_enter_mkdir
SEC("tracepoint/syscalls/sys_enter_mkdir")
int sys_enter_mkdir() {
    u32 key = 0;
    u64 init_val = 1, *count;

    count = bpf_map_lookup_elem(&syscall_count, &key);
    if (!count) {
        bpf_map_update_elem(&syscall_count, &key, &init_val, BPF_ANY);
        return 0;
    }
    __sync_fetch_and_add(count, 1);

    return 0;
}
