package config

import (
	"context"
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"github.com/skruger/privatestudio/transcoder/config"
	"github.com/skruger/privatestudio/web/dao"
	"os"
	"path"
)

type StandardEncodingProfile struct {
	Outputs          config.Outputs `mapstructure:"Outputs"`
	GlobalArgs       []string       `mapstructure:"GlobalArgs"`
	OverwriteOutputs bool           `mapstructure:"OverwriteOutputs"`
}

type ConfigFile struct {
	AssetHome        string                             `mapstructure:"AssetHome"`
	DatabaseURL      string                             `mapstructure:"DatabaseURL"`
	EncodingProfiles map[string]StandardEncodingProfile `mapstructure:"EncodingProfiles"`
}

type Config struct {
	ConfigFile     ConfigFile
	configBasePath string
	DB *sql.DB
}

func LoadWebConfig(args []string) (*Config, error) {
	configFileData := DefaultConfig

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	basePath := path.Join(cwd, configFileData.AssetHome)

	db, err := dao.NewSqliteDb(path.Join(basePath, "sqlite3.db"))
	if err != nil {
		return nil, err
	}
	conn, err := db.Conn(context.TODO())
	if err != nil {
		return nil, err
	}
	conn.Close()

	db.Ping()
	if err != nil {
		return nil, err
	}

	c := &Config{
		ConfigFile:     DefaultConfig,
		configBasePath: basePath,
		DB: db,
	}

	return c, nil
}

func (c Config) GetSourcePath() string {
	return path.Join(c.configBasePath, "source")
}

func (c Config) GetTranscodePath() string {
	return path.Join(c.configBasePath, "transcode")
}

func (c Config) GetManifestPath() string {
	return path.Join(c.configBasePath, "manifest")
}
