package helper

// Helper collect helper functions
type Helper struct {
	Log  *Loger
	Conf *Config
}

//Destroy closing opened files, databases and etc
func (hlp *Helper) Destroy() {
	hlp.Conf.Destroy()
	hlp.Log.Destroy()
}

//NewHelper Create new helper
func NewHelper(configFile string) (helper *Helper) {
	helper = new(Helper)
	helper.Log = NewLoger()
	helper.Conf = NewConfig(configFile)
	return helper
}
