package model

type Default struct {
	AddressRB string `mapstructure:"ADDRESSRB"`
	UserRB string `mapstructure:"USERRB"`
	PassRB string `mapstructure:"PASSRB"`
	PortRB int `mapstructure:"PORTRB"`
	Queue string `mapstructure:"QUEUE"`
	Addr string `mapstructure:"ADDRESS"`
	Port string `mapstructure:"PORT"`
	Loglevel bool `mapstructure:"LOGLEVEL"`
	Botname string `mapstructure:"BOTNAME"`
	Rootpath string `mapstructure:"ROOTPATH"`
	GitlabURL string `mapstructure:"GITLABURL"`
	GitlabToken string `mapstructure:"GITLABTOKEN"`
	Clusterenv string `mapstructure:"CLUSTERENV"`
}