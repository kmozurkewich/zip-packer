package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"

    "zip-packer/internal/packer"
)

func main() {
    start := time.Now()

    sizeMiB := flag.Int64("size", 0, "minimum target size per archive in MiB (required)")
    outBase := flag.String("out", "archive", "output file base name (archive_N.zip)")
    flag.Parse()

       if flag.NArg() != 1 || *sizeMiB <= 0 {
        fmt.Fprintf(os.Stderr, "usage: %s -size <MiB> [-out base] <inputDir>", filepath.Base(os.Args[0]))
        flag.PrintDefaults()
        os.Exit(2)
    }

    absRoot, err := filepath.Abs(flag.Arg(0))
    if err != nil {
        log.Fatalf("resolve input path: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
    defer cancel()

    files, err := packer.Gather(ctx, absRoot)
    if err != nil {
        log.Fatal(err)
    }

    var totalBytes int64
    for _, f := range files {
        totalBytes += f.Size
    }

    partitions, err := packer.Partition(files, *sizeMiB*1024*1024)
    if err != nil {
        log.Fatal(err)
    }

    if err := packer.BuildZips(ctx, partitions, *outBase); err != nil {
        log.Fatal(err)
    }

    elapsed := time.Since(start)

    fmt.Printf("Done. Archives: %d \nFiles: %d \nTotal: %.2f MiB \nElapsed: %s", len(partitions), len(files), float64(totalBytes)/(1024*1024), elapsed)
}