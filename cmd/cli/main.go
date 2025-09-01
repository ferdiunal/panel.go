package main

import (
	"context"
	"fmt"
	"log"
	"os"

	entsql "entgo.io/ent/dialect/sql"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ferdiunal/go.utils/database"
	goutils "github.com/ferdiunal/go.utils/database/interfaces"
	"panel.go/internal/ent"
	"panel.go/internal/repository"
	"panel.go/internal/service"
	"panel.go/shared/encrypt"
)

type step int

const (
	inputName step = iota
	inputEmail
	inputPassword
	creating
	done
)

type model struct {
	step        step
	name        string
	email       string
	password    string
	err         error
	created     bool
	authService *service.AuthService
	db          goutils.DatabaseService
}

func initialModel() *model {
	dbService, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB("postgres", dbService.Db())
	ent := ent.NewClient(ent.Driver(drv))

	// Initialize repositories
	userRepo := repository.NewUserRepository(ent)
	accountRepo := repository.NewAccountRepository(ent)
	sessionRepo := repository.NewSessionRepository(ent)

	// Initialize encrypt service
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	fmt.Println(encryptionKey)
	encryptService := encrypt.NewCrypt(encryptionKey)

	// Initialize auth service
	authService := service.NewAuthService(accountRepo, userRepo, sessionRepo, encryptService)

	return &model{
		step:        inputName,
		authService: authService,
		db:          dbService,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle("User Creator")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.step {
			case inputName:
				if m.name != "" {
					m.step = inputEmail
				}
			case inputEmail:
				if m.email != "" {
					m.step = inputPassword
				}
			case inputPassword:
				if m.password != "" {
					m.step = creating
					return m, m.createUser()
				}
			case done:
				return m, tea.Quit
			}
		case "backspace":
			if m.step == inputName && len(m.name) > 0 {
				m.name = m.name[:len(m.name)-1]
			} else if m.step == inputEmail && len(m.email) > 0 {
				m.email = m.email[:len(m.email)-1]
			} else if m.step == inputPassword && len(m.password) > 0 {
				m.password = m.password[:len(m.password)-1]
			}
		default:
			// Sadece yazdırılabilir karakterleri kabul et
			if len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] <= 126 {
				switch m.step {
				case inputName:
					m.name += msg.String()
				case inputEmail:
					m.email += msg.String()
				case inputPassword:
					m.password += msg.String()
				}
			}
		}
	case createUserMsg:
		m.step = done
		m.created = true
		m.err = msg.err
	}

	return m, nil
}

func (m *model) View() string {
	switch m.step {
	case inputName:
		cursor := ""
		if len(m.name) == 0 {
			cursor = "_"
		}
		return "=== User Creator ===\n\n" +
			"Name: " + m.name + cursor + "\n\n" +
			"Press Enter to continue..."

	case inputEmail:
		cursor := ""
		if len(m.email) == 0 {
			cursor = "_"
		}
		return "=== User Creator ===\n\n" +
			"Name: " + m.name + "\n" +
			"Email: " + m.email + cursor + "\n\n" +
			"Press Enter to continue..."

	case inputPassword:
		cursor := ""
		if len(m.password) == 0 {
			cursor = "_"
		}
		// Şifreyi * ile gizle
		maskedPassword := ""
		for i := 0; i < len(m.password); i++ {
			maskedPassword += "*"
		}
		return "=== User Creator ===\n\n" +
			"Name: " + m.name + "\n" +
			"Email: " + m.email + "\n" +
			"Password: " + maskedPassword + cursor + "\n\n" +
			"Press Enter to create user..."

	case creating:
		return "=== User Creator ===\n\n" +
			"Creating user...\n" +
			"Name: " + m.name + "\n" +
			"Email: " + m.email + "\n\n" +
			"Please wait..."

	case done:
		if m.err != nil {
			return "=== User Creator ===\n\n" +
				"ERROR: " + m.err.Error() + "\n\n" +
				"Press Enter to exit..."
		}
		return "=== User Creator ===\n\n" +
			"SUCCESS: User created!\n" +
			"Name: " + m.name + "\n" +
			"Email: " + m.email + "\n\n" +
			"Press Enter to exit..."
	}

	return ""
}

type createUserMsg struct {
	err error
}

func (m *model) createUser() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := m.authService.RegisterCLI(ctx, m.name, m.email, m.password)
		return createUserMsg{err: err}
	}
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
