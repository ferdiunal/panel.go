# Secrets Management Guide

## Overview

This guide provides best practices and implementation patterns for managing secrets in the panel.go application. Proper secrets management is critical for security and compliance.

---

## Current State

The application currently uses environment variables for configuration, including sensitive data like:
- Database credentials
- Encryption keys
- API keys
- Session secrets

**Security Risk**: Environment variables can be exposed through process listings, logs, and error messages.

---

## Recommended Approach

### Option 1: HashiCorp Vault (Production-Grade)

**Best for**: Production environments, enterprise deployments, compliance requirements

#### Setup

```bash
# Install Vault
brew install vault  # macOS
# or download from https://www.vaultproject.io/downloads

# Start Vault dev server (for testing)
vault server -dev

# Set Vault address
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN='<dev-token>'
```

#### Store Secrets

```bash
# Store database credentials
vault kv put secret/panel/database \
  host=localhost \
  port=5432 \
  username=panel_user \
  password=secure_password \
  database=panel_db

# Store encryption key
vault kv put secret/panel/encryption \
  key=<64-char-hex-key>

# Store session secret
vault kv put secret/panel/session \
  secret=<random-secret>
```

#### Go Integration

```go
package config

import (
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Address string
	Token   string
}

type SecretsManager struct {
	client *vault.Client
}

func NewSecretsManager(config VaultConfig) (*SecretsManager, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = config.Address

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(config.Token)

	return &SecretsManager{client: client}, nil
}

func (sm *SecretsManager) GetDatabaseCredentials() (map[string]string, error) {
	secret, err := sm.client.Logical().Read("secret/data/panel/database")
	if err != nil {
		return nil, fmt.Errorf("failed to read database credentials: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no database credentials found")
	}

	data := secret.Data["data"].(map[string]interface{})
	credentials := make(map[string]string)
	for k, v := range data {
		credentials[k] = v.(string)
	}

	return credentials, nil
}

func (sm *SecretsManager) GetEncryptionKey() (string, error) {
	secret, err := sm.client.Logical().Read("secret/data/panel/encryption")
	if err != nil {
		return "", fmt.Errorf("failed to read encryption key: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no encryption key found")
	}

	data := secret.Data["data"].(map[string]interface{})
	return data["key"].(string), nil
}

func (sm *SecretsManager) GetSessionSecret() (string, error) {
	secret, err := sm.client.Logical().Read("secret/data/panel/session")
	if err != nil {
		return "", fmt.Errorf("failed to read session secret: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("no session secret found")
	}

	data := secret.Data["data"].(map[string]interface{})
	return data["secret"].(string), nil
}
```

#### Usage in Application

```go
// In main.go or initialization
vaultConfig := config.VaultConfig{
	Address: os.Getenv("VAULT_ADDR"),
	Token:   os.Getenv("VAULT_TOKEN"),
}

secretsManager, err := config.NewSecretsManager(vaultConfig)
if err != nil {
	log.Fatalf("Failed to initialize secrets manager: %v", err)
}

// Get database credentials
dbCreds, err := secretsManager.GetDatabaseCredentials()
if err != nil {
	log.Fatalf("Failed to get database credentials: %v", err)
}

// Get encryption key
encryptionKey, err := secretsManager.GetEncryptionKey()
if err != nil {
	log.Fatalf("Failed to get encryption key: %v", err)
}

// Use in configuration
config := panel.Config{
	Database: panel.DatabaseConfig{
		Host:     dbCreds["host"],
		Port:     dbCreds["port"],
		Username: dbCreds["username"],
		Password: dbCreds["password"],
		Database: dbCreds["database"],
	},
	Encryption: panel.EncryptionConfig{
		KeyHex: encryptionKey,
	},
}
```

---

### Option 2: AWS Secrets Manager (AWS Deployments)

**Best for**: AWS-hosted applications, serverless deployments

#### Setup

```bash
# Install AWS CLI
brew install awscli  # macOS

# Configure AWS credentials
aws configure
```

#### Store Secrets

```bash
# Store database credentials
aws secretsmanager create-secret \
  --name panel/database \
  --secret-string '{"host":"localhost","port":"5432","username":"panel_user","password":"secure_password","database":"panel_db"}'

# Store encryption key
aws secretsmanager create-secret \
  --name panel/encryption \
  --secret-string '{"key":"<64-char-hex-key>"}'
```

#### Go Integration

```go
package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSSecretsManager struct {
	client *secretsmanager.Client
}

func NewAWSSecretsManager(ctx context.Context) (*AWSSecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	return &AWSSecretsManager{client: client}, nil
}

func (sm *AWSSecretsManager) GetSecret(ctx context.Context, secretName string) (map[string]string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := sm.client.GetSecretValue(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}

	var secretData map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return secretData, nil
}

func (sm *AWSSecretsManager) GetDatabaseCredentials(ctx context.Context) (map[string]string, error) {
	return sm.GetSecret(ctx, "panel/database")
}

func (sm *AWSSecretsManager) GetEncryptionKey(ctx context.Context) (string, error) {
	secret, err := sm.GetSecret(ctx, "panel/encryption")
	if err != nil {
		return "", err
	}
	return secret["key"], nil
}
```

---

### Option 3: Environment Variables with .env Files (Development)

**Best for**: Local development, testing environments

#### Setup

```bash
# Install godotenv
go get github.com/joho/godotenv
```

#### Create .env File

```bash
# .env (NEVER commit this file to git!)
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=panel_user
DATABASE_PASSWORD=secure_password
DATABASE_NAME=panel_db

ENCRYPTION_KEY=<64-char-hex-key>
SESSION_SECRET=<random-secret>

# Add to .gitignore
echo ".env" >> .gitignore
```

#### Go Integration

```go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file in development
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}
}

func main() {
	config := panel.Config{
		Database: panel.DatabaseConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Username: os.Getenv("DATABASE_USERNAME"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Database: os.Getenv("DATABASE_NAME"),
		},
		Encryption: panel.EncryptionConfig{
			KeyHex: os.Getenv("ENCRYPTION_KEY"),
		},
	}

	// ... rest of application
}
```

---

## Secret Rotation

### Encryption Key Rotation

```go
package config

import (
	"time"
)

type KeyRotationConfig struct {
	Enabled          bool
	RotationInterval time.Duration
	OldKeys          []string // Keep old keys for decryption
}

func (c *Config) RotateEncryptionKey() error {
	// 1. Generate new key
	newKey := generateNewKey()

	// 2. Store old key for decryption
	c.Encryption.OldKeys = append(c.Encryption.OldKeys, c.Encryption.KeyHex)

	// 3. Update to new key
	c.Encryption.KeyHex = newKey

	// 4. Re-encrypt sensitive data with new key
	if err := c.reencryptData(); err != nil {
		return err
	}

	// 5. Remove old keys after grace period (e.g., 90 days)
	c.cleanupOldKeys()

	return nil
}

func (c *Config) reencryptData() error {
	// Re-encrypt all sensitive data in database
	// This should be done in batches to avoid performance issues
	return nil
}

func (c *Config) cleanupOldKeys() {
	// Remove keys older than grace period
	gracePerio := 90 * 24 * time.Hour
	// Implementation depends on how you track key age
}
```

---

## Best Practices

### 1. Never Hardcode Secrets

❌ **Bad**:
```go
const DatabasePassword = "my_password"
const EncryptionKey = "0123456789abcdef..."
```

✅ **Good**:
```go
password := os.Getenv("DATABASE_PASSWORD")
encryptionKey := secretsManager.GetEncryptionKey()
```

### 2. Use Strong Encryption Keys

```bash
# Generate strong encryption key (256-bit)
openssl rand -hex 32

# Generate session secret
openssl rand -base64 32
```

### 3. Implement Least Privilege Access

- Only grant access to secrets that are needed
- Use separate credentials for different environments
- Rotate credentials regularly

### 4. Audit Secret Access

```go
func (sm *SecretsManager) GetSecret(name string) (string, error) {
	// Log secret access for audit
	log.Printf("Secret accessed: %s by user: %s", name, getCurrentUser())

	// Get secret
	secret, err := sm.client.GetSecret(name)
	if err != nil {
		log.Printf("Failed to access secret: %s, error: %v", name, err)
		return "", err
	}

	return secret, nil
}
```

### 5. Never Log Secrets

❌ **Bad**:
```go
log.Printf("Database password: %s", password)
log.Printf("Encryption key: %s", encryptionKey)
```

✅ **Good**:
```go
log.Printf("Database connection established")
log.Printf("Encryption initialized")
```

### 6. Use Secret Scanning Tools

```bash
# Install gitleaks
brew install gitleaks

# Scan repository for secrets
gitleaks detect --source . --verbose

# Add pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
gitleaks protect --staged --verbose
EOF
chmod +x .git/hooks/pre-commit
```

---

## Implementation Checklist

- [ ] Choose secrets management solution (Vault, AWS Secrets Manager, or .env)
- [ ] Install and configure secrets management tool
- [ ] Migrate existing secrets from environment variables
- [ ] Update application code to use secrets manager
- [ ] Implement secret rotation policies
- [ ] Set up audit logging for secret access
- [ ] Add secret scanning to CI/CD pipeline
- [ ] Document secret management procedures
- [ ] Train team on secrets management best practices
- [ ] Test secret rotation procedures
- [ ] Set up monitoring and alerting for secret access

---

## Security Considerations

### 1. Secret Storage

- **Never** commit secrets to version control
- **Never** store secrets in plain text files
- **Always** use encrypted storage for secrets
- **Always** use access controls for secret storage

### 2. Secret Transmission

- **Always** use TLS/HTTPS for secret transmission
- **Never** send secrets via email or chat
- **Never** include secrets in URLs or query parameters

### 3. Secret Lifecycle

- **Rotate** secrets regularly (every 90 days minimum)
- **Revoke** secrets immediately when compromised
- **Monitor** secret access and usage
- **Audit** secret access logs regularly

### 4. Development vs Production

- **Use** different secrets for each environment
- **Never** use production secrets in development
- **Test** secret rotation in staging before production
- **Document** secret management procedures

---

## Troubleshooting

### Secret Not Found

```go
secret, err := secretsManager.GetSecret("panel/database")
if err != nil {
	if errors.Is(err, ErrSecretNotFound) {
		log.Fatal("Secret not found. Please create it first.")
	}
	log.Fatalf("Failed to get secret: %v", err)
}
```

### Secret Access Denied

```bash
# Check Vault token
vault token lookup

# Check AWS IAM permissions
aws sts get-caller-identity
aws secretsmanager list-secrets
```

### Secret Rotation Failed

```go
func (c *Config) RotateEncryptionKey() error {
	// Create backup before rotation
	if err := c.backupCurrentKey(); err != nil {
		return fmt.Errorf("failed to backup key: %w", err)
	}

	// Attempt rotation
	if err := c.performRotation(); err != nil {
		// Rollback on failure
		if rollbackErr := c.rollbackKey(); rollbackErr != nil {
			return fmt.Errorf("rotation failed and rollback failed: %w, %w", err, rollbackErr)
		}
		return fmt.Errorf("rotation failed (rolled back): %w", err)
	}

	return nil
}
```

---

## References

- [HashiCorp Vault Documentation](https://www.vaultproject.io/docs)
- [AWS Secrets Manager Documentation](https://docs.aws.amazon.com/secretsmanager/)
- [OWASP Secrets Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [12-Factor App: Config](https://12factor.net/config)

---

## Next Steps

1. **Immediate**: Implement .env file support for local development
2. **Short-term**: Set up HashiCorp Vault or AWS Secrets Manager for production
3. **Medium-term**: Implement secret rotation policies
4. **Long-term**: Integrate with enterprise secrets management solution
