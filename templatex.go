package templatex

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ParseError represents an error which can occure when trying to parse a template
type ParseError struct {
	Path string
	Err  error
}

func (e ParseError) Error() string {
	return fmt.Sprintf(`error while parsing "%v": %v`, e.Path, e.Err)
}

func (e ParseError) Unwrap() error {
	return e.Err
}

// NotFound represents an error which can occure when trying to execute a template,
// which does not exist
type NotFoundError struct {
	Template string
}

func (e NotFoundError) Error() string {
	return "template not found: " + e.Template
}

// ExecuteError represents an error which can occure while trying to execute a template
type ExecuteError struct {
	Template string
	Err      error
}

func (e ExecuteError) Error() string {
	return fmt.Sprintf(`error executing template "%v": %v`, e.Template, e.Err)
}

func (e ExecuteError) Unwrap() error {
	return e.Err
}

// New creates a new Template with sane default values for directories like:
// templates/
//   layout.html
//
//   includes/
//     header.html
//     footer.html
//
//   profile/
//     view.html
//     edit.html
func New() *Template {
	return &Template{
		Layout:     "layout.html",
		IncludeDir: "includes",
	}
}

// Template represents a container for multiple templates parsed from a directory
type Template struct {
	// Layout specifies the filename of the layout files in a directory
	// Most commonly: "layout.html" or "base.html"
	Layout string

	// IncludeDir specifies the directory name where partial templates can be found
	// Most commonly: "includes", "include" or "inc"
	IncludeDir string

	// FuncMap is a map of functions, given to the templates while parsing
	FuncMap template.FuncMap

	// templates is a map of template identifiers to executable templates
	templates map[string]*template.Template
}

// ParseDir parses all templates inside a given directory
func (t *Template) ParseDir(dir string) (err error) {
	t.templates, err = parseFS(os.DirFS("."), dir, t.IncludeDir, t.Layout, t.FuncMap)
	return
}

// ParseFS parses all templates inside a given fs.FS
func (t *Template) ParseFS(files fs.FS, dir string) (err error) {
	t.templates, err = parseFS(files, dir, t.IncludeDir, t.Layout, t.FuncMap)
	return
}

// ExecuteTemplate executes a template by its name to a io.Writer with any given data
func (t *Template) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return NotFoundError{Template: name}
	}
	if err := tmpl.Execute(w, data); err != nil {
		return ExecuteError{Template: name, Err: err}
	}
	return nil
}

// parseDir builds templates inside a given directory
func parseFS(files fs.FS, dir, includeDir, layout string, funcMap template.FuncMap) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	// Collect template parsing information of the given directory
	templateInfos, err := findTemplates(files, dir, includeDir, layout)
	if err != nil {
		return nil, err
	}

	for _, info := range templateInfos {
		var err error

		// Create a new empty layout with the name of the layout file
		t := template.New(layout).Funcs(funcMap)

		// Use ParseGlob to parse all partial templates from the include directories
		for _, f := range info.includes {
			gf, err := fs.Glob(files, f+string(filepath.Separator)+"*")
			if err != nil {
				return nil, ParseError{Path: f + string(filepath.Separator) + "*", Err: err}
			}
			t, err = t.ParseFS(files, gf...)
			if err != nil {
				return nil, ParseError{Path: f + string(filepath.Separator) + "*", Err: err}
			}
		}

		// Parse the rest of the templates
		t, err = t.ParseFS(files, info.files...)
		if err != nil {
			return nil, ParseError{Path: fmt.Sprintf("%v", info.files), Err: err}
		}

		// Add the parsed template to the template map
		templates[info.id] = t
	}

	return templates, nil
}

// templateInfo contains all template information neccessary to parse a
// final template with its dependencies (layout templates, include templates)
// It also contains an identifier for the resulting template to execute
type templateInfo struct {
	id       string
	includes []string
	files    []string
}

// fileTemplates returns a list of all executable templates
// with their respective layout dependencies and include templates
func findTemplates(files fs.FS, dir, includeDir, layout string) ([]templateInfo, error) {
	// Cleans trailing slashs from directories
	dir = filepath.Clean(dir)
	includeDir = filepath.Clean(includeDir)

	// Slices to hold all found files and directories
	includeDirs := []string{}
	layouts := []string{}
	templates := []string{}

	// walkfn finds all files and directories inside of dir
	walkfn := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Check if the found directory is an include directory
			if filepath.Base(path) == includeDir {
				includeDirs = append(includeDirs, path)
			}
			return nil
		}

		// Skip all templates found in an include directory
		if filepath.Base(filepath.Dir(path)) == includeDir {
			return nil
		}

		// Determine if the template is a base/layout template or normal template
		if strings.HasSuffix(path, string(filepath.Separator)+layout) {
			layouts = append(layouts, path)
		} else {
			templates = append(templates, path)
		}

		return nil
	}
	if err := fs.WalkDir(files, dir, walkfn); err != nil {
		return nil, ParseError{Path: dir, Err: err}
	}

	// Sort all base/layout templates by their directory depth (shallow to deep)
	sort.Slice(layouts, func(i, j int) bool {
		return len(layouts[i]) < len(layouts[j])
	})

	// Sort all include directories by their directory depth (shallow to deep)
	sort.Slice(includeDirs, func(i, j int) bool {
		return len(includeDirs[i]) < len(includeDirs[j])
	})

	// For each found normal template, build a list of dependencies to parse
	templateInfos := []templateInfo{}
	for _, t := range templates {
		files := []string{}
		includes := []string{}

		// Add all include directories which lie in the same directory hirachy
		for _, i := range includeDirs {
			if strings.HasPrefix(t, filepath.Dir(i)) {
				includes = append(includes, i)
			}
		}

		// Add all base/layout templates which lie in the same directory hirachy
		for _, l := range layouts {
			if strings.HasPrefix(t, filepath.Dir(l)) {
				files = append(files, l)
			}
		}

		// Add the final template as the last entry
		files = append(files, t)

		// Build the template identifier based on the path of the final template
		// e.g. <dir>/profile/view.html -> profile/view
		id := strings.TrimPrefix(t, dir+string(filepath.Separator))
		id = strings.TrimSuffix(id, filepath.Ext(id))

		templateInfos = append(templateInfos, templateInfo{
			id:       id,
			includes: includes,
			files:    files,
		})
	}

	return templateInfos, nil
}
