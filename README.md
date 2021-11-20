# MyDocker

`mydocker` is source code I wrote after reading `Write Docker by Yourself`.

## How to run it

> It is better run it by root account.
> And it only tested on linux(42~20.04.1-Ubuntu)

### Build

```
go build
```

### Prepare

```
cp ./buxybox.tar /root/
```

### Run

```
mydocker run -ti busybox top
```