package conf

// DbConfiguration represents a set of attributes required to form database connection string.
type DbConfiguration struct {
	Driver   string `yaml:"driver"`
	Protocol string `yaml:"protocol"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DbName   string `yaml:"dbname"`
}
