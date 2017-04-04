GOOS=linux GOARCH=amd64  go build -o namerctl_linux_amd64 github.com/linkerd/namerctl
GOOS=linux GOARCH=386    go build -o namerctl_linux_i386  github.com/linkerd/namerctl
GOOS=darwin GOARCH=amd64 go build -o namerctl_darwin      github.com/linkerd/namerctl
echo "releases built:"
ls namerctl_*
