
## before build 
login to SWR repo


## build
```shell
cd example
export GOPROXY=https://goproxy.io
GO111MODULE=on go mod vendor
export REPO={swr_repo}/{org}
build_image.sh
```




