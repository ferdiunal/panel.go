// Package plugin, Panel.go plugin sistemi için UI build işlemlerini sağlar.
//
// Bu paket, web-ui clone, dependency yükleme, build ve output kopyalama
// işlemlerini içerir.
package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// WebUIRepoURL, panel.web repository URL'si
	WebUIRepoURL = "https://github.com/ferdiunal/panel.web"

	// WebUIDefaultPath, web-ui'nin varsayılan yolu
	WebUIDefaultPath = "web-ui"

	// AssetsOutputPath, build output'unun kopyalanacağı yol
	AssetsOutputPath = "assets/ui"
)

// BuildUI, UI build alır.
//
// Bu fonksiyon, web-ui'yi build eder ve output'u assets/ui/'ye kopyalar.
// İlk kez çalıştırıldığında web-ui'yi clone eder.
//
// ## Parametreler
//   - opts: Build seçenekleri (DevMode, WatchMode)
//
// ## Dönüş Değeri
//   - error: Build hatası varsa hata, aksi takdirde nil
//
// ## Davranış
//  1. web-ui var mı kontrol eder (yoksa clone eder)
//  2. Package manager detect eder (pnpm > npm)
//  3. Dependencies yükler
//  4. Build alır (production/development)
//  5. Output'u assets/ui/'ye kopyalar
//
// ## Kullanım Örneği
//
//	if err := BuildUI(BuildUIOptions{}); err != nil {
//	    log.Fatal(err)
//	}
func BuildUI(opts BuildUIOptions) error {
	webUIPath := WebUIDefaultPath

	// web-ui var mı kontrol et
	if _, err := os.Stat(webUIPath); os.IsNotExist(err) {
		fmt.Println("web-ui bulunamadı, clone ediliyor...")
		if err := cloneWebUI(webUIPath); err != nil {
			return fmt.Errorf("web-ui clone edilemedi: %w", err)
		}
		fmt.Printf("✓ web-ui clone edildi: %s\n", webUIPath)
	}

	// Package manager detect et
	pkgManager, err := detectPackageManager(webUIPath)
	if err != nil {
		return fmt.Errorf("package manager detect edilemedi: %w", err)
	}

	fmt.Printf("✓ Package manager: %s\n", pkgManager)

	// Dependencies yükle
	fmt.Println("✓ Dependencies yükleniyor...")
	if err := runCommand(webUIPath, pkgManager, "install"); err != nil {
		return fmt.Errorf("dependencies yüklenemedi: %w", err)
	}

	fmt.Println("✓ Dependencies yüklendi")

	// Build al
	if opts.WatchMode {
		// Watch mode: Continuous build
		fmt.Println("✓ Watch mode başlatılıyor...")
		fmt.Println("  (Ctrl+C ile durdurun)")
		return runCommand(webUIPath, pkgManager, "run", "dev")
	}

	// Production/Development build
	buildMode := "build"
	if opts.DevMode {
		buildMode = "build:dev"
	}

	fmt.Printf("✓ Build alınıyor (%s)...\n", buildMode)
	if err := runCommand(webUIPath, pkgManager, "run", buildMode); err != nil {
		return fmt.Errorf("build hatası: %w", err)
	}

	fmt.Println("✓ Build tamamlandı")

	// Output'u kopyala
	distPath := filepath.Join(webUIPath, "dist")
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		return fmt.Errorf("build output bulunamadı: %s", distPath)
	}

	fmt.Printf("✓ Output kopyalanıyor: %s -> %s\n", distPath, AssetsOutputPath)
	if err := copyDir(distPath, AssetsOutputPath); err != nil {
		return fmt.Errorf("output kopyalanamadı: %w", err)
	}

	fmt.Printf("✓ Output kopyalandı: %s\n", AssetsOutputPath)

	return nil
}

// cloneWebUI, panel.web repository'sini clone eder.
//
// Bu fonksiyon, GitHub'dan panel.web repository'sini belirtilen yola clone eder.
//
// ## Parametreler
//   - targetPath: Clone edilecek yol
//
// ## Dönüş Değeri
//   - error: Clone hatası varsa hata, aksi takdirde nil
func cloneWebUI(targetPath string) error {
	// Git clone komutu
	cmd := exec.Command("git", "clone", WebUIRepoURL, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone hatası: %w", err)
	}

	return nil
}

// detectPackageManager, package manager'ı detect eder.
//
// Bu fonksiyon, web-ui dizininde hangi package manager'ın kullanıldığını
// detect eder. Öncelik sırası: pnpm > npm
//
// ## Parametreler
//   - webUIPath: web-ui dizini
//
// ## Dönüş Değeri
//   - string: Package manager adı ("pnpm" veya "npm")
//   - error: Detect hatası varsa hata, aksi takdirde nil
func detectPackageManager(webUIPath string) (string, error) {
	// pnpm-lock.yaml var mı kontrol et
	pnpmLock := filepath.Join(webUIPath, "pnpm-lock.yaml")
	if _, err := os.Stat(pnpmLock); err == nil {
		return "pnpm", nil
	}

	// package-lock.json var mı kontrol et
	npmLock := filepath.Join(webUIPath, "package-lock.json")
	if _, err := os.Stat(npmLock); err == nil {
		return "npm", nil
	}

	// Varsayılan: pnpm
	return "pnpm", nil
}

// runCommand, komut çalıştırır.
//
// Bu fonksiyon, belirtilen dizinde komut çalıştırır ve output'u gösterir.
//
// ## Parametreler
//   - dir: Çalışma dizini
//   - name: Komut adı
//   - args: Komut argümanları
//
// ## Dönüş Değeri
//   - error: Komut hatası varsa hata, aksi takdirde nil
func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("komut hatası (%s %s): %w", name, strings.Join(args, " "), err)
	}

	return nil
}

// copyDir, dizini kopyalar.
//
// Bu fonksiyon, kaynak dizini hedef dizine kopyalar. Mevcut hedef dizini
// varsa önce siler.
//
// ## Parametreler
//   - src: Kaynak dizin
//   - dst: Hedef dizin
//
// ## Dönüş Değeri
//   - error: Kopyalama hatası varsa hata, aksi takdirde nil
func copyDir(src, dst string) error {
	// Hedef dizini sil (varsa)
	if _, err := os.Stat(dst); err == nil {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("hedef dizin silinemedi: %w", err)
		}
	}

	// Hedef dizini oluştur
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("hedef dizin oluşturulamadı: %w", err)
	}

	// Kaynak dizini oku
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("kaynak dizin okunamadı: %w", err)
	}

	// Her entry'yi kopyala
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Dizin: Recursive kopyala
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Dosya: Kopyala
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile, dosyayı kopyalar.
//
// Bu fonksiyon, kaynak dosyayı hedef dosyaya kopyalar.
//
// ## Parametreler
//   - src: Kaynak dosya
//   - dst: Hedef dosya
//
// ## Dönüş Değeri
//   - error: Kopyalama hatası varsa hata, aksi takdirde nil
func copyFile(src, dst string) error {
	// Kaynak dosyayı aç
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("kaynak dosya açılamadı: %w", err)
	}
	defer srcFile.Close()

	// Hedef dosyayı oluştur
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("hedef dosya oluşturulamadı: %w", err)
	}
	defer dstFile.Close()

	// Kopyala
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("dosya kopyalanamadı: %w", err)
	}

	// Permissions kopyala
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("kaynak dosya bilgisi alınamadı: %w", err)
	}

	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("dosya izinleri ayarlanamadı: %w", err)
	}

	return nil
}
