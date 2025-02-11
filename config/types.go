package config

type ManagerConfig struct {
	ScreenName string   `json:"screen_name"`
	Servers    []Server `json:"servers"`
}

type Server struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	MaxRAM     int        `json:"max_ram"`
	JarPath    string     `json:"jar_path"`
	JavaPath   string     `json:"java_path"`
	Type       ServerType `json:"type"`
	IsSelected bool       `json:"is_selected"`
}

type ServerType int

const (
	ServerType_Vanilla ServerType = 0
	ServerType_Paper   ServerType = 1
	ServerType_Fabric  ServerType = 2
)
