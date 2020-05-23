package seeds

import (
	"gohome_self/tracingbook/infrastructure"
	"gohome_self/tracingbook/models"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func randomInt(min, max int) int {

	return rand.Intn(max-min) + min
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func seedAdmin(db *gorm.DB) {
	count := 0
	adminRole := models.Role{Name: "ROLE_ADMIN", Description: "Only for admin"}
	query := db.Model(&models.Role{}).Where("name = ?", "ROLE_ADMIN")
	query.Count(&count)

	if count == 0 {
		db.Create(&adminRole)
	} else {
		query.First(&adminRole)
	}

	adminRoleUsers := 0
	var adminUsers []models.User
	db.Model(&adminRole).Related(&adminUsers, "Users")

	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminRoleUsers)
	if adminRoleUsers == 0 {

		// query.First(&adminRole) // First would fetch the Role admin because the query status name='ROLE_ADMIN'
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		// Approach 1
		user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: string(password)}
		user.Roles = append(user.Roles, adminRole)

		// Do not try to update the adminRole
		db.Set("gorm:association_autoupdate", false).Create(&user)

		// Approach 2
		// user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: "password"}
		// user.Roles = append(user.Roles, adminRole)
		// db.NewRecord(user)
		// db.Set("gorm:association_autoupdate", false).Save(&user)

		if db.Error != nil {
			print(db.Error)
		}
	}
}

func seedUsers(db *gorm.DB) {
	count := 0
	role := models.Role{Name: "ROLE_USER", Description: "Only for standard users"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_USER")
	q.Count(&count)

	if count == 0 {
		db.Create(&role)
	} else {
		q.First(&role)
	}

	var standardUsers []models.User
	db.Model(&role).Related(&standardUsers, "Users")
	usersCount := len(standardUsers)
	usersToSeed := 20
	usersToSeed -= usersCount
	if usersToSeed > 0 {
		for i := 0; i < usersToSeed; i++ {
			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			user := models.User{FirstName: fake.FirstName(), LastName: fake.LastName(), Email: fake.EmailAddress(), Username: fake.UserName(),
				Password: string(password)}
			// No need to add the role as we did for seedAdmin, it is added by the BeforeSave hook
			db.Set("gorm:association_autoupdate", false).Create(&user)
		}
	}
}

func Seed() {
	db := infrastructure.GetDb()
	rand.Seed(time.Now().UnixNano())
	seedAdmin(db)
	seedUsers(db)
}
