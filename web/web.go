package web

import assetfs "github.com/elazarl/go-bindata-assetfs"

//go:generate ember build -prod

//go:generate go-bindata-assetfs -pkg web dist/...

// AssetFS retruns http.FileSystem
func AssetFS() *assetfs.AssetFS {
	return assetFS()
}

func init() {
}
