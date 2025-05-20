package symphony

type Options struct {
	RenderId   string
	PlayerId   string
	Prg        string
	Cartridges string
	Drives     string
	Disks      string
	NoJiffy    bool
}

func NewOptions() *Options {
	return &Options{}
}
