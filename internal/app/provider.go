package app

import (
	auth_module "github.com/Fi44er/cloud-store-api/internal/modules/auth"
	file_module "github.com/Fi44er/cloud-store-api/internal/modules/files"
)

type moduleProvider struct {
	app *App

	authModule *auth_module.AuthModule
	fileModule *file_module.FileModule
}

func NewModuleProvider(app *App) (*moduleProvider, error) {
	provider := &moduleProvider{
		app: app,
	}

	err := provider.initDeps()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *moduleProvider) initDeps() error {
	inits := []func() error{
		p.AuthModule,
		p.FileModule,
	}
	for _, init := range inits {
		err := init()
		if err != nil {
			p.app.logger.Errorf("%s", "✖ Failed to initialize module: "+err.Error())
			return err
		}
	}
	return nil
}

func (p *moduleProvider) AuthModule() error {
	p.authModule = auth_module.NewAuthModule(
		p.app.logger,
	)
	p.authModule.Init()
	return nil
}

func (p *moduleProvider) FileModule() error {
	p.fileModule = file_module.NewFileModule(
		p.app.logger,
		p.app.db,
	)
	p.fileModule.Init()
	return nil
}
