# PipeCheck (pc) [![build-and-test](https://github.com/goodmustache/pc/workflows/build-and-test/badge.svg?branch=main)](https://github.com/goodmustache/pc/actions/workflows/build-test.yml)
Pipecheck is a tool to help debug / test pipe'd commands on your terminal:

```bash
$ echo "foo\nbar\nbaz" | pc | xargs -n 1 echo -
==================================================================
 PipeCheck: The following was read in and will be passed through:
==================================================================
foo
bar
baz
===============================
 Proceed with this data (y/N):y
===============================
- foo
- bar
- baz
```

*Note*:
When an `n` or any other key is provided, the command will terminate and send an _empty buffer_ to the subsequent commands.
IE *the string of pipes do not terminate*, but continue one with an empty buffer.

## How to install

Currently the only way to install:

```
go install github.com/goodmustache/pc@latest
```

### What's the difference between this and `tee /dev/stderr`?

Honestly?
Not much - you get a pause to validate if you want things to continue.
I don't know how much this helps others, but I found it useful.
