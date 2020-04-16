# Example

This example demonstrates how to use the `github.com/philippta/templatex` package.

## Code

Create a new `templatex.Template` with custom parameters:
```go
tmpl := &templatex.Template{
    // Use "layout.html" as base layout name
    Layout:     "layout.html",

    // Use "includes" as directory for partial templates
    IncludeDir: "includes",

    // Add strings.ToUpper as a template function
    FuncMap: template.FuncMap{
        "uppercase": strings.ToUpper,
    },
}
```

Parse the `templates/` directory:
```go
tmpl.ParseDir("templates/")
```

Execute the templates:
```go
tmpl.ExecuteTemplate(os.Stdout, "profile/view", viewParams)
tmpl.ExecuteTemplate(os.Stdout, "profile/edit", editParams)
tmpl.ExecuteTemplate(os.Stdout, "profile/payment/methods", paymentMethodsParams)
```

## Result

The result of these three demo pages will look like this (whitespaces have been cleaned up):

### profile/view
```html
<!DOCTYPE html>
<head>
    <title>templatex Demo</title>
</head>

<body>

    <h1>PROFILE</h1>

    <div class="content">
        <div class="profile">
            <dl>
                <dt>Name</dt>
                <dd>John Doe</dd>

                <dt>Address</dt>
                <dd>Mainstreet 1st</dd>

                <dt>Subscription</dt>
                <dd>Premium</dd>
            </dl>
        </div>
    </div>

    <div class="copyright">
        &copy Company 2020. All rights reserved.
    </div>

</body>
```

### profile/edit
```html
<!DOCTYPE html>
<head>
    <title>templatex Demo</title>
</head>

<body>

    <h1>PROFILE EDIT</h1>

    <div class="content">
        <div class="profile">
            <form>

                <label>
                    Name:
                    <input type="text" value="John Doe" />
                </label>

                <label>
                    Address:
                    <input type="text" value="Mainstreet 1st" />
                </label>

            </form>
        </div>
    </div>


    <div class="copyright">
        &copy Company 2020. All rights reserved.
    </div>

</body>
```

### profile/payment/methods
```html
<!DOCTYPE html>
<head>
    <title>templatex Demo</title>
</head>

<body>

    <h1>PAYMENT METHODS</h1>

    <div class="content">
        <div class="profile">
            <div class="payments">

                <h2>The following payment methods are available</h2>

                <div class="method">
                    <div id="method-debit">
                        Debit card
                    </div>
                </div>

                <div class="method">
                    <div id="method-cc">
                        Credit card
                    </div>
                </div>

                <div class="method">
                    <div id="method-paypal">
                        Paypal
                    </div>
                </div>

            </div>
        </div>
    </div>


    <div class="copyright">
        &copy Company 2020. All rights reserved.
    </div>

</body>
```
