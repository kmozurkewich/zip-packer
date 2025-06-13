package packer

import (
    "errors"
    "sort"
)

// Partition groups files into as few bins as possible such that every bin's total
// size is >= size.  It returns the resulting partitions in order of creation.
// Last bin may be < floor **only** if the total data is < floor; otherwise we
// rebalance or merge so all bins (except possibly the lone final one) meet the threshold.
func Partition(files []File, floor int64) ([][]File, error) {
    if floor <= 0 {
        return nil, errors.New("floor must be > 0")
    }

    var total int64
    for _, f := range files {
        total += f.Size
    }
    if total < floor {
        return nil, errors.New("total data smaller than floor; nothing to do")
    }

    // Largest‑first ordering works best for balancing plus guarantees big files
    // get their own spot immediately.
    sort.Slice(files, func(i, j int) bool { return files[i].Size > files[j].Size })

    var (
        bins     [][]File
        curBin   []File
        curTotal int64
    )

    pushBin := func() {
        if len(curBin) > 0 {
            bins = append(bins, curBin)
            curBin = nil
            curTotal = 0
        }
    }

    for _, f := range files {
        // If adding this file would push us well past the floor and we already
        // satisfied the floor, start a new bin.
        if curTotal >= floor && curTotal+f.Size > floor {
            pushBin()
        }
        curBin = append(curBin, f)
        curTotal += f.Size
    }
    pushBin()

    // If the final bin is below the threshold and we have >1 bins, try to
    // rebalance by moving the smallest transferable file(s) from earlier bins.
    last := len(bins) - 1
    lastSize := sumSize(bins[last])

    if len(bins) > 1 && lastSize < floor {
        // iterate bins from first to penultimate
        for i := 0; i < last && lastSize < floor; i++ {
            // move smallest file from bin i if after moving both bins stay >= floor
            sort.Slice(bins[i], func(a, b int) bool { return bins[i][a].Size < bins[i][b].Size })
            for idx, file := range bins[i] {
                sizeI := sumSize(bins[i])
                if sizeI-file.Size >= floor {
                    // move it
                    bins[last] = append(bins[last], file)
                    bins[i] = append(bins[i][:idx], bins[i][idx+1:]...)
                    lastSize += file.Size
                    break
                }
            }
        }
    }

    // If balancing still failed, merge last two bins (guaranteed >= floor)
    if len(bins) > 1 && sumSize(bins[len(bins)-1]) < floor {
        bins[len(bins)-2] = append(bins[len(bins)-2], bins[len(bins)-1]...)
        bins = bins[:len(bins)-1]
    }

    return bins, nil
}

func sumSize(bin []File) int64 {
    var s int64
    for _, f := range bin {
        s += f.Size
    }
    return s
}