package utilities

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Adebusy/receiptonboardingsvc/obj"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file" + err.Error())
	}
	return os.Getenv(key)
}

// HashPassword hashes a given password and returns the hashed password or an error
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// CheckPasswordHash verifies the password against the hashed password and returns if it's correct or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DeactivateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 0).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Logout(token, username string) {
	// ttl := time.Now().Add(time.Hour * 0).Unix()

	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	// deleted, delErr := DeleteAuth(au.AccessUuid)

	// if delErr != nil || deleted == 0 {
	// 	c.JSON(http.StatusUnauthorized, "unauthorized")
	// 	return
	// }

	//c.JSON(http.StatusOK, "Successfully logged out")
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// isEmailValid checks if the email provided is valid by regex.
func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// isEmailValid checks if the email provided is valid by regex.
func IsNumberValid(e string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	if re.MatchString(e) {
		return true
	} else {
		return false
	}
}

func SendEmail(toEmail, mailBody string) string {

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	sender := os.Getenv("SMTP_SENDER")
	recipient := toEmail
	from := "From: " + sender + "\n"
	to := "To: " + recipient + "\n"
	subject := "Subject: Digital Company update\n"
	body := mailBody
	message := []byte(from + to + subject + "\n" + body)
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{recipient}, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err.Error())
		return "01"
	} else {
		return "00"
	}
}

const (
	letterBytes               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes              = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes                  = "0123456789"
	companyReceiptsPath       = "./companyReceipts/"
	suggestedReceiptPath      = "./suggestedReceipt/"
	companyIssuedReceiptsPath = "./companyIssuedReceipts/"
)

func TempPassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	b := make([]byte, length)
	for i := range b {
		if useLetters {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[rand.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[rand.Intn(len(numBytes))]
		}
	}
	return string(b)
}

func TrimString(req string) string {
	clean := strings.TrimSpace(req)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimSuffix(clean, "```")
	return clean
}

func saveBase64Image(base64Str string, outputPath string) error {
	imageBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, imageBytes, 0644)
}

type LogoConcept struct {
	ConceptName string `json:"concept_name"`
	LogoID      string `json:"logo_id"`
	LogoBase64  struct {
		Extension string `json:"extension"`
		Data      string `json:"data"`
	} `json:"logo_base64"`
}

func ConvertBase64ToImage(imageName, base64Image string) bool {
	pathAndName := fmt.Sprintf("./logos/%s.svg", imageName)
	err := saveBase64Image(base64Image, pathAndName)
	if err != nil {
		fmt.Print(err)
		return false
	} else {
		return true
	}
}

func ReconstructJSON(lines []string) string {
	var builder strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		builder.WriteString(line)
	}

	return builder.String()
}

func SaveBase64Image(base64Data, outputDir, filename, extension string) error {
	// Ensure directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %w", err)
	}

	// Build full file path
	fullPath := filepath.Join(outputDir, filename+"."+extension)

	// Write file
	if err := os.WriteFile(fullPath, decoded, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func ParseItems(input string) ([]Item, float64) {
	// input: "bread 2 1000, rice 2 100, milk 3 100"
	var items []Item
	var retPrice int
	entries := strings.Split(input, ",")

	for _, entry := range entries {
		parts := strings.Fields(strings.TrimSpace(entry))
		if len(parts) >= 3 {
			name := parts[0]
			qty, _ := strconv.Atoi(parts[1])
			price, _ := strconv.Atoi(parts[2])

			items = append(items, Item{
				Name:     name,
				Quantity: qty,
				Price:    float64(price),
			})
			retPrice += price
		}
	}

	return items, float64(retPrice)
}

func LoadAndUpdateHTMLTemplates(path string, data obj.TemplateData) ([]string, error) {
	var results []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .html files
		if info.IsDir() || filepath.Ext(info.Name()) != ".html" {
			return nil
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		updated := string(content)

		// Replace placeholders
		updated = strings.ReplaceAll(updated, "{{COMPANY_NAME}}", data.CompanyName)
		updated = strings.ReplaceAll(updated, "{{SIGNATURE}}", data.Signature)
		updated = strings.ReplaceAll(updated, "{{LOGO_URL}}", data.LogoURL)

		results = append(results, updated)
		return nil
	})
	return results, err
}

func PickSelectedFile(sourceFileName string) string {
	src := suggestedReceiptPath + sourceFileName

	// Open source file
	in, err := os.Open(src)
	if err != nil {
		panic(err)
	}

	defer in.Close()
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(companyReceiptsPath, 0755); err != nil {
		panic(err)
	}

	// Create destination file
	dst := filepath.Join(companyReceiptsPath, filepath.Base(src))
	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}

	defer out.Close()

	// Copy file content
	if _, err := io.Copy(out, in); err != nil {
		panic(err)
	}
	return "00"
}

func GenerateNewReceipt(sourceFileName, replacements, customerName, totalAmount string) string {
	// input := "bread 2 1000 , rice 2 100, milk 3 100"
	src := companyReceiptsPath + sourceFileName
	content, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}

	html := string(content)
	html = strings.ReplaceAll(html, "TRANSACTION_TOTAL", totalAmount)
	html = strings.ReplaceAll(html, "TRANSACTION_DATE", string(time.Now().String()))
	html = strings.ReplaceAll(html, "CUSTOMER_NAME", customerName)
	html = strings.ReplaceAll(html, "CONTENT_DETAILS", replacements)

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(companyIssuedReceiptsPath, 0755); err != nil {
		panic(err)
	}

	dst := filepath.Join(companyIssuedReceiptsPath, filepath.Base(src))

	if err := os.WriteFile(dst, []byte(html), 0644); err != nil {
		panic(err)
	}
	return "00"
}
