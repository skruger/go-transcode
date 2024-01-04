package web

import (
	"github.com/labstack/echo/v4"
	"github.com/skruger/privatestudio/web/datastore"
)

type ServerConfig struct {
	assetStorage datastore.AssetStorage
}

func NewServerConfig(assetStorage datastore.AssetStorage) *ServerConfig {
	return &ServerConfig{
		assetStorage: assetStorage,
	}
}

func (sc *ServerConfig) getAdmin(c echo.Context) error {
	return c.Render(200, "base.html,asset.html", contextMap(c))
}

func (sc *ServerConfig) getAssetList(c echo.Context) error {
	sourceAssets, err := sc.assetStorage.FindSourceAssets(map[string]string{}, "")
	if err != nil {
		c.Logger().Infof("unable to get source assets: %s", err)
	}

	context := contextMap(c)
	context["sourceAssets"] = sourceAssets

	return c.Render(200, "asset_list_fragment.html", context)
}

func (sc *ServerConfig) getProfiles(c echo.Context) error {
	return c.Render(200, "base.html,profiles.html", contextMap(c))
}
