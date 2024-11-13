# tsv

Pretty print tab separated values (TSVs)

## Installation

You can install `tsv` by running the install script which will download
the [latest release](https://github.com/mskelton/tsv/releases/latest).

```bash
curl -LSfs https://go.mskelton.dev/tsv/install | sh
```

Or you can build from source.

```bash
git clone git@github.com:mskelton/tsv.git
cd tsv
go install .
```

## Usage

```bash
echo -e "Id\tName\tAge\n18\tAlice\t30\n2\tBob\t25" | tsv \
    --column Id \
    --column Age type=number \
    --column Name
```
