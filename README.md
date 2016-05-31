## Godtrace
dtrace in Go

### Enable kernel modules in EI Capitan

If you got the same error as follows:

```shell
$ sudo dtrace -n 'io:::{}'
dtrace: invalid probe specifier io:::{}: probe description io::: does not match any probes
```

Then you need to disable rootless for dtrace.

```
Hold ⌘R during reboot
From the Utilities menu, run Terminal
csrutil enable --without dtrace
reboot
```

### Libdtrace document
http://dev.lrem.net/tcldtrace/wiki/LibDtrace

### Dynamic Tracing Guide
http://dtrace.org/guide/bookinfo.html

### Dtrace Tools
http://www.brendangregg.com/dtrace.html
