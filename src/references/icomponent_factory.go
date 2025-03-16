package references

type IComponentFactory interface {
	Create(IComponent, string, string) (IComponent, error)
}
