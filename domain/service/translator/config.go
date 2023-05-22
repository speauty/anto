package translator

type ImplConfig interface {
	GetAK() string // app-key or access-key
	GetSK() string // secret-key
	GetPK() string // project-key
	GetQPS() int
	GetTML() int // text-max-length
	GetPM() int  // proc-max
	SetAK(ak string) error
	SetSK(sk string) error
	SetPK(pk string) error
	SetQPS(qps int) error
	SetTML(tml int) error
	SetPM(pm int) error
	Default() ImplConfig
	Sync() error
}

type DefaultConfig struct{}

func (d DefaultConfig) GetAK() string { return "" }

func (d DefaultConfig) GetSK() string { return "" }

func (d DefaultConfig) GetPK() string { return "" }

func (d DefaultConfig) GetQPS() int { return 0 }

func (d DefaultConfig) GetTML() int { return 0 }

func (d DefaultConfig) GetPM() int { return 0 }

func (d DefaultConfig) SetAK(ak string) error { return nil }

func (d DefaultConfig) SetSK(sk string) error { return nil }

func (d DefaultConfig) SetPK(pk string) error { return nil }

func (d DefaultConfig) SetQPS(qps int) error { return nil }

func (d DefaultConfig) SetTML(tml int) error { return nil }

func (d DefaultConfig) SetPM(pm int) error { return nil }

func (d DefaultConfig) Default() ImplConfig { return nil }

func (d DefaultConfig) Sync() error { return nil }
