// Bu paket, Panel CLI uygulamasının ana giriş noktasıdır.
//
// Panel, Go tabanlı bir kod oluşturma aracıdır (code generator) ve aşağıdaki
// komutları destekler:
//   - make:resource: Yeni bir resource (kaynak) oluşturur
//   - make:page: Yeni bir sayfa oluşturur
//   - make:model: Yeni bir model (veri modeli) oluşturur
//   - plugin:create: Yeni plugin oluşturur
//   - plugin:add: Git repository'den plugin ekler
//   - plugin:remove: Plugin'i siler
//   - plugin:list: Yüklü plugin'leri listeler
//   - plugin:build: UI build alır
//
// Tüm komutlar, gömülü stub dosyalarından şablonlar kullanarak dosyalar oluşturur.
package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ferdiunal/panel.go/pkg/plugin"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Bu değişken, stubs klasöründeki tüm .stub dosyalarını gömülü dosya sistemi
// olarak içerir. Go'nun embed özelliği sayesinde, bu dosyalar derleme zamanında
// ikili dosyaya dahil edilir ve çalışma zamanında erişilebilir hale gelir.
//
// # Kullanım Senaryosu
//
// Stub dosyaları, yeni kaynaklar, sayfalar ve modeller oluştururken şablon
// olarak kullanılır. Bu sayede, tutarlı ve standartlaştırılmış kod yapısı
// sağlanır.
//
// # Önemli Notlar
//
//   - Stub dosyaları, Go template sözdizimini kullanır
//   - Dosyalar, stubs/ klasöründe *.stub uzantısıyla saklanır
//   - Gömülü dosyalar, derleme zamanında sabitlenir ve değiştirilemez
//
//go:embed stubs/*.stub
var stubsFS embed.FS

// skillsFS, skills dizinindeki tüm skill dosyalarını gömülü dosya sistemi
// olarak içerir. SDK kullanıcıları bu skill'leri kendi projelerine kopyalayabilir.
//
//go:embed skills/**/*
var skillsFS embed.FS

// rootCmd, Panel CLI'nin root command'ı.
var rootCmd = &cobra.Command{
	Use:   "panel",
	Short: "Panel.go CLI - Code generator ve plugin yönetimi",
	Long: `Panel.go CLI, Go tabanlı admin panel için kod oluşturma ve plugin yönetimi aracıdır.

Resource, page ve model oluşturabilir, plugin'leri yönetebilir ve UI build alabilirsiniz.`,
}

// Bu fonksiyon, Panel CLI uygulamasının ana giriş noktasıdır.
//
// Cobra CLI framework kullanarak komutları yönetir ve çalıştırır.
func main() {
	// Make komutları
	rootCmd.AddCommand(newMakeResourceCommand())
	rootCmd.AddCommand(newMakePageCommand())
	rootCmd.AddCommand(newMakeModelCommand())

	// Publish komutları
	rootCmd.AddCommand(newPublishStubsCommand())
	rootCmd.AddCommand(newPublishSkillsCommand())

	// Init komutu
	rootCmd.AddCommand(newInitCommand())

	// Plugin komutları
	rootCmd.AddCommand(plugin.NewPluginCommand())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// newMakeResourceCommand, make:resource komutunu oluşturur.
func newMakeResourceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:resource <name>",
		Short: "Yeni bir resource oluşturur",
		Long:  "Yeni bir resource (kaynak) oluşturur. Resource, policy ve repository dosyalarını oluşturur.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			makeResource(args[0])
		},
	}
}

// newMakePageCommand, make:page komutunu oluşturur.
func newMakePageCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:page <name>",
		Short: "Yeni bir sayfa oluşturur",
		Long:  "Yeni bir sayfa oluşturur. Sayfa tanımı ve konfigürasyonunu içerir.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			makePage(args[0])
		},
	}
}

// newMakeModelCommand, make:model komutunu oluşturur.
func newMakeModelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "make:model <name>",
		Short: "Yeni bir model oluşturur",
		Long:  "Yeni bir model (veri modeli) oluşturur. Model tanımı ve veri yapısını içerir.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			makeModel(args[0])
		},
	}
}

// newPublishStubsCommand, publish:stubs komutunu oluşturur.
func newPublishStubsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "publish:stubs",
		Short: "Stub dosyalarını .panel/stubs/ dizinine kopyalar",
		Long:  "SDK'daki stub dosyalarını kullanıcının projesine kopyalar.",
		Run: func(cmd *cobra.Command, args []string) {
			publishStubs()
		},
	}
}

// newPublishSkillsCommand, publish:skills komutunu oluşturur.
func newPublishSkillsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "publish:skills",
		Short: "Skill dosyalarını .claude/skills/ dizinine kopyalar",
		Long:  "SDK'daki skill dosyalarını kullanıcının projesine kopyalar.",
		Run: func(cmd *cobra.Command, args []string) {
			publishSkills()
		},
	}
}

// newInitCommand, init komutunu oluşturur.
func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Projeyi başlatır (stubs + skills + starter files)",
		Long:  "Yeni bir Panel.go projesini başlatır. Starter dosyaları, stub ve skill dosyalarını oluşturur.",
		Run: func(cmd *cobra.Command, args []string) {
			database, _ := cmd.Flags().GetString("database")
			initProject(database)
		},
	}
	cmd.Flags().StringP("database", "d", "", "Database driver (sqlite, postgres, mysql)")
	return cmd
}

// makeResource, yeni bir resource (kaynak) oluşturur.
func makeResource(name string) {
	// İsim normalizasyonu
	caser := cases.Title(language.English)
	resourceName := caser.String(name)        // Blog
	packageName := strings.ToLower(name)      // blog
	identifier := strings.ToLower(name) + "s" // blogs
	label := resourceName + "s"               // Blogs
	modelName := resourceName                 // Blog

	// Dizin: internal/resource/<name>
	dir := filepath.Join("internal", "resource", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Şablonlar için veri
	data := map[string]string{
		"PackageName":  packageName,
		"ResourceName": resourceName,
		"ModelName":    modelName,
		"Slug":         identifier,
		"Title":        label,
		"Label":        label,
		"Identifier":   identifier,
		"Group":        "Content",
		"Icon":         "circle",
	}

	// İşlenecek stub'lar
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

// makePage, yeni bir sayfa (page) oluşturur.
func makePage(name string) {
	// İsim normalizasyonu
	caser := cases.Title(language.English)
	pageName := caser.String(name)       // Dashboard
	packageName := strings.ToLower(name) // dashboard
	slug := strings.ToLower(name)        // dashboard
	title := pageName                    // Dashboard

	dir := filepath.Join("internal", "page")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, fmt.Sprintf("%s.go", packageName))

	// Şablonlar için veri
	data := map[string]string{
		"PackageName": "page",
		"PageName":    pageName,
		"Slug":        slug,
		"Title":       title,
		"Group":       "System",
		"Icon":        "circle",
	}

	createFileFromStub("page.stub", targetPath, data)
	fmt.Printf("Page %s generated successfully at %s\n", pageName, targetPath)
}

// makeModel, yeni bir model (veri modeli) oluşturur.
func makeModel(name string) {
	// İsim normalizasyonu
	caser := cases.Title(language.English)
	modelName := caser.String(name)      // Blog
	packageName := strings.ToLower(name) // blog

	// Dizin: internal/domain/<name>
	dir := filepath.Join("internal", "domain", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, "entity.go")

	// Şablonlar için veri
	data := map[string]string{
		"PackageName": packageName,
		"ModelName":   modelName,
	}

	createFileFromStub("model.stub", targetPath, data)
	fmt.Printf("Model %s generated successfully at %s\n", modelName, targetPath)
}

// createFileFromStub, stub dosyasından şablon işleyerek yeni bir dosya oluşturur.
func createFileFromStub(stubName, targetPath string, data map[string]string) {
	// Stub dosyasını gömülü dosya sisteminden oku
	path := stubName
	if !strings.HasPrefix(path, "stubs/") {
		path = filepath.Join("stubs", stubName)
	}

	content, err := stubsFS.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading stub %s: %v\n", path, err)
		return
	}

	// Şablonu işle
	tmpl, err := template.New(stubName).Parse(string(content))
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", stubName, err)
		return
	}

	// Dosya oluştur
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

// publishStubs, SDK'daki stub dosyalarını kullanıcının projesine kopyalar.
func publishStubs() {
	targetDir := filepath.Join(".panel", "stubs")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Stub dosyalarını listele
	stubs := []string{
		"model.stub",
		"resource.stub",
		"policy.stub",
		"repository.stub",
		"page.stub",
		"field_resolver.stub",
		"card_resolver.stub",
	}

	for _, stub := range stubs {
		sourcePath := filepath.Join("stubs", stub)
		content, err := stubsFS.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("Error reading stub %s: %v\n", stub, err)
			continue
		}

		targetPath := filepath.Join(targetDir, stub)
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			fmt.Printf("Error writing stub %s: %v\n", stub, err)
			continue
		}

		fmt.Printf("✓ Copied: %s\n", targetPath)
	}

	fmt.Println("\n✅ Stubs published successfully to .panel/stubs/")
	fmt.Println("You can now customize these stubs for your project.")
}

// publishSkills, SDK'daki skill dosyalarını kullanıcının projesine kopyalar.
func publishSkills() {
	targetDir := ".claude/skills"

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Skill dizinlerini listele
	skills := []string{
		"panel-go-resource",
		"panel-go-field-resolver",
		"panel-go-policy",
		"panel-go-relationship",
		"panel-go-migration",
	}

	for _, skill := range skills {
		// Skill dizinini oluştur
		skillTargetDir := filepath.Join(targetDir, skill)
		if err := os.MkdirAll(skillTargetDir, 0755); err != nil {
			fmt.Printf("Error creating skill directory %s: %v\n", skill, err)
			continue
		}

		// SKILL.md dosyasını kopyala
		sourcePath := filepath.Join("skills", skill, "SKILL.md")
		content, err := skillsFS.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("Error reading skill %s: %v\n", skill, err)
			continue
		}

		targetPath := filepath.Join(skillTargetDir, "SKILL.md")
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			fmt.Printf("Error writing skill %s: %v\n", skill, err)
			continue
		}

		fmt.Printf("✓ Copied: %s\n", targetPath)
	}

	fmt.Println("\n✅ Skills published successfully to .claude/skills/")
	fmt.Println("Claude Code will automatically load these skills.")
}

// initProject, yeni bir Panel.go projesini başlatır.
func initProject(database string) {
	fmt.Println("🚀 Initializing Panel.go project...\n")

	// Proje adını al (mevcut dizin adı)
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}
	projectName := filepath.Base(cwd)

	// Database seçimi (flag yoksa kullanıcıya sor)
	if database == "" {
		database = promptDatabaseSelection()
	}

	// Database'i normalize et
	database = strings.ToLower(strings.TrimSpace(database))
	if database != "sqlite" && database != "postgres" && database != "mysql" {
		fmt.Printf("⚠️  Invalid database driver: %s, using sqlite\n", database)
		database = "sqlite"
	}

	fmt.Printf("📦 Creating project files (database: %s)...\n", database)
	createProjectFiles(projectName, database)

	fmt.Println("\n📦 Publishing stubs...")
	publishStubs()

	fmt.Println("\n🎯 Publishing skills...")
	publishSkills()

	fmt.Println("\n✅ Project initialized successfully!")
	fmt.Println("\nProject structure:")
	fmt.Println("  ├── main.go           # Application entry point")
	fmt.Println("  ├── go.mod            # Go module definition")
	fmt.Println("  ├── .env              # Environment configuration")
	fmt.Println("  ├── .panel/stubs/     # Code generation templates")
	fmt.Println("  └── .claude/skills/   # Claude Code skills")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Update .env with your configuration")
	fmt.Println("  2. Run: go mod tidy")
	fmt.Println("  3. Run: go run main.go")
	fmt.Println("  4. Create a resource: panel make:resource blog")
	fmt.Println("  5. Use Claude Code skills with /panel-go-resource")
}

// promptDatabaseSelection, kullanıcıya database seçimi için interactive prompt gösterir.
func promptDatabaseSelection() string {
	fmt.Println("Select database driver:")
	fmt.Println("  1. SQLite (default, file-based)")
	fmt.Println("  2. PostgreSQL (recommended for production)")
	fmt.Println("  3. MySQL")
	fmt.Print("\nEnter choice [1-3] (default: 1): ")

	var choice string
	fmt.Scanln(&choice)

	switch strings.TrimSpace(choice) {
	case "2":
		return "postgres"
	case "3":
		return "mysql"
	default:
		return "sqlite"
	}
}

// createProjectFiles, proje başlangıç dosyalarını oluşturur.
func createProjectFiles(projectName, database string) {
	// COOKIE_ENCRYPTION_KEY oluştur (openssl rand -base64 32)
	encryptionKey, err := generateEncryptionKey()
	if err != nil {
		fmt.Printf("Warning: Failed to generate encryption key: %v\n", err)
		encryptionKey = "PLEASE-GENERATE-YOUR-OWN-KEY-WITH-OPENSSL"
	}

	// main.go oluştur (database'e göre)
	mainData := map[string]string{
		"ProjectName": projectName,
		"Database":    database,
	}

	// Database'e göre farklı stub kullan
	var mainStub string
	switch database {
	case "postgres":
		mainStub = "main-postgres.stub"
	case "mysql":
		mainStub = "main-mysql.stub"
	default:
		mainStub = "main.stub" // SQLite
	}

	// Eğer database-specific stub yoksa, generic stub kullan
	if _, err := stubsFS.ReadFile(filepath.Join("stubs", mainStub)); err != nil {
		mainStub = "main.stub"
		mainData["DatabaseDriver"] = database
	}

	createFileFromStub(mainStub, "main.go", mainData)

	// go.mod oluştur
	modData := map[string]string{
		"ModuleName": projectName,
	}
	createFileFromStub("go.mod.stub", "go.mod", modData)

	// .env oluştur (database'e göre)
	envData := map[string]string{
		"ProjectName":   projectName,
		"EncryptionKey": encryptionKey,
		"Database":      database,
	}
	createFileFromStub("env.stub", ".env", envData)

	// permissions.toml oluştur
	permissionsContent, err := stubsFS.ReadFile("stubs/permissions.toml.stub")
	if err != nil {
		fmt.Printf("Error reading permissions.toml.stub: %v\n", err)
	} else {
		if err := os.WriteFile("permissions.toml", permissionsContent, 0644); err != nil {
			fmt.Printf("Error creating permissions.toml: %v\n", err)
		} else {
			fmt.Printf("Created: permissions.toml\n")
		}
	}

	// .gitignore oluştur (eğer yoksa)
	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		gitignoreContent := `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
*.db

# Test binary
*.test

# Output
*.out

# Go workspace file
go.work

# Environment
.env

# Storage
storage/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~
`
		if err := os.WriteFile(".gitignore", []byte(gitignoreContent), 0644); err != nil {
			fmt.Printf("Error creating .gitignore: %v\n", err)
		} else {
			fmt.Printf("Created: .gitignore\n")
		}
	}
}

// generateEncryptionKey, openssl kullanarak 32-byte encryption key oluşturur.
func generateEncryptionKey() (string, error) {
	cmd := exec.Command("openssl", "rand", "-base64", "32")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
