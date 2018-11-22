package utils

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir   string = "./views/layout/"
	TemplateDir string = "./views/"
	TemplateExt string = ".html"
	ContextVar         = make(map[string]interface{})
)

type MetaData struct {
	Title string
	Date  string
}

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(layout string, lang string, lang2 string, files ...string) *View {
	addTemplatePath(files) // variadic params are always taken as slices
	addTemplateExt(files)  // slices are passed by reference. so changing them
	// in another func will change original string in the slice

	files = append(files, layoutFiles()...)

	bundle := GetBundle()
	localizer := i18n.NewLocalizer(bundle, lang, lang2)
	/*
		defines localizer function for easy usage in html files
		example usage:
		{{localizer "Previous"}}
	*/
	t, err := template.New("template").Funcs(template.FuncMap{
		"localizer": func(s string) string {
			return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: s})
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}

}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, map[string]interface{}{
		"data": data,
	})
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files

}

// Takes in a slice of strings representing filepaths for templates
// and it prepends the templatedir directory to each string in the slice
// e.g input {"home"} would become {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
