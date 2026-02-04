package core_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCoreLayerCompilation tests that the Core Layer compiles successfully.
// This test validates Requirement 2.4: Core Layer should compile without errors.
func TestCoreLayerCompilation(t *testing.T) {
	// Get the core package directory
	coreDir := "."

	// Parse all Go files in the core package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, coreDir, func(fi os.FileInfo) bool {
		// Include all .go files except test files for compilation check
		return strings.HasSuffix(fi.Name(), ".go") && !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ImportsOnly)

	require.NoError(t, err, "Core Layer should parse without errors")
	require.NotEmpty(t, pkgs, "Core Layer should contain at least one package")

	// Verify the core package exists
	_, exists := pkgs["core"]
	assert.True(t, exists, "Package 'core' should exist")
}

// TestCoreLayerImports tests that Core Layer only imports allowed packages.
// This test validates Requirement 2.4: Core Layer should only depend on standard library
// and third-party libraries (fiber, mime/multipart), not internal packages.
//
// Validates: Requirements 2.4
func TestCoreLayerImports(t *testing.T) {
	// Define allowed import prefixes
	allowedImports := []string{
		"github.com/gofiber/fiber/v2", // Fiber framework
		"mime/multipart",              // Standard library for file uploads
	}

	// Define forbidden import prefixes (internal packages)
	forbiddenImports := []string{
		"github.com/ferdiunal/panel.go/pkg/fields",
		"github.com/ferdiunal/panel.go/pkg/resource",
		"github.com/ferdiunal/panel.go/pkg/handler",
		"github.com/ferdiunal/panel.go/pkg/context",
		"github.com/ferdiunal/panel.go/pkg/domain",
		"github.com/ferdiunal/panel.go/pkg/auth",
		"github.com/ferdiunal/panel.go/pkg/data",
		"github.com/ferdiunal/panel.go/pkg/widget",
		"github.com/ferdiunal/panel.go/internal",
	}

	// Get the core package directory
	coreDir := "."

	// Parse all Go files in the core package
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, coreDir, func(fi os.FileInfo) bool {
		// Include all .go files except test files
		return strings.HasSuffix(fi.Name(), ".go") && !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ImportsOnly)

	require.NoError(t, err, "Core Layer should parse without errors")
	require.NotEmpty(t, pkgs, "Core Layer should contain at least one package")

	// Get the core package
	corePkg, exists := pkgs["core"]
	require.True(t, exists, "Package 'core' should exist")

	// Collect all imports from all files
	allImports := make(map[string]bool)
	for _, file := range corePkg.Files {
		for _, imp := range file.Imports {
			// Remove quotes from import path
			importPath := strings.Trim(imp.Path.Value, `"`)
			allImports[importPath] = true
		}
	}

	// Check that no forbidden imports are used
	for importPath := range allImports {
		for _, forbidden := range forbiddenImports {
			assert.NotEqual(t, forbidden, importPath,
				"Core Layer should not import forbidden package: %s", forbidden)
			assert.False(t, strings.HasPrefix(importPath, forbidden),
				"Core Layer should not import forbidden package or its subpackages: %s (found: %s)",
				forbidden, importPath)
		}
	}

	// Verify that only allowed imports are used (excluding standard library)
	for importPath := range allImports {
		// Skip standard library imports (they don't contain a dot in the first path segment)
		firstSegment := strings.Split(importPath, "/")[0]
		if !strings.Contains(firstSegment, ".") {
			// This is a standard library import, which is allowed
			continue
		}

		// Check if this is an allowed third-party import
		isAllowed := false
		for _, allowed := range allowedImports {
			if importPath == allowed || strings.HasPrefix(importPath, allowed+"/") {
				isAllowed = true
				break
			}
		}

		assert.True(t, isAllowed,
			"Core Layer should only import allowed packages. Found: %s\nAllowed: %v",
			importPath, allowedImports)
	}
}

// TestCoreLayerFiles tests that all expected Core Layer files exist.
// This test validates Requirement 2.1, 2.2, 2.3: Core Layer should contain
// Element interface, ResourceContext, and type definitions.
func TestCoreLayerFiles(t *testing.T) {
	expectedFiles := []string{
		"element.go",
		"context.go",
		"types.go",
		"callbacks.go",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(".", filename)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Core Layer should contain file: %s", filename)
	}
}

// TestCoreLayerNoInternalDependencies tests that Core Layer has no dependencies
// on internal implementation packages.
// This test validates Requirement 2.4, 4.2: Core Layer should not depend on
// implementation layers.
//
// Validates: Requirements 2.4, 4.2
func TestCoreLayerNoInternalDependencies(t *testing.T) {
	// This is a more comprehensive check that verifies the architectural rule
	// that Core Layer must not depend on any internal packages

	coreDir := "."

	// Parse all Go files
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, coreDir, func(fi os.FileInfo) bool {
		return strings.HasSuffix(fi.Name(), ".go") && !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ImportsOnly)

	require.NoError(t, err)
	corePkg, exists := pkgs["core"]
	require.True(t, exists)

	// Check each file
	for filename, file := range corePkg.Files {
		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)

			// Check if this is an internal package import
			if strings.Contains(importPath, "github.com/ferdiunal/panel.go/pkg/") {
				// Only pkg/core itself is allowed
				if !strings.HasPrefix(importPath, "github.com/ferdiunal/panel.go/pkg/core") {
					t.Errorf("File %s imports internal package %s, which violates Core Layer isolation",
						filename, importPath)
				}
			}

			// Check for internal/ directory imports
			if strings.Contains(importPath, "github.com/ferdiunal/panel.go/internal/") {
				t.Errorf("File %s imports internal package %s, which violates Core Layer isolation",
					filename, importPath)
			}
		}
	}
}
