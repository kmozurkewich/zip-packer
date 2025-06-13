# zip‑packer

Prototype CLI that partitions an arbitrary set of files into **k** ZIP archives, each at least a minimum size.

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


The tool applies a greedy best‑fit‑decreasing bin‑packing heuristic followed by a balancing pass to ensure every archive ≥ floor.

## Example


```
 zip-packer -in ./dataset -k 3 -size 512 -out backup 
 ```

Creates:

backup_1.zip
backup_2.zip
backup_3.zip




