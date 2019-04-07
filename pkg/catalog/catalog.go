package catalog

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// Entry contains the properties of a project which are recorded as part of the catalog
type Entry struct {
	Origin      string `json:"origin"`
	ProjectName string `json:"project_name"`
	Author      string `json:"author"`
	URL         string `json:"url"`
}

// Header describes the file type and schema version for compatibility purposes
type Header struct {
	Version int    `json:"version"`
	Kind    string `json:"kind"`
}

// catalogContent represents the contents of a catalog file
type catalogContent struct {
	Header  Header           `json:"file_info"`
	Entries map[string]Entry `json:"entries"`
}

// Catalog is collected information about
type Catalog struct {
	content *catalogContent
}

var (
	defaultHeader Header
)

func init() {
	defaultHeader = Header{
		Version: 1,
		Kind:    "script_catalog",
	}
}

// New creates a new empty catalog
func New() *Catalog {
	return &Catalog{
		content: &catalogContent{
			Header:  defaultHeader,
			Entries: make(map[string]Entry),
		},
	}
}

// Load loads an existing catalog in json form from the given reader
func Load(r io.Reader) (*Catalog, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// TODO: deal with upgrading the
	var content catalogContent
	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	return &Catalog{content: &content}, nil
}

// Store writes the catalog in json form to the given writer
func (c *Catalog) Store(w io.Writer) error {
	if c.content.Header != defaultHeader {
		return errors.New("catalog header version/kind does not match")
	}
	data, err := json.Marshal(c.content)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// Insert places the given entry in the catalog potentially replacing an entry
// with the name name
func (c *Catalog) Insert(entry *Entry) {
	c.content.Entries[entry.ProjectName] = *entry
}

// Entries returns all the entries in the catalog
func (c *Catalog) Entries() []Entry {
	es := make([]Entry, len(c.content.Entries))
	i := 0
	for _, v := range c.content.Entries {
		es[i] = v
		i++
	}
	return es
}
