# templatex

The missing function for parsing nested templates in Go, like you know it from Laravel, Django, etc.
* Use the whole feature set of Go's templating library
* No new syntax or keywords. Just Go's templating functions in the right order
* No dependencies, just Go's standard library
* Super easy to use

Your template structure can now look like this and parsing will be done for you:

```
templates/
  layout.html

  includes/
    header.html
    footer.html

  dashboard/
    includes/
      widgets.html
    view.html
    edit.html

  profile/
    view.html
    edit.html

    payments/
      methods.html
      add.html
```

## Installation

```
go get -u github.com/philippta/templatex
```

## Usage

templatex is very easy to use as it only has one type and three methods.
The most basic way to use it is with default options.
This will parse the `templates/` dire

```go
t := templatex.New()

t.ParseDir("templates/")

t.ExecuteTemplate(w, "profile/view", userdata)
```

The parser also has options to use different names for the includes directory and layout files. Additional template functions can also be supplied:

```go
t := &templatex.Template{
    Layout:     "layout.html",
    IncludeDir: "includes",
    FuncMap: template.FuncMap{
        "uppercase": strings.ToUpper,
    },
}
```

## Defining nested templates

Defining nested templates in Go is relatively simple. There are just three instructions you should keep in mind:
1. `block` for creating an new block
2. `define` for filling a block
3. `template` for rendering a partial/include template

### Step 1: Creating a block
Creating a block is the essential part for nested templates, as these parts can be later filled by a child template. It is created by the `block` command.
```html
<!-- layout.html -->
<body>
    <h1>
        {{block "title" .}}Default Title{{end}}
    </h1>

    {{block "content" .}}
        <p>
            This renders, when the block has not been filled.
            Otherwise this message is discarded.
            Good for default content.
        </p>
    {{end}}
</body>
```

### Step 2: Filling a block
Filling a block which was created in a parent template is done with the `define` command. So basically, we redefine the contents of a block.
```html
<!-- profile/view.html -->
{{define "title"}}
    You're on profile/view
{{end}}

{{define "content"}}
    <p>
        This will render inside the "content" block.
        There is no need to place any html outside of this block
    </p>
{{end}}
```

### Step 3: Using partial/include templates
Creating an include template is similar to step 2. Here we're not redefining a existing block. We rather define a new block which is not rendered immediately.
```html
<!-- includes/footer.html -->
{{define "footer"}}
    &copy; Company 2020. All rights reserved
{{end}}
```

Then render it like this with the `template` command:
```html
<!-- layout.html again -->
...
{{template "footer" .}} <!-- notice: no {{end}} here -->
</body>
```


## Example

Make sure to checkout the [example](examples/README.md) in this repository.

## License
MIT


## Enjoy :-)
