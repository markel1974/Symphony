package modules

// IModuleGetter enables implementing dynamic module loading.
type IModuleGetter interface {
	Get(name string) Importable
}
