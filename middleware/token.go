package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func getSecretKey() []byte {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("Chave secreta não encontrada no arquivo .env")
		os.Exit(1)
	}
	return []byte(secretKey)
}

func generateToken() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}
func generateTokenTest() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 9999).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}

// @Summary Criar token de auth
// @Description Cria um toke para auth do usuario
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} string "token make:"
// @Router /token [get]
func generateTokenHandler(c *gin.Context) {
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}
	token = "Bearer " + token
	c.Set("Response", token)
	c.Status(http.StatusOK)
}

type resTest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// @Summary Test Token
// @Description Cria um toke para auth do usuario
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Token de autenticação" default(Bearer <token>)
// @Success 200 {object} string "token make:"
// @Router /api/test-token [post]
func testeToken(c *gin.Context) {
	res := resTest{
		Name:        "token",
		Description: "This token is valid",
	}
	statusCode := http.StatusOK

	c.Set("Response", res)
	c.Status(statusCode)
}
