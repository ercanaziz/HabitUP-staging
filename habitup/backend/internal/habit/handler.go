package habit

import (
	"context"
	"net/http"
	"sort"
	"time"

	"habitup/internal/db"
	"habitup/internal/queue"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func userID(c *gin.Context) (primitive.ObjectID, bool) {
	raw, _ := c.Get("userID")
	id, err := primitive.ObjectIDFromHex(raw.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Geçersiz kullanıcı"})
		return primitive.NilObjectID, false
	}
	return id, true
}

// Gereksinim 3: Yeni Alışkanlık Tanımlama
func Create(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	var req CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	habit := Habit{
		ID:          primitive.NewObjectID(),
		UserID:      uid,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if _, err := db.Col("habits").InsertOne(context.Background(), habit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Alışkanlık oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, habit)
}

// Gereksinim 4: Alışkanlıkları Listeleme
func List(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	cursor, err := db.Col("habits").Find(context.Background(), bson.M{"userId": uid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Alışkanlıklar alınamadı"})
		return
	}
	defer cursor.Close(context.Background())

	var habits []Habit
	if err := cursor.All(context.Background(), &habits); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Veri okuma hatası"})
		return
	}

	if habits == nil {
		habits = []Habit{}
	}

	c.JSON(http.StatusOK, habits)
}

// Gereksinim 6: Alışkanlık Güncelleme
func Update(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	habitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Geçersiz ID"})
		return
	}

	var req UpdateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}

	res, err := db.Col("habits").UpdateOne(
		context.Background(),
		bson.M{"_id": habitID, "userId": uid},
		bson.M{"$set": update},
	)
	if err != nil || res.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Alışkanlık bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alışkanlık güncellendi"})
}

// Gereksinim 7: Alışkanlık Silme
func Delete(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	habitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Geçersiz ID"})
		return
	}

	ctx := context.Background()
	res, err := db.Col("habits").DeleteOne(ctx, bson.M{"_id": habitID, "userId": uid})
	if err != nil || res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Alışkanlık bulunamadı"})
		return
	}

	// Bağlı check kayıtlarını da sil
	_, _ = db.Col("checks").DeleteMany(ctx, bson.M{"habitId": habitID, "userId": uid})

	c.JSON(http.StatusNoContent, nil)
}

// Gereksinim 5: Alışkanlık Durumu Güncelleme (tamamlandı işaretle)
func Check(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	habitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Geçersiz ID"})
		return
	}

	today := time.Now().Format("2006-01-02")
	ctx := context.Background()

	// Bugün zaten işaretli mi?
	var existing CheckRecord
	findErr := db.Col("checks").FindOne(ctx, bson.M{
		"habitId": habitID,
		"userId":  uid,
		"date":    today,
	}).Decode(&existing)

	if findErr == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Bugün zaten işaretlendi"})
		return
	}
	if findErr != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Sunucu hatası"})
		return
	}

	record := CheckRecord{
		ID:        primitive.NewObjectID(),
		HabitID:   habitID,
		UserID:    uid,
		Date:      today,
		CheckedAt: time.Now(),
	}

	if _, err := db.Col("checks").InsertOne(ctx, record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "İşaretleme kaydedilemedi"})
		return
	}

	// RabbitMQ event publish
	_ = queue.PublishHabitChecked(ctx, uid.Hex(), habitID.Hex())

	c.JSON(http.StatusOK, gin.H{"message": "Alışkanlık tamamlandı olarak işaretlendi"})
}

// Gereksinim 8: İşaretlemeyi Geri Alma
func Uncheck(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	habitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Geçersiz ID"})
		return
	}

	today := time.Now().Format("2006-01-02")
	res, err := db.Col("checks").DeleteOne(context.Background(), bson.M{
		"habitId": habitID,
		"userId":  uid,
		"date":    today,
	})

	if err != nil || res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bugün için işaretleme bulunamadı"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Gereksinim 9: İstatistik ve Seri Takibi
func Stats(c *gin.Context) {
	uid, ok := userID(c)
	if !ok {
		return
	}

	habitID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Geçersiz ID"})
		return
	}

	ctx := context.Background()

	// Alışkanlığın oluşturulma tarihini bul
	var h Habit
	if err := db.Col("habits").FindOne(ctx, bson.M{"_id": habitID, "userId": uid}).Decode(&h); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Alışkanlık bulunamadı"})
		return
	}

	// Tüm check kayıtlarını getir
	cursor, err := db.Col("checks").Find(ctx, bson.M{"habitId": habitID, "userId": uid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "İstatistik alınamadı"})
		return
	}
	defer cursor.Close(ctx)

	var records []CheckRecord
	_ = cursor.All(ctx, &records)

	// Tarihleri sırala
	dates := make([]string, 0, len(records))
	for _, r := range records {
		dates = append(dates, r.Date)
	}
	sort.Strings(dates)

	totalChecks := len(dates)
	daysSinceCreation := int(time.Since(h.CreatedAt).Hours()/24) + 1
	var completionRate float64
	if daysSinceCreation > 0 {
		completionRate = float64(totalChecks) / float64(daysSinceCreation) * 100
	}

	currentStreak, longestStreak := calculateStreaks(dates)

	c.JSON(http.StatusOK, HabitStats{
		CurrentStreak:  currentStreak,
		LongestStreak:  longestStreak,
		CompletionRate: completionRate,
		TotalChecks:    totalChecks,
	})
}

func calculateStreaks(sortedDates []string) (current, longest int) {
	if len(sortedDates) == 0 {
		return 0, 0
	}

	today := time.Now().Format("2006-01-02")
	streak := 1
	maxStreak := 1

	for i := len(sortedDates) - 1; i > 0; i-- {
		t1, _ := time.Parse("2006-01-02", sortedDates[i])
		t2, _ := time.Parse("2006-01-02", sortedDates[i-1])
		if t1.Sub(t2).Hours() == 24 {
			streak++
			if streak > maxStreak {
				maxStreak = streak
			}
		} else {
			break
		}
	}

	// Bugün veya dün işaretlenmişse current streak aktif
	last := sortedDates[len(sortedDates)-1]
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if last != today && last != yesterday {
		streak = 0
	}

	return streak, maxStreak
}
