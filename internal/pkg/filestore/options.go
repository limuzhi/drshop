package filestore

type FileStoreOptions interface {
    NewOptions() error
	UpLoad(pathName string, localFile string) error
}
