package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/philippta/templatex"
)

type User struct {
	Name         string
	Address      string
	Subscription string
}

type PaymentMethod struct {
	ID    string
	Title string
}

type ViewTemplateParams struct {
	Title string
	User  User
}

type EditTemplateParams struct {
	Title string
	User  User
}

type PaymentMethodsTemplateParams struct {
	Title          string
	PaymentMethods []PaymentMethod
}

func main() {
	// For the default parser without any template functions or custom layout / include names:
	// tmpl := templatex.New()

	tmpl := &templatex.Template{
		Layout:     "layout.html",
		IncludeDir: "includes",
		FuncMap: template.FuncMap{
			"uppercase": strings.ToUpper,
		},
	}

	var err error

	// Parse templates/ directory
	err = tmpl.ParseDir("templates/")
	check(err)

	// Execute profile/view template to stdout
	viewParams := ViewTemplateParams{
		Title: "Profile",
		User: User{
			Name:         "John Doe",
			Address:      "Mainstreet 1st",
			Subscription: "Premium",
		},
	}
	err = tmpl.ExecuteTemplate(os.Stdout, "profile/view", viewParams)
	check(err)

	// Execute profile/edit template to stdout
	editParams := EditTemplateParams{
		Title: "Profile Edit",
		User: User{
			Name:         "John Doe",
			Address:      "Mainstreet 1st",
			Subscription: "Premium",
		},
	}
	err = tmpl.ExecuteTemplate(os.Stdout, "profile/edit", editParams)
	check(err)

	// Execute profile/payment/methods template to stdout
	paymentMethodsParams := PaymentMethodsTemplateParams{
		Title: "Payment methods",
		PaymentMethods: []PaymentMethod{
			{
				ID:    "debit",
				Title: "Debit card",
			},
			{
				ID:    "cc",
				Title: "Credit card",
			},
			{
				ID:    "paypal",
				Title: "Paypal",
			},
		},
	}
	err = tmpl.ExecuteTemplate(os.Stdout, "profile/payment/methods", paymentMethodsParams)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
