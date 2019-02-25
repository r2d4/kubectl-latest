# kubectl-latest

Describe, get, and logs output from the last deployed Kubernetes object in one command.

* All resources types are supported (including custom resources)
* Arbitrary flags can be passed to the underlying commands
* "get", "describe", and "logs" are the only kubectl output subcommands supported currently
* Automatically be used as a kubectl subcommand through plugin system (installed manually or with krew)

### Installing

The only requirement is that `kubectl` must be installed on your system.

Linux
```
curl -Lo kubectl-latest https://github.com/r2d4/kubectl-latest/releases/download/v0.0.1/kubectl-latest-linux-amd64 && chmod +x kubectl-latest && sudo mv kubectl-latest /usr/local/bin
```

macOS
```
curl -Lo kubectl-latest https://github.com/r2d4/kubectl-latest/releases/download/v0.0.1/kubectl-latest-darwin-amd64 && chmod +x kubectl-latest && sudo mv kubectl-latest /usr/local/bin
```

Windows

https://github.com/r2d4/kubectl-latest/releases/download/v0.0.1/kubectl-latest-windows-amd64.exe



### Examples

You can invoke the binary directly `kubectl-latest` or, as long as the binary is on your path, through kubectl with `kubectl latest`.

```bash
# Return the "get" output of the most recent resource (across all types)
kubectl latest get 

# Return the "get" output of the most recent pod, using with the pod short syntax "po"
kubectl latest po

# or equivalently
kubectl latest get po

# Return the logs of the most recently pod
kubectl latest logs

# Returns the "get" output in yaml format of the most recent deployment. 
# kubectl-latest will pass on arbitrary flags to kubectl
kubectl latest deployment -o yaml

# Return the "describe" output of the most recent service.
kubectl latest describe svc
```

### Building from source

`make`

`make install` will install the binary into $GOBIN

`make cross` will build binaries and checksum files for all targets

### Future Work (Currently not supported but contributions welcome!)

* Namespaces other than the default configured one
* Other subcommands that make sense like label, edit, delete, patch, exec, patch, cp.
* Following logs