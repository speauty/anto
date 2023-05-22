package volcengine

type Cfg struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

func (customC *Cfg) GetAK() string {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) GetSK() string {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) GetPK() string {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) GetQPS() int {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) GetTML() int {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) GetPM() int {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) SetAK(ak string) error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) SetSK(sk string) error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) SetPK(pk string) error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) SetQPS(qps int) error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) SetTML(tml int) error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) SetPM(pm int) error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) Sync() error {
	//TODO implement me
	panic("implement me")
}

func (customC *Cfg) Default() *Cfg {
	return &Cfg{AccessKey: "", SecretKey: ""}
}
