// Package plugin, Panel.go plugin sistemi için CLI komutlarını sağlar.
//
// Bu paket, plugin oluşturma, ekleme, silme, listeleme ve build işlemleri için
// Cobra-based CLI komutları içerir.
package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewPluginCommand, plugin komut grubunun root command'ını oluşturur.
//
// Bu fonksiyon, tüm plugin alt komutlarını (create, add, remove, list, build)
// içeren ana plugin command'ını döndürür.
//
// ## Kullanım
//
//	rootCmd.AddCommand(NewPluginCommand())
//
// ## Alt Komutlar
//   - create: Yeni plugin oluşturur
//   - add: Git repository'den plugin ekler
//   - remove: Plugin'i siler
//   - list: Yüklü plugin'leri listeler
//   - build: UI build alır
func NewPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin yönetimi komutları",
		Long:  "Panel.go plugin'lerini oluşturmak, eklemek, silmek ve yönetmek için komutlar.",
	}

	// Alt komutları ekle
	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newAddCommand())
	cmd.AddCommand(newRemoveCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newBuildCommand())

	return cmd
}

// newCreateCommand, plugin:create komutunu oluşturur.
//
// Bu komut, yeni bir plugin scaffold eder. Backend ve frontend dosyalarını
// oluşturur, workspace config'i günceller ve build alır.
//
// ## Kullanım
//
//	panel plugin create <plugin-name> [flags]
//
// ## Flags
//   - --path: Plugin dizini (default: ./plugins)
//   - --no-frontend: Frontend scaffold etme
//   - --no-build: Otomatik build yapma
func newCreateCommand() *cobra.Command {
	var (
		pluginPath  string
		noFrontend  bool
		noBuild     bool
		withExample bool
	)

	cmd := &cobra.Command{
		Use:   "create <plugin-name>",
		Short: "Yeni plugin oluşturur",
		Long:  "Yeni bir plugin scaffold eder. Backend ve frontend dosyalarını oluşturur.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := args[0]

			fmt.Printf("🚀 Plugin oluşturuluyor: %s\n\n", pluginName)

			// Plugin oluştur
			opts := CreatePluginOptions{
				Name:        pluginName,
				Path:        pluginPath,
				NoFrontend:  noFrontend,
				NoBuild:     noBuild,
				WithExample: withExample,
			}

			if err := CreatePlugin(opts); err != nil {
				return fmt.Errorf("plugin oluşturma hatası: %w", err)
			}

			fmt.Printf("\n✅ Plugin '%s' başarıyla oluşturuldu!\n\n", pluginName)
			fmt.Println("Sonraki adımlar:")
			fmt.Printf("  1. Backend implement et: %s/%s/plugin.go\n", pluginPath, pluginName)
			if !noFrontend {
				fmt.Printf("  2. Frontend field'ları ekle: %s/%s/frontend/fields/\n", pluginPath, pluginName)
			}
			fmt.Printf("  3. Plugin'i import et: import _ \"your-module/%s/%s\"\n", strings.TrimPrefix(pluginPath, "./"), pluginName)
			if !noBuild {
				fmt.Println("  4. Rebuild: panel plugin build")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&pluginPath, "path", "./plugins", "Plugin dizini")
	cmd.Flags().BoolVar(&noFrontend, "no-frontend", false, "Frontend scaffold etme")
	cmd.Flags().BoolVar(&noBuild, "no-build", false, "Otomatik build yapma")
	cmd.Flags().BoolVar(&withExample, "with-example", false, "Tüm relationship türlerini içeren örnek entity'ler ekle")

	return cmd
}

// newAddCommand, plugin:add komutunu oluşturur.
//
// Bu komut, Git repository'den plugin ekler. Repository'yi clone eder,
// validate eder, workspace config'i günceller ve build alır.
//
// ## Kullanım
//
//	panel plugin add <git-url> [flags]
//
// ## Flags
//   - --path: Plugin dizini (default: ./plugins)
//   - --branch: Git branch (default: main)
//   - --no-build: Otomatik build yapma
func newAddCommand() *cobra.Command {
	var (
		pluginPath string
		branch     string
		noBuild    bool
	)

	cmd := &cobra.Command{
		Use:   "add <git-url>",
		Short: "Git repository'den plugin ekler",
		Long:  "Git repository'den plugin clone eder, validate eder ve workspace'e ekler.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gitURL := args[0]

			fmt.Printf("📦 Plugin ekleniyor: %s\n\n", gitURL)

			// Plugin ekle
			opts := AddPluginOptions{
				GitURL:  gitURL,
				Path:    pluginPath,
				Branch:  branch,
				NoBuild: noBuild,
			}

			pluginName, err := AddPlugin(opts)
			if err != nil {
				return fmt.Errorf("plugin ekleme hatası: %w", err)
			}

			fmt.Printf("\n✅ Plugin '%s' başarıyla eklendi!\n\n", pluginName)
			fmt.Println("Sonraki adımlar:")
			fmt.Printf("  1. Plugin'i import et: import _ \"your-module/%s/%s\"\n", strings.TrimPrefix(pluginPath, "./"), pluginName)
			if !noBuild {
				fmt.Println("  2. Rebuild: panel plugin build")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&pluginPath, "path", "./plugins", "Plugin dizini")
	cmd.Flags().StringVar(&branch, "branch", "main", "Git branch")
	cmd.Flags().BoolVar(&noBuild, "no-build", false, "Otomatik build yapma")

	return cmd
}

// newRemoveCommand, plugin:remove komutunu oluşturur.
//
// Bu komut, plugin'i siler. Workspace reference'ı kaldırır, plugin dosyalarını
// siler ve build alır.
//
// ## Kullanım
//
//	panel plugin remove <plugin-name> [flags]
//
// ## Flags
//   - --path: Plugin dizini (default: ./plugins)
//   - --keep-files: Plugin dosyalarını silme
//   - --no-build: Otomatik build yapma
func newRemoveCommand() *cobra.Command {
	var (
		pluginPath string
		keepFiles  bool
		noBuild    bool
	)

	cmd := &cobra.Command{
		Use:   "remove <plugin-name>",
		Short: "Plugin'i siler",
		Long:  "Plugin'i workspace'den kaldırır ve dosyalarını siler.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := args[0]

			fmt.Printf("🗑️  Plugin siliniyor: %s\n\n", pluginName)

			// Plugin sil
			opts := RemovePluginOptions{
				Name:      pluginName,
				Path:      pluginPath,
				KeepFiles: keepFiles,
				NoBuild:   noBuild,
			}

			if err := RemovePlugin(opts); err != nil {
				return fmt.Errorf("plugin silme hatası: %w", err)
			}

			fmt.Printf("\n✅ Plugin '%s' başarıyla silindi!\n", pluginName)

			return nil
		},
	}

	cmd.Flags().StringVar(&pluginPath, "path", "./plugins", "Plugin dizini")
	cmd.Flags().BoolVar(&keepFiles, "keep-files", false, "Plugin dosyalarını silme")
	cmd.Flags().BoolVar(&noBuild, "no-build", false, "Otomatik build yapma")

	return cmd
}

// newListCommand, plugin:list komutunu oluşturur.
//
// Bu komut, yüklü plugin'leri listeler. Plugin metadata'sını okur ve
// tablo formatında gösterir.
//
// ## Kullanım
//
//	panel plugin list [flags]
//
// ## Flags
//   - --path: Plugin dizini (default: ./plugins)
//   - --json: JSON output
func newListCommand() *cobra.Command {
	var (
		pluginPath string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Yüklü plugin'leri listeler",
		Long:  "Yüklü plugin'leri metadata ile birlikte listeler.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Plugin'leri listele
			plugins, err := ListPlugins(pluginPath)
			if err != nil {
				return fmt.Errorf("plugin listeleme hatası: %w", err)
			}

			if len(plugins) == 0 {
				fmt.Println("Yüklü plugin bulunamadı.")
				return nil
			}

			if jsonOutput {
				// JSON output
				return printPluginsJSON(plugins)
			}

			// Tablo output
			return printPluginsTable(plugins)
		},
	}

	cmd.Flags().StringVar(&pluginPath, "path", "./plugins", "Plugin dizini")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "JSON output")

	return cmd
}

// newBuildCommand, plugin:build komutunu oluşturur.
//
// Bu komut, UI build alır. web-ui'yi clone eder (ilk kez), dependencies
// yükler, build alır ve output'u assets/ui/'ye kopyalar.
//
// ## Kullanım
//
//	panel plugin build [flags]
//
// ## Flags
//   - --dev: Development build (no minification)
//   - --watch: Watch mode (continuous build)
func newBuildCommand() *cobra.Command {
	var (
		devMode   bool
		watchMode bool
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "UI build alır",
		Long:  "web-ui'yi build eder ve output'u assets/ui/'ye kopyalar.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔨 UI build alınıyor...")
			fmt.Println()

			// Build UI
			opts := BuildUIOptions{
				DevMode:   devMode,
				WatchMode: watchMode,
			}

			if err := BuildUI(opts); err != nil {
				return fmt.Errorf("build hatası: %w", err)
			}

			if !watchMode {
				fmt.Println("\n✅ Build başarıyla tamamlandı!")
				fmt.Println("\nBuild output: assets/ui/")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&devMode, "dev", false, "Development build (no minification)")
	cmd.Flags().BoolVar(&watchMode, "watch", false, "Watch mode (continuous build)")

	return cmd
}

// printPluginsTable, plugin'leri tablo formatında yazdırır.
func printPluginsTable(plugins []PluginInfo) error {
	fmt.Println("Yüklü Plugin'ler:")
	fmt.Println()
	fmt.Printf("%-20s %-10s %-20s %-10s %-10s\n", "NAME", "VERSION", "AUTHOR", "FRONTEND", "STATUS")
	fmt.Println(strings.Repeat("-", 80))

	for _, p := range plugins {
		frontend := "No"
		if p.HasFrontend {
			frontend = "Yes"
		}

		status := "Active"
		if !p.Valid {
			status = "Invalid"
		}

		fmt.Printf("%-20s %-10s %-20s %-10s %-10s\n",
			truncate(p.Name, 20),
			truncate(p.Version, 10),
			truncate(p.Author, 20),
			frontend,
			status,
		)
	}

	fmt.Printf("\nToplam: %d plugin\n", len(plugins))
	return nil
}

// printPluginsJSON, plugin'leri JSON formatında yazdırır.
func printPluginsJSON(plugins []PluginInfo) error {
	data, err := json.MarshalIndent(plugins, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// truncate, string'i belirtilen uzunlukta keser.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// CreatePluginOptions, plugin oluşturma seçenekleri.
type CreatePluginOptions struct {
	Name        string
	Path        string
	NoFrontend  bool
	NoBuild     bool
	WithExample bool
}

// AddPluginOptions, plugin ekleme seçenekleri.
type AddPluginOptions struct {
	GitURL  string
	Path    string
	Branch  string
	NoBuild bool
}

// RemovePluginOptions, plugin silme seçenekleri.
type RemovePluginOptions struct {
	Name      string
	Path      string
	KeepFiles bool
	NoBuild   bool
}

// BuildUIOptions, UI build seçenekleri.
type BuildUIOptions struct {
	DevMode   bool
	WatchMode bool
}

// PluginInfo, plugin metadata bilgisi.
type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HasFrontend bool   `json:"has_frontend"`
	Valid       bool   `json:"valid"`
	Path        string `json:"path"`
}

// CreatePlugin, yeni plugin oluşturur.
func CreatePlugin(opts CreatePluginOptions) error {
	// Plugin dizini oluştur
	pluginDir := filepath.Join(opts.Path, opts.Name)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("plugin dizini oluşturulamadı: %w", err)
	}

	fmt.Printf("✓ Plugin dizini oluşturuldu: %s\n", pluginDir)

	// Backend dosyaları oluştur
	if err := generateBackendFiles(pluginDir, opts.Name, opts.WithExample); err != nil {
		return fmt.Errorf("backend dosyaları oluşturulamadı: %w", err)
	}

	fmt.Println("✓ Backend dosyaları oluşturuldu: plugin.go, plugin.yaml")

	// Frontend dosyaları oluştur (eğer --no-frontend değilse)
	if !opts.NoFrontend {
		if err := generateFrontendFiles(pluginDir, opts.Name); err != nil {
			return fmt.Errorf("frontend dosyaları oluşturulamadı: %w", err)
		}

		fmt.Println("✓ Frontend dosyaları oluşturuldu: index.ts, package.json, tsconfig.json")
	}

	// web-ui clone (ilk kez)
	webUIPath := "web-ui"
	if _, err := os.Stat(webUIPath); os.IsNotExist(err) {
		fmt.Println("✓ web-ui clone ediliyor...")
		if err := cloneWebUI(webUIPath); err != nil {
			return fmt.Errorf("web-ui clone edilemedi: %w", err)
		}
		fmt.Printf("✓ web-ui clone edildi: %s\n", webUIPath)
	}

	// Workspace config güncelle
	if !opts.NoFrontend {
		if err := updateWorkspaceConfig(webUIPath, opts.Name, pluginDir); err != nil {
			return fmt.Errorf("workspace config güncellenemedi: %w", err)
		}

		fmt.Println("✓ Workspace config güncellendi: web-ui/pnpm-workspace.yaml")

		// Plugin workspace reference oluştur
		if err := createPluginSymlink(webUIPath, opts.Name, pluginDir); err != nil {
			return fmt.Errorf("workspace reference oluşturulamadı: %w", err)
		}

		fmt.Printf("✓ Workspace reference oluşturuldu: web-ui/plugins/%s\n", opts.Name)
	}

	// Build (eğer --no-build değilse)
	if !opts.NoBuild {
		fmt.Println("✓ UI build alınıyor...")
		if err := BuildUI(BuildUIOptions{}); err != nil {
			return fmt.Errorf("build hatası: %w", err)
		}
		fmt.Println("✓ Build tamamlandı: assets/ui/")
	}

	return nil
}

// AddPlugin, Git repository'den plugin ekler.
func AddPlugin(opts AddPluginOptions) (string, error) {
	// Git URL'den plugin adını çıkar
	pluginName, err := parsePluginNameFromGitURL(opts.GitURL)
	if err != nil {
		return "", fmt.Errorf("git URL parse edilemedi: %w", err)
	}

	// Plugin dizini
	pluginDir := filepath.Join(opts.Path, pluginName)

	// Plugin clone
	fmt.Printf("✓ Plugin clone ediliyor: %s\n", opts.GitURL)
	if err := CloneRepository(opts.GitURL, pluginDir, opts.Branch); err != nil {
		return "", fmt.Errorf("plugin clone edilemedi: %w", err)
	}

	fmt.Printf("✓ Plugin clone edildi: %s\n", pluginDir)

	// Plugin validate
	if err := validatePlugin(pluginDir); err != nil {
		return "", fmt.Errorf("plugin geçersiz: %w", err)
	}

	fmt.Println("✓ Plugin validate edildi")

	// web-ui clone (ilk kez)
	webUIPath := "web-ui"
	if _, err := os.Stat(webUIPath); os.IsNotExist(err) {
		fmt.Println("✓ web-ui clone ediliyor...")
		if err := cloneWebUI(webUIPath); err != nil {
			return "", fmt.Errorf("web-ui clone edilemedi: %w", err)
		}
		fmt.Printf("✓ web-ui clone edildi: %s\n", webUIPath)
	}

	// Frontend var mı kontrol et
	frontendPath := filepath.Join(pluginDir, "frontend")
	hasFrontend := false
	if _, err := os.Stat(frontendPath); err == nil {
		hasFrontend = true
	}

	// Workspace config güncelle (eğer frontend varsa)
	if hasFrontend {
		if err := updateWorkspaceConfig(webUIPath, pluginName, pluginDir); err != nil {
			return "", fmt.Errorf("workspace config güncellenemedi: %w", err)
		}

		fmt.Println("✓ Workspace config güncellendi")

		// Plugin workspace reference oluştur
		if err := createPluginSymlink(webUIPath, pluginName, pluginDir); err != nil {
			return "", fmt.Errorf("workspace reference oluşturulamadı: %w", err)
		}

		fmt.Printf("✓ Workspace reference oluşturuldu: web-ui/plugins/%s\n", pluginName)
	}

	// Build (eğer --no-build değilse)
	if !opts.NoBuild {
		fmt.Println("✓ UI build alınıyor...")
		if err := BuildUI(BuildUIOptions{}); err != nil {
			return "", fmt.Errorf("build hatası: %w", err)
		}
		fmt.Println("✓ Build tamamlandı: assets/ui/")
	}

	return pluginName, nil
}

// RemovePlugin, plugin'i siler.
func RemovePlugin(opts RemovePluginOptions) error {
	pluginDir := filepath.Join(opts.Path, opts.Name)

	// Plugin var mı kontrol et
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return fmt.Errorf("plugin bulunamadı: %s", opts.Name)
	}

	// Workspace reference sil
	webUIPath := "web-ui"
	symlinkPath := filepath.Join(webUIPath, "plugins", opts.Name)
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			return fmt.Errorf("workspace reference silinemedi: %w", err)
		}
		fmt.Printf("✓ Workspace reference silindi: %s\n", symlinkPath)
	}

	// Workspace config güncelle
	if err := removeFromWorkspaceConfig(webUIPath, opts.Name); err != nil {
		return fmt.Errorf("workspace config güncellenemedi: %w", err)
	}

	fmt.Println("✓ Workspace config güncellendi")

	// Plugin dosyalarını sil (eğer --keep-files değilse)
	if !opts.KeepFiles {
		if err := os.RemoveAll(pluginDir); err != nil {
			return fmt.Errorf("plugin dosyaları silinemedi: %w", err)
		}
		fmt.Printf("✓ Plugin dosyaları silindi: %s\n", pluginDir)
	}

	// Build (eğer --no-build değilse)
	if !opts.NoBuild {
		fmt.Println("✓ UI build alınıyor...")
		if err := BuildUI(BuildUIOptions{}); err != nil {
			return fmt.Errorf("build hatası: %w", err)
		}
		fmt.Println("✓ Build tamamlandı: assets/ui/")
	}

	return nil
}

// ListPlugins, yüklü plugin'leri listeler.
func ListPlugins(pluginPath string) ([]PluginInfo, error) {
	// Plugin dizini var mı kontrol et
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, nil
	}

	// Plugin dizinlerini oku
	entries, err := os.ReadDir(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("plugin dizini okunamadı: %w", err)
	}

	plugins := []PluginInfo{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginDir := filepath.Join(pluginPath, entry.Name())

		// Plugin metadata oku
		metadata, err := readPluginMetadata(pluginDir)
		if err != nil {
			// Geçersiz plugin, skip
			continue
		}

		// Frontend var mı kontrol et
		frontendPath := filepath.Join(pluginDir, "frontend")
		hasFrontend := false
		if _, err := os.Stat(frontendPath); err == nil {
			hasFrontend = true
		}

		plugins = append(plugins, PluginInfo{
			Name:        metadata.Name,
			Version:     metadata.Version,
			Author:      metadata.Author,
			Description: metadata.Description,
			HasFrontend: hasFrontend,
			Valid:       true,
			Path:        pluginDir,
		})
	}

	return plugins, nil
}

// parsePluginNameFromGitURL, Git URL'den plugin adını çıkarır.
func parsePluginNameFromGitURL(gitURL string) (string, error) {
	// URL'den son path segment'i al
	// Örnek: github.com/user/plugin-name -> plugin-name
	parts := strings.Split(strings.TrimSuffix(gitURL, ".git"), "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("geçersiz git URL: %s", gitURL)
	}

	return parts[len(parts)-1], nil
}

// validatePlugin, plugin'in geçerli olup olmadığını kontrol eder.
func validatePlugin(pluginDir string) error {
	// plugin.yaml var mı kontrol et
	pluginYAML := filepath.Join(pluginDir, "plugin.yaml")
	if _, err := os.Stat(pluginYAML); os.IsNotExist(err) {
		return fmt.Errorf("plugin.yaml bulunamadı")
	}

	// plugin.go var mı kontrol et
	pluginGo := filepath.Join(pluginDir, "plugin.go")
	if _, err := os.Stat(pluginGo); os.IsNotExist(err) {
		return fmt.Errorf("plugin.go bulunamadı")
	}

	return nil
}

// removeFromWorkspaceConfig, workspace config'den plugin'i kaldırır.
func removeFromWorkspaceConfig(webUIPath, pluginName string) error {
	workspaceYAMLPath := filepath.Join(webUIPath, "pnpm-workspace.yaml")
	if _, err := os.Stat(workspaceYAMLPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(workspaceYAMLPath)
	if err != nil {
		return fmt.Errorf("workspace config okunamadı: %w", err)
	}

	var workspaceConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &workspaceConfig); err != nil {
		return fmt.Errorf("workspace config parse edilemedi: %w", err)
	}

	rawPackages, ok := workspaceConfig["packages"]
	if !ok {
		return nil
	}

	packages, ok := rawPackages.([]interface{})
	if !ok {
		return nil
	}

	pluginsDir := filepath.Join(webUIPath, "plugins")
	shouldKeepPluginWorkspacePath := false
	if entries, err := os.ReadDir(pluginsDir); err == nil {
		for _, entry := range entries {
			name := strings.TrimSpace(entry.Name())
			if name == "" || strings.HasPrefix(name, ".") || name == pluginName {
				continue
			}
			shouldKeepPluginWorkspacePath = true
			break
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("plugins dizini okunamadı: %w", err)
	}

	pluginWildcardPath := "../plugins/*/frontend"
	pluginSpecificPath := fmt.Sprintf("../plugins/%s/frontend", pluginName)

	filtered := make([]interface{}, 0, len(packages))
	changed := false

	for _, pkg := range packages {
		pkgStr, ok := pkg.(string)
		if !ok {
			filtered = append(filtered, pkg)
			continue
		}

		if pkgStr == pluginSpecificPath {
			changed = true
			continue
		}

		if pkgStr == pluginWildcardPath && !shouldKeepPluginWorkspacePath {
			changed = true
			continue
		}

		filtered = append(filtered, pkg)
	}

	if !changed {
		return nil
	}

	workspaceConfig["packages"] = filtered

	updatedData, err := yaml.Marshal(workspaceConfig)
	if err != nil {
		return fmt.Errorf("workspace config marshal edilemedi: %w", err)
	}

	if err := os.WriteFile(workspaceYAMLPath, updatedData, 0644); err != nil {
		return fmt.Errorf("workspace config yazılamadı: %w", err)
	}

	return nil
}
