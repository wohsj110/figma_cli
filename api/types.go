package api

type Me struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Handle string `json:"handle"`
	ImgURL string `json:"img_url"`
}

type File struct {
	Name         string               `json:"name"`
	LastModified string               `json:"lastModified"`
	ThumbnailURL string               `json:"thumbnailUrl"`
	Version      string               `json:"version"`
	Document     Node                 `json:"document"`
	Components   map[string]Component `json:"components,omitempty"`
	Styles       map[string]Style     `json:"styles,omitempty"`
}

type Node struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	Type                string         `json:"type"`
	Visible             *bool          `json:"visible,omitempty"`
	Characters          string         `json:"characters,omitempty"`
	Children            []Node         `json:"children,omitempty"`
	AbsoluteBoundingBox *Rectangle     `json:"absoluteBoundingBox,omitempty"`
	Style               map[string]any `json:"style,omitempty"`
	Fills               []Paint        `json:"fills,omitempty"`
}

type Rectangle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Paint struct {
	Type    string `json:"type"`
	Visible *bool  `json:"visible,omitempty"`
	Color   *Color `json:"color,omitempty"`
}

type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

type NodesResponse struct {
	Name         string              `json:"name"`
	LastModified string              `json:"lastModified"`
	Version      string              `json:"version"`
	Nodes        map[string]NodeWrap `json:"nodes"`
}

type NodeWrap struct {
	Document Node `json:"document"`
}

type ImagesResponse struct {
	Err    string            `json:"err"`
	Images map[string]string `json:"images"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type User struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
	ImgURL string `json:"img_url"`
}

type Component struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Remote         bool   `json:"remote"`
	ComponentSetID string `json:"componentSetId"`
}

type Style struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Remote      bool   `json:"remote"`
	StyleType   string `json:"styleType"`
}

type VariablesResponse struct {
	Status int           `json:"status"`
	Error  bool          `json:"error"`
	Meta   VariablesMeta `json:"meta"`
}

type VariablesMeta struct {
	Variables           map[string]Variable           `json:"variables"`
	VariableCollections map[string]VariableCollection `json:"variableCollections"`
}

type Variable struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	Key                  string         `json:"key"`
	VariableCollectionID string         `json:"variableCollectionId"`
	ResolvedType         string         `json:"resolvedType"`
	ValuesByMode         map[string]any `json:"valuesByMode"`
	Scopes               []string       `json:"scopes"`
	Remote               bool           `json:"remote"`
	Description          string         `json:"description"`
	HiddenFromPublishing bool           `json:"hiddenFromPublishing"`
}

type VariableCollection struct {
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	Key                  string         `json:"key"`
	Modes                []VariableMode `json:"modes"`
	DefaultModeID        string         `json:"defaultModeId"`
	Remote               bool           `json:"remote"`
	HiddenFromPublishing bool           `json:"hiddenFromPublishing"`
}

type VariableMode struct {
	ModeID string `json:"modeId"`
	Name   string `json:"name"`
}
