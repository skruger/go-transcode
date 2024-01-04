package localfile

import (
	"encoding/json"
	"fmt"
	"github.com/skruger/privatestudio/web/datastore"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type LocalFileAssets struct {
	SourceRoot     string
	TranscodedRoot string
	ManifestRoot   string
}


func NewLocalFileAssetStorage(sourcePath string, transcodedPath string, manifestPath string) (datastore.AssetStorage, error) {
	return &LocalFileAssets{
		SourceRoot:     sourcePath,
		TranscodedRoot: transcodedPath,
		ManifestRoot:   manifestPath,
	}, nil
}

type Source struct {
	Key      string
	Filename string
}

func (s *Source) GetKey() string {
	return s.Key
}

func (s *Source) GetLocalFilename(_ bool) (string, error) {
	return s.Filename, nil
}

func (s Source) GetAssetMetadata() (*datastore.AssetMetadata, error) {
	return GetAssetMetadata(fmt.Sprintf("%s.metadata", s.Filename))
}

func (s Source) SaveAssetMetadata(metadata datastore.AssetMetadata) error {
	metadata.Type = "source"
	if fileInfo, err := os.Stat(s.Filename); err == nil {
		metadata.UpdatedAt = fileInfo.ModTime()
	}
	metadataFile := fmt.Sprintf("%s.metadata", s.Filename)
	return SaveAssetMetadata(metadataFile, metadata)
}

func GetAssetMetadata(filename string) (*datastore.AssetMetadata, error) {
	var metadata *datastore.AssetMetadata

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read metadata: %s", err)
	}
	err = json.Unmarshal(data, metadata)
	return metadata, err
}

func SaveAssetMetadata(filename string, metadata datastore.AssetMetadata) error {
	data, err := json.MarshalIndent(metadata, "  ", " ")
	if err != nil {
		return fmt.Errorf("unable to serialize metadata: %s", err)
	}

	return os.WriteFile(filename, data, 666)
}

func (a LocalFileAssets) GetSourceAsset(key string) (datastore.SourceAsset, error) {
	filename := path.Join(a.SourceRoot, key)
	_, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to find source asset: err")
	}
	source := &Source{
		Key:      key,
		Filename: filename,
	}
	return source, nil
}

func (a LocalFileAssets) FindSourceAssets(filters map[string]string, offset string) ([]datastore.SourceAsset, error) {
	fc := fileCollector{}
	filepath.Walk(a.SourceRoot, fc.addFile)

	result := []datastore.SourceAsset{}

	for _, file := range fc.fileList {
		_, baseName := filepath.Split(file.path)
		relativePath, err := filepath.Rel(a.SourceRoot, file.path)
		if err != nil {
			relativePath = baseName
		}
		result = append(result, &Source{
			Key: relativePath,
			Filename: file.path,
		})
	}

	return result, nil
}

func (a LocalFileAssets) GetTranscodedAsset(key string) (datastore.TranscodedAsset, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a LocalFileAssets) FindTranscodedAssets(filters map[string]string, offset string) ([]datastore.TranscodedAsset, error) {
	return []datastore.TranscodedAsset{}, fmt.Errorf("not implemented")
}

func (a LocalFileAssets) GetManifestAsset(key string) (datastore.ManifestAsset, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a LocalFileAssets) FindManifestAssets(filters map[string]string, offset string) ([]datastore.ManifestAsset, error) {
	return []datastore.ManifestAsset{}, fmt.Errorf("not implemented")
}

type fileEntry struct {
	path string
	info fs.FileInfo
}

type fileCollector struct {
	fileList []fileEntry
}

func (fc *fileCollector) addFile(path string, info fs.FileInfo, err error) error {
	if err == nil && !info.IsDir(){
		fc.fileList = append(fc.fileList, fileEntry{path: path, info: info})
	}
	return nil
}
