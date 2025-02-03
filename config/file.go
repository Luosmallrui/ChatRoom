package config

import "chatroom/pkg/filesystem"

type Filesystem struct {
	Default string                       `json:"default" yaml:"default"`
	Local   filesystem.LocalSystemConfig `json:"local" yaml:"local"`
	Minio   filesystem.MinioSystemConfig `json:"minio" yaml:"minio"`
}

func NewFilesystem(conf *Config) filesystem.IFilesystem {
	if conf.Filesystem.Default == filesystem.MinioDriver {
		return filesystem.NewMinioFilesystem(conf.Filesystem.Minio)
	}

	return filesystem.NewLocalFilesystem(conf.Filesystem.Local)
}
