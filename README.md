# golang simple api

## Run
```bash
go run main.go
```
or
```bash
make compose-up
```

## Containerize

```bash
make build-image
```
then
```bash
make run-container
```



# NOTE

## Runtime Metrics
## Alloc
- Bytes of allocated heap objects. This represents the amount of memory currently in use by your Go program.

## TotalAlloc
- Cumulative bytes allocated for heap objects. This is a running total of all heap memory allocated since the program started, regardless of whether it's still in use. It will always increase, even if memory is freed.

## Sys
- Total bytes of memory obtained from the OS. This includes memory allocated for the heap, stacks, and other internal data structures. It's important to note that this doesn't necessarily represent the amount of physical memory used by the program.

## HeapAlloc
- Bytes of allocated heap objects. This is identical to Alloc. It's included for compatibility reasons.
- Note: The NumGC field represents the number of garbage collection cycles completed. While it can indirectly indicate memory pressure, it's not a direct measure of memory usage.

## Notes
- Memory Units: The values are in bytes. You've converted them to megabytes (MB) for better readability.
- Garbage Collection: The Go garbage collector reclaims unused memory, but the TotalAlloc value continues to increase.
- OS Memory Management: The actual memory usage of your program might differ from the reported values due to operating system memory management techniques.
