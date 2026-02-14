// Bu paket, Panel CLI uygulamasının ana giriş noktasıdır.
//
// Panel, Go tabanlı bir kod oluşturma aracıdır (code generator) ve aşağıdaki
// komutları destekler:
//   - make:resource: Yeni bir resource (kaynak) oluşturur
//   - make:lens: Resource için yeni bir lens oluşturur
//   - make:action: Resource için yeni bir action oluşturur
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
//go:embed stubs/*.stub stubs/*.yaml
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
	rootCmd.AddCommand(newMakeLensCommand())
	rootCmd.AddCommand(newMakeActionCommand())
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
		Long:  "Yeni bir resource (kaynak) oluşturur. Resource, policy, repository, field resolver ve card resolver dosyalarını oluşturur.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			makeResource(args[0])
		},
	}
}

// newMakeLensCommand, make:lens komutunu oluşturur.
func newMakeLensCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:lens <name>",
		Short: "Belirli bir resource için yeni bir lens oluşturur",
		Long:  "Belirli bir resource package'i için lens dosyası oluşturur. Resource adı --resource flag'i ile verilir veya etkileşimli olarak sorulur.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resourceName, _ := cmd.Flags().GetString("resource")
			if strings.TrimSpace(resourceName) == "" {
				resourceName = promptRequiredInput("Resource name")
			}
			makeLens(args[0], resourceName)
		},
	}

	cmd.Flags().StringP("resource", "r", "", "Hedef resource package adı (örn: blog)")
	return cmd
}

// newMakeActionCommand, make:action komutunu oluşturur.
func newMakeActionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:action <name>",
		Short: "Belirli bir resource için yeni bir action oluşturur",
		Long:  "Belirli bir resource package'i için action dosyası oluşturur. Resource adı --resource flag'i ile verilir veya etkileşimli olarak sorulur.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resourceName, _ := cmd.Flags().GetString("resource")
			if strings.TrimSpace(resourceName) == "" {
				resourceName = promptRequiredInput("Resource name")
			}
			makeAction(args[0], resourceName)
		},
	}

	cmd.Flags().StringP("resource", "r", "", "Hedef resource package adı (örn: blog)")
	return cmd
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

// getModulePath, go.mod dosyasından module path'ini okur.
func getModulePath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Printf("Warning: Could not read go.mod: %v\n", err)
		return "your-module-path"
	}

	// "module " ile başlayan satırı bul
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	return "your-module-path"
}

// makeResource, yeni bir resource (kaynak) oluşturur.
// Model dosyası da otomatik olarak internal/domain/<name>/entity.go'da oluşturulur.
func makeResource(name string) {
	// İsim normalizasyonu
	caser := cases.Title(language.English)
	resourceName := caser.String(name)        // Blog
	packageName := strings.ToLower(name)      // blog
	identifier := strings.ToLower(name) + "s" // blogs
	label := resourceName + "s"               // Blogs
	modelName := resourceName                 // Blog
	tableName := identifier                   // blogs (plural tablo ismi)

	// Module path'ini al (go.mod'dan)
	modulePath := getModulePath()

	// 1. Model dosyasını oluştur: internal/entity/entity.go
	modelDir := filepath.Join("internal", "entity")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		fmt.Printf("Error creating model directory: %v\n", err)
		return
	}

	modelPath := filepath.Join(modelDir, "entity.go")
	modelData := map[string]string{
		"PackageName": "entity",
		"ModelName":   modelName,
		"TableName":   tableName,
	}

	// Model dosyası zaten varsa append et, yoksa oluştur
	if _, err := os.Stat(modelPath); err == nil {
		appendFileFromStub("model_struct.stub", modelPath, modelData)
		fmt.Printf("Modified: %s (appended struct)\n", modelPath)
	} else {
		createFileFromStub("model.stub", modelPath, modelData)
	}

	// 2. Resource dosyalarını oluştur: internal/resource/<name>/
	dir := filepath.Join("internal", "resource", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating resource directory: %v\n", err)
		return
	}

	// Şablonlar için veri
	data := map[string]string{
		"PackageName":     packageName,
		"ResourceName":    resourceName,
		"ModelName":       modelName,
		"ModelPkg":        "entity",                        // Model package adı
		"ModulePath":      modulePath,                      // go.mod'dan okunan module path
		"ModelImportPath": modulePath + "/internal/entity", // Model import path
		"Slug":            identifier,
		"Title":           label,
		"Label":           label,
		"Identifier":      identifier,
		"Group":           "Content",
		"Icon":            "circle",
		"TableName":       tableName,
	}

	// İşlenecek stub'lar
	stubs := map[string]string{
		"resource.stub":       filepath.Join(dir, fmt.Sprintf("%s_resource.go", packageName)),
		"policy.stub":         filepath.Join(dir, fmt.Sprintf("%s_policy.go", packageName)),
		"repository.stub":     filepath.Join(dir, fmt.Sprintf("%s_repository.go", packageName)),
		"field_resolver.stub": filepath.Join(dir, fmt.Sprintf("%s_field_resolver.go", packageName)),
		"card_resolver.stub":  filepath.Join(dir, fmt.Sprintf("%s_card_resolver.go", packageName)),
	}

	for stub, target := range stubs {
		createFileFromStub(stub, target, data)
	}

	fmt.Printf("\n✅ Resource %s generated successfully!\n", resourceName)
	fmt.Printf("   Model:    %s\n", modelPath)
	fmt.Printf("   Resource: %s\n", dir)
	fmt.Printf("   Table:    %s (plural)\n", tableName)
	fmt.Printf("   Import:   %s\n", modulePath+"/internal/domain/"+packageName)

	// 3. main.go dosyasına import ekle
	importPath := modulePath + "/internal/resource/" + packageName
	addImportToMain("main.go", importPath)
}

func normalizeClassName(name string) string {
	normalized := strings.NewReplacer("-", " ", "_", " ").Replace(strings.TrimSpace(name))
	caser := cases.Title(language.English)
	return strings.ReplaceAll(caser.String(normalized), " ", "")
}

func normalizeSlug(name string) string {
	slug := strings.NewReplacer("_", "-", " ", "-").Replace(strings.ToLower(strings.TrimSpace(name)))
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}

func normalizeFileName(name string) string {
	file := strings.NewReplacer("-", "_", " ", "_").Replace(strings.ToLower(strings.TrimSpace(name)))
	for strings.Contains(file, "__") {
		file = strings.ReplaceAll(file, "__", "_")
	}
	return strings.Trim(file, "_")
}

// makeLens, belirli bir resource için lens dosyası oluşturur.
func makeLens(name, resourceName string) {
	lensName := normalizeClassName(name)
	resourcePkg := strings.ToLower(strings.TrimSpace(resourceName))
	lensSlug := normalizeSlug(name)
	fileBase := normalizeFileName(name)

	if lensName == "" || resourcePkg == "" || fileBase == "" {
		fmt.Println("Error: Invalid lens name or resource name")
		return
	}

	dir := filepath.Join("internal", "resource", resourcePkg)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating resource directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, fmt.Sprintf("%s_lens.go", fileBase))
	data := map[string]string{
		"PackageName":  resourcePkg,
		"ResourceSlug": resourcePkg,
		"LensName":     lensName,
		"LensSlug":     lensSlug,
	}

	createFileFromStub("lens.stub", targetPath, data)
	fmt.Printf("\n✅ Lens %s generated successfully for resource %s\n", lensName, resourcePkg)
	fmt.Printf("   File: %s\n", targetPath)
}

// makeAction, belirli bir resource için action dosyası oluşturur.
func makeAction(name, resourceName string) {
	actionName := normalizeClassName(name)
	resourcePkg := strings.ToLower(strings.TrimSpace(resourceName))
	actionSlug := normalizeSlug(name)
	fileBase := normalizeFileName(name)

	if actionName == "" || resourcePkg == "" || fileBase == "" {
		fmt.Println("Error: Invalid action name or resource name")
		return
	}

	dir := filepath.Join("internal", "resource", resourcePkg)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating resource directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, fmt.Sprintf("%s_action.go", fileBase))
	data := map[string]string{
		"PackageName":  resourcePkg,
		"ResourceSlug": resourcePkg,
		"ActionName":   actionName,
		"ActionSlug":   actionSlug,
	}

	createFileFromStub("action.stub", targetPath, data)
	fmt.Printf("\n✅ Action %s generated successfully for resource %s\n", actionName, resourcePkg)
	fmt.Printf("   File: %s\n", targetPath)
}

// addImportToMain, main.go dosyasına anonymous import ekler.
func addImportToMain(mainPath, importPath string) {
	content, err := os.ReadFile(mainPath)
	if err != nil {
		// main.go bulunamazsa sessizce geç (belki proje kök dizininde değiliz)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	imported := false
	inImportBlock := false
	targetImport := fmt.Sprintf("\t_ \"%s\"", importPath)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Zaten ekli mi kontrol et
		if strings.Contains(line, importPath) {
			imported = true
		}

		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
		}

		// Import bloğunun sonuna ekle (bloğu kapatan parantezden hemen önce)
		if inImportBlock && strings.HasPrefix(trimmed, ")") && !imported {
			newLines = append(newLines, targetImport)
			imported = true
			inImportBlock = false
		}

		newLines = append(newLines, line)
	}

	if imported {
		if err := os.WriteFile(mainPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
			fmt.Printf("Error updating %s: %v\n", mainPath, err)
		} else {
			fmt.Printf("Updated %s with import: %s\n", mainPath, importPath)
		}
	} else {
		fmt.Printf("Warning: Could not automatically add import to %s. Please add: _ \"%s\"\n", mainPath, importPath)
	}
}

// makePage, yeni bir sayfa (page) oluşturur.
func makePage(name string) {
	// İsim normalizasyonu
	caser := cases.Title(language.English)
	pageName := caser.String(name)       // Dashboard
	packageName := strings.ToLower(name) // dashboard
	slug := strings.ToLower(name)        // dashboard
	title := pageName                    // Dashboard

	dir := filepath.Join("internal", "pages")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, fmt.Sprintf("%s.go", packageName))

	// Özel stub varsa kullan (dashboard.stub, settings.stub, account.stub)
	stubName := "page.stub"
	switch slug {
	case "dashboard":
		stubName = "dashboard.stub"
	case "settings":
		stubName = "settings.stub"
	case "account":
		stubName = "account.stub"
	}

	// Module path'ini al
	modulePath := getModulePath()

	// Şablonlar için veri
	data := map[string]string{
		"PackageName": "pages",
		"PageName":    pageName,
		"Slug":        slug,
		"Title":       title,
		"Group":       "System",
		"Icon":        "circle",
		"ModulePath":  modulePath,
	}

	createFileFromStub(stubName, targetPath, data)
	fmt.Printf("Page %s generated successfully at %s\n", pageName, targetPath)
}

// makeModel, yeni bir model (veri modeli) oluşturur.
func makeModel(name string) {
	// İsim normalizasyonu
	caser := cases.Title(language.English)
	modelName := caser.String(name)      // Blog
	packageName := strings.ToLower(name) // blog
	tableName := packageName + "s"       // blogs (plural tablo ismi)

	// Dizin: internal/entity
	dir := filepath.Join("internal", "entity")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, "entity.go")

	// Şablonlar için veri
	data := map[string]string{
		"PackageName": "entity",
		"ModelName":   modelName,
		"TableName":   tableName,
	}

	// Model dosyası zaten varsa append et, yoksa oluştur
	if _, err := os.Stat(targetPath); err == nil {
		appendFileFromStub("model_struct.stub", targetPath, data)
		fmt.Printf("Model %s appended successfully to %s (table: %s)\n", modelName, targetPath, tableName)
	} else {
		createFileFromStub("model.stub", targetPath, data)
		fmt.Printf("Model %s generated successfully at %s (table: %s)\n", modelName, targetPath, tableName)
	}
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

// appendFileFromStub, stub dosyasından şablon işleyerek mevcut dosyanın sonuna ekler.
func appendFileFromStub(stubName, targetPath string, data map[string]string) {
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

	// Dosyayı append modunda aç
	f, err := os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", targetPath, err)
		return
	}
	defer f.Close()

	// Bir satır boşluk ekle
	if _, err := f.WriteString("\n"); err != nil {
		fmt.Printf("Error writing newline to file %s: %v\n", targetPath, err)
	}

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Printf("Error executing template %s: %v\n", stubName, err)
	}
	fmt.Printf("Appended: %s\n", targetPath)
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
		"model_struct.stub",
		"resource.stub",
		"lens.stub",
		"action.stub",
		"policy.stub",
		"repository.stub",
		"page.stub",
		"field_resolver.stub",
		"card_resolver.stub",
		"dashboard.stub",
		"settings.stub",
		"account.stub",
		"i18n-pages-example.yaml",
		"i18n-pages-example-en.yaml",
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

func promptRequiredInput(label string) string {
	for {
		fmt.Printf("%s: ", label)
		var value string
		fmt.Scanln(&value)
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
		fmt.Printf("%s cannot be empty.\n", label)
	}
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
	fmt.Println("🚀 Initializing Panel.go project...")
	fmt.Println()

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

	// Module path'ini al
	modulePath := getModulePath()

	fmt.Printf("📦 Creating project files (database: %s)...\n", database)
	createProjectFiles(projectName, database)

	fmt.Println("\n📄 Creating default pages...")
	createDefaultPages(modulePath)

	fmt.Println("\n🌍 Creating locale files...")
	createLocaleFiles()

	fmt.Println("\n📦 Publishing stubs...")
	publishStubs()

	fmt.Println("\n🎯 Publishing skills...")
	publishSkills()

	fmt.Println("\n✅ Project initialized successfully!")
	fmt.Println("\nProject structure:")
	fmt.Println("  ├── main.go              # Application entry point")
	fmt.Println("  ├── go.mod               # Go module definition")
	fmt.Println("  ├── .env                 # Environment configuration")
	fmt.Println("  ├── internal/pages/      # Custom pages (Dashboard, Settings, Account)")
	fmt.Println("  ├── locales/             # i18n translation files")
	fmt.Println("  ├── .panel/stubs/        # Code generation templates")
	fmt.Println("  └── .claude/skills/      # Claude Code skills")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Update .env with your configuration")
	fmt.Println("  2. Run: go mod tidy")
	fmt.Println("  3. Run: go run main.go")
	fmt.Println("  4. Create a resource: panel make:resource blog")
	fmt.Println("  5. Create a page: panel make:page analytics")
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
	modulePath := getModulePath()
	mainData := map[string]string{
		"ProjectName": projectName,
		"Database":    database,
		"ModulePath":  modulePath,
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

// createDefaultPages, varsayılan Dashboard, Settings ve Account sayfalarını oluşturur.
func createDefaultPages(modulePath string) {
	dir := filepath.Join("internal", "pages")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating pages directory: %v\n", err)
		return
	}

	// Şablonlar için veri
	pages := []struct {
		StubName string
		FileName string
		PageName string
		Slug     string
		Title    string
		Icon     string
	}{
		{"dashboard.stub", "dashboard.go", "Dashboard", "dashboard", "Dashboard", "home"},
		{"settings.stub", "settings.go", "Settings", "settings", "Settings", "settings"},
		{"account.stub", "account.go", "Account", "account", "Account", "user"},
	}

	for _, p := range pages {
		targetPath := filepath.Join(dir, p.FileName)

		// Dosya zaten varsa atla
		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("⏩ Skipped (already exists): %s\n", targetPath)
			continue
		}

		data := map[string]string{
			"PackageName": "pages",
			"PageName":    p.PageName,
			"Slug":        p.Slug,
			"Title":       p.Title,
			"Group":       "System",
			"Icon":        p.Icon,
			"ModulePath":  modulePath,
		}
		createFileFromStub(p.StubName, targetPath, data)
	}
}

// createLocaleFiles, i18n dil dosyalarını locales/ dizinine kopyalar.
func createLocaleFiles() {
	localesDir := "locales"
	if err := os.MkdirAll(localesDir, 0755); err != nil {
		fmt.Printf("Error creating locales directory: %v\n", err)
		return
	}

	// i18n dosyalarını kopyala
	localeFiles := map[string]string{
		"i18n-pages-example.yaml":    filepath.Join(localesDir, "tr.yaml"),
		"i18n-pages-example-en.yaml": filepath.Join(localesDir, "en.yaml"),
	}

	for stub, target := range localeFiles {
		// Dosya zaten varsa atla
		if _, err := os.Stat(target); err == nil {
			fmt.Printf("⏩ Skipped (already exists): %s\n", target)
			continue
		}

		sourcePath := filepath.Join("stubs", stub)
		content, err := stubsFS.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", stub, err)
			continue
		}

		if err := os.WriteFile(target, content, 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", target, err)
			continue
		}
		fmt.Printf("✓ Created: %s\n", target)
	}
}
