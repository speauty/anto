package translator

type ImplConfig interface {
	GetAK() string // app-key or access-key
	GetSK() string // secret-key
	GetPK() string // project-key
	GetQPS() int
	GetMaxSingleTextLength() int
	GetMaxCoroutineNum() int
	SetAK(ak string) error
	SetSK(sk string) error
	SetPK(pk string) error
	SetQPS(qps int) error
	SetMaxSingleTextLength(textLen int) error
	SetMaxCoroutineNum(coroutineNum int) error
	Default() ImplConfig
	Sync() error
}

type DefaultConfig struct{}

func (d *DefaultConfig) GetAK() string               { return "" }
func (d *DefaultConfig) GetSK() string               { return "" }
func (d *DefaultConfig) GetPK() string               { return "" }
func (d *DefaultConfig) GetMaxSingleTextLength() int { return 0 }
func (d *DefaultConfig) GetQPS() int                 { return 0 }
func (d *DefaultConfig) GetMaxCoroutineNum() int     { return 0 }

func (d *DefaultConfig) SetAK(ak string) error                     { return nil }
func (d *DefaultConfig) SetSK(sk string) error                     { return nil }
func (d *DefaultConfig) SetPK(pk string) error                     { return nil }
func (d *DefaultConfig) SetMaxSingleTextLength(textLen int) error  { return nil }
func (d *DefaultConfig) SetQPS(qps int) error                      { return nil }
func (d *DefaultConfig) SetMaxCoroutineNum(coroutineNum int) error { return nil }

func (d *DefaultConfig) Default() ImplConfig { return nil }
func (d *DefaultConfig) Sync() error         { return nil }
