package packer

import (
    "archive/zip"
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

func BuildZips(ctx context.Context, partitions [][]File, outBase string) error {
    for i, part := range partitions {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        zipName := fmt.Sprintf("%s_%d.zip", outBase, i+1)
        if err := writeOne(zipName, part); err != nil {
            return err
        }
    }
    return nil
}

func writeOne(zipPath string, files []File) error {
    f, err := os.Create(zipPath)
    if err != nil {
        return err
    }
    defer f.Close()

    zw := zip.NewWriter(f)
    defer zw.Close()

    for _, file := range files {
        rel := filepath.Base(file.Path)
        hdr, err := zip.FileInfoHeader(&fakeInfo{size: file.Size, name: rel})
        if err != nil {
            return err
        }
        hdr.Method = zip.Deflate
        w, err := zw.CreateHeader(hdr)
        if err != nil {
            return err
        }
        src, err := os.Open(file.Path)
        if err != nil {
            return err
        }
        if _, err := io.Copy(w, src); err != nil {
            src.Close()
            return err
        }
        src.Close()
    }
    return nil
}

type fakeInfo struct {
    size int64
    name string
}

func (fi *fakeInfo) Name() string       { return fi.name }
func (fi *fakeInfo) Size() int64        { return fi.size }
func (fi *fakeInfo) Mode() os.FileMode   { return 0444 }
func (fi *fakeInfo) ModTime() time.Time  { return time.Time{} }
func (fi *fakeInfo) IsDir() bool         { return false }
func (fi *fakeInfo) Sys() any            { return nil }