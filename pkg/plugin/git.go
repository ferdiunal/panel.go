// Package plugin, Panel.go plugin sistemi için Git işlemlerini sağlar.
//
// Bu paket, Git repository clone ve URL parsing işlemlerini içerir.
package plugin

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneRepository, Git repository'sini clone eder.
//
// Bu fonksiyon, belirtilen Git URL'den repository'yi belirtilen yola clone eder.
//
// ## Parametreler
//   - url: Git repository URL'si
//   - targetPath: Clone edilecek yol
//   - branch: Clone edilecek branch (boşsa default branch)
//
// ## Dönüş Değeri
//   - error: Clone hatası varsa hata, aksi takdirde nil
//
// ## Kullanım Örneği
//
//	if err := CloneRepository("https://github.com/user/repo", "./plugins/repo", "main"); err != nil {
//	    log.Fatal(err)
//	}
func CloneRepository(url, targetPath, branch string) error {
	opts := &git.CloneOptions{
		URL:      url,
		Progress: nil, // Progress gösterme (şimdilik kapalı)
	}

	// Branch belirtilmişse
	if branch != "" {
		opts.ReferenceName = plumbing.NewBranchReferenceName(branch)
	}

	// Clone
	_, err := git.PlainClone(targetPath, false, opts)
	if err != nil {
		return fmt.Errorf("git clone hatası: %w", err)
	}

	return nil
}
