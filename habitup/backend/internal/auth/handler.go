package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"habitup/internal/cache"
	"habitup/internal/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	col := db.Col("users")
	ctx := context.Background()

	var existing User
	if err := col.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bu e-posta zaten kayıtlı"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Şifre hashlenemedi"})
		return
	}

	user := User{
		ID:       primitive.NewObjectID(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}

	if _, err := col.InsertOne(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Kullanıcı oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kullanıcı başarıyla oluşturuldu"})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	col := db.Col("users")
	ctx := context.Background()

	var user User
	if err := col.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Hatalı e-posta veya şifre"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Sunucu hatası"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Hatalı e-posta veya şifre"})
		return
	}

	token, err := GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token oluşturulamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	header := c.GetHeader("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")

	claims, err := ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Çıkış yapıldı"})
		return
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	_ = cache.BlacklistToken(context.Background(), token, ttl)

	c.JSON(http.StatusOK, gin.H{"message": "Başarıyla çıkış yapıldı"})
}
