package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
	Memory      uint32 // Memory cost in KB
	Iterations  uint32 // Time cost (number of iterations)
	Parallelism uint8  // Number of threads
	SaltLength  uint32 // Length of salt in bytes
	KeyLength   uint32 // Length of generated keys in bytes
}

func DefaultConfig() *PasswordConfig {
	return &PasswordConfig{
		Memory:      64 * 1024, // 64MB
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16, // 16 bytes salt
		KeyLength:   32, // 32 bytes key
	}
}

// securely hashes a password using Argon2id
// returns the encoded hash string that can be stores in database
func HashPassword(password string) (string, error) {
	return HashPasswordWithConfig(password, DefaultConfig())
}

func HashPasswordWithConfig(password string, passwordConfig *PasswordConfig) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	salt := make([]byte, passwordConfig.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		passwordConfig.Iterations,
		passwordConfig.Memory,
		passwordConfig.Parallelism,
		passwordConfig.KeyLength,
	)

	// encode in PHC format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	encodeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		passwordConfig.Memory,
		passwordConfig.Iterations,
		passwordConfig.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodeHash, nil
}

// VerifyPassword compares a plainText password with a hashed password
// Returns true if the password matches the hash, false otherwise
func VerifyPassword(password, hashedPassword string) (bool, error) {
	if password == "" {
		return false, errors.New("password cannot be empty")
	}
	if hashedPassword == "" {
		return false, errors.New("hashed password cannot be empty")
	}

	// parse the stored hash to extract parameters and components
	passwordConfig, salt, hash, err := parseHash(hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to parse the hash: %w", err)
	}

	// hash the input password using the same parameters
	inputHash := argon2.IDKey(
		[]byte(password),
		salt,
		passwordConfig.Iterations,
		passwordConfig.Memory,
		passwordConfig.Parallelism,
		passwordConfig.KeyLength,
	)

	// use constant time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(hash, inputHash) == 1, nil
}

// parseHash parses an Argon2id hash string and extracts the parameters and components
// Expected format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
func parseHash(encodedHash string) (*PasswordConfig, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	// Check algorithm
	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("unsupported algorithm")
	}

	// Parse version
	if !strings.HasPrefix(parts[2], "v=") {
		return nil, nil, nil, errors.New("invalid version format")
	}
	version, err := strconv.Atoi(parts[2][2:])
	if err != nil || version != argon2.Version {
		return nil, nil, nil, errors.New("unsupported version")
	}

	// Parse parameters (m=memory,t=time,p=parallelism)
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return nil, nil, nil, errors.New("invalid parameters format")
	}

	config := &PasswordConfig{}

	// Parse memory
	if !strings.HasPrefix(params[0], "m=") {
		return nil, nil, nil, errors.New("invalid memory parameter")
	}
	if memory, err := strconv.ParseUint(params[0][2:], 10, 32); err != nil {
		return nil, nil, nil, errors.New("invalid memory value")
	} else {
		config.Memory = uint32(memory)
	}

	// Parse iterations
	if !strings.HasPrefix(params[1], "t=") {
		return nil, nil, nil, errors.New("invalid time parameter")
	}
	if iterations, err := strconv.ParseUint(params[1][2:], 10, 32); err != nil {
		return nil, nil, nil, errors.New("invalid time value")
	} else {
		config.Iterations = uint32(iterations)
	}

	// Parse parallelism
	if !strings.HasPrefix(params[2], "p=") {
		return nil, nil, nil, errors.New("invalid parallelism parameter")
	}
	if parallelism, err := strconv.ParseUint(params[2][2:], 10, 8); err != nil {
		return nil, nil, nil, errors.New("invalid parallelism value")
	} else {
		config.Parallelism = uint8(parallelism)
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	config.SaltLength = uint32(len(salt))

	// Decode hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode hash: %w", err)
	}
	config.KeyLength = uint32(len(hash))

	return config, salt, hash, nil
}
