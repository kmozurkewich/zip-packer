# zip‑packer

Prototype CLI that partitions an arbitrary set of files into ZIP archives, each at least a minimum size.

## Build

```bash
go build ./cmd/zip-packer
```

## Usage

```
zip-packer <inputDir> -size 2048 -out backup
```

* **`<inputDir>`** Positional argument: root directory containing files to pack.
* **`-size`**     Minimum size *per* archive in **MiB**. Archives will be at least this size (except possibly the last one if the remaining data is smaller).
* **`-out`**       Base name for output archives (`backup_1.zip`, `backup_2.zip`, etc…).


The tool applies a greedy best‑fit‑decreasing bin‑packing heuristic followed by a balancing pass to ensure every archive ≥ **size**.  If the first file of an archive is > **size** then the archive will have a single file.

## Example


```
 zip-packer -size 64 /data/
 ```

Creates:

archive_1.zip
archive_2.zip
archive_3.zip




