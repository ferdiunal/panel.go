package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/cli/main.go make:resource <name>")
		return
	}

	command := os.Args[1]
	name := os.Args[2]

	if command == "make:resource" {
		makeResource(name)
	} else if command == "make:page" {
		makePage(name)
	} else if command == "make:model" {
		makeModel(name)
	} else {
		fmt.Println("Unknown command")
	}
}

func makeResource(name string) {
	// Normalize name
	// e.g. "blog" -> "Blog"
	caser := cases.Title(language.English)
	resourceName := caser.String(name)        // Blog
	packageName := strings.ToLower(name)      // blog
	identifier := strings.ToLower(name) + "s" // blogs
	label := resourceName + "s"               // Blogs
	modelName := resourceName                 // Blog (Assumes model exists or will be created)

	// Directort: internal/resource/<name>
	dir := filepath.Join("internal", "resource", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Data for templates
	data := map[string]string{
		"PackageName":  packageName,
		"ResourceName": resourceName,
		"ModelName":    modelName,
		"ModelPkg":     strings.ToLower(modelName),
		"Slug":         identifier,
		"Title":        label,
		"Label":        label,
		"Identifier":   identifier,
		"Group":        "Content", // Default group
		"Icon":         "circle",  // Default icon
	}

	// Stubs to process
	stubs := map[string]string{
		"resource.stub":   filepath.Join(dir, fmt.Sprintf("%s_resource.go", packageName)),
		"policy.stub":     filepath.Join(dir, fmt.Sprintf("%s_policy.go", packageName)),
		"repository.stub": filepath.Join(dir, fmt.Sprintf("%s_repository.go", packageName)),
	}

	for stub, target := range stubs {
		createFileFromStub(stub, target, data)
	}

	fmt.Printf("Resource %s generated successfully in %s\n", resourceName, dir)
}

func makePage(name string) {
	// Normalize name
	// e.g. "dashboard" -> "Dashboard"
	caser := cases.Title(language.English)
	pageName := caser.String(name)       // Dashboard
	packageName := strings.ToLower(name) // dashboard
	slug := strings.ToLower(name)        // dashboard
	title := pageName                    // Dashboard

	// Directory: internal/page/<name>
	// Actually, pages are usually just in internal/page/ or internal/page/<name> if complex.
	// Looking at settings.go, it's directly in internal/page.
	// But let's create a separate file for the page in internal/page/

	dir := filepath.Join("internal", "page")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, fmt.Sprintf("%s.go", packageName))

	// Data for templates
	data := map[string]string{
		"PackageName": "page", // Using 'page' package to match existing structure
		"PageName":    pageName,
		"Slug":        slug,
		"Title":       title,
		"Group":       "System",
		"Icon":        "circle",
	}

	createFileFromStub("page.stub", targetPath, data)
	fmt.Printf("Page %s generated successfully at %s\n", pageName, targetPath)
}

func makeModel(name string) {
	// Normalize name
	// e.g. "blog" -> "Blog"
	caser := cases.Title(language.English)
	modelName := caser.String(name)      // Blog
	packageName := strings.ToLower(name) // blog

	// Directory: internal/domain/<name>
	dir := filepath.Join("internal", "domain", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, "entity.go")

	// Data for templates
	data := map[string]string{
		"PackageName": packageName,
		"ModelName":   modelName,
	}

	createFileFromStub("model.stub", targetPath, data)
	fmt.Printf("Model %s generated successfully at %s\n", modelName, targetPath)
}

func createFileFromStub(stubName, targetPath string, data map[string]string) {
	// Read stub
	stubPath := filepath.Join("stubs", stubName)
	content, err := os.ReadFile(stubPath)
	if err != nil {
		fmt.Printf("Error reading stub %s: %v\n", stubName, err)
		return
	}

	// Process template
	tmpl, err := template.New(stubName).Parse(string(content))
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", stubName, err)
		return
	}

	// Create file
	f, err := os.Create(targetPath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", targetPath, err)
		return
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Printf("Error executing template %s: %v\n", stubName, err)
	}
	fmt.Printf("Created: %s\n", targetPath)
}
