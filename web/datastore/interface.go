package datastore

import (
	"time"
)

type UserInfo interface {
	GetUsername() string
	SetPassword(password string) error
	GetGroups() []string
}

type Login interface {
	CheckLogin(username string, password string) (UserInfo, error)
	CreateUser(username string) (UserInfo, error)
	GetUser(username string) (UserInfo, error)
}

type UpstreamAsset struct {
	Type string `json:"Type"`
	Key  string `json:"Key"`
}

type AssetMetadata struct {
	Name          string            `json:"Name"`
	Type          string            `json:"Type"`
	CreatedAt     time.Time         `json:"CreatedAt"`
	UpdatedAt     time.Time         `json:"UpdatedAt"`
	Tags          map[string]string `json:"Tags"`
	UpstreamAsset *UpstreamAsset    `json:"UpstreamAsset"`
}

type BaseAsset interface {
	GetKey() string
	GetLocalFilename(download bool) (string, error)
	GetAssetMetadata() (*AssetMetadata, error)
	SaveAssetMetadata(metadata AssetMetadata) error
}

//type SourceAsset interface {
//	BaseAsset
//}

type TranscodedAsset interface {
	BaseAsset
}

type ManifestAsset interface {
	BaseAsset
}

type AssetStorage interface {
	//GetSourceAsset(key string) (SourceAsset, error)
	//FindSourceAssets(filters map[string]string, offset string) ([]SourceAsset, error)

	GetTranscodedAsset(key string) (TranscodedAsset, error)
	FindTranscodedAssets(filters map[string]string, offset string) ([]TranscodedAsset, error)

	GetManifestAsset(key string) (ManifestAsset, error)
	FindManifestAssets(filters map[string]string, offset string) ([]ManifestAsset, error)
}

type ProfileStorage interface {
	GetOutputSet(key string)
}
