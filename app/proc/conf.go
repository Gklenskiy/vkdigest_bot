package proc

// Conf for sources config yml
type Conf struct {
	Sources map[string]Source `yaml:"sources"`
}

// Source defines config section
type Source struct {
	Title      string     `yaml:"title"`
	BaseURL    string     `yaml:"base_url"`
	APIVersion string     `yaml:"api_version"` // nolint
	Domains    []VkDomain `yaml:"domains"`
}

// VkDomain explains configuration vk record
type VkDomain struct {
	Title string `yaml:"title"`
	Name  string `yaml:"name"`
}
