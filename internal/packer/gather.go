package packer

import (
    "context"
    "io/fs"
    "path/filepath"
)

type File struct {
    Path string
    Size int64
}

func Gather(ctx context.Context, dir string) ([]File, error) {
    var files []File
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        if d.Type().IsRegular() {
            info, err := d.Info()
            if err != nil {
                return err
            }
            files = append(files, File{Path: path, Size: info.Size()})
        }
        return nil
    })
    return files, err
}