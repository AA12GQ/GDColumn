package migrations

import (
    "database/sql"
    "GDColumn/app/models"
    "GDColumn/pkg/migrate"

    "gorm.io/gorm"
)

func init() {

    type User struct {
        models.BaseModel

        Name     string `gorm:"type:varchar(255);not null;index"`
        Email    string `gorm:"type:varchar(255);index;default:null"`
        Phone    string `gorm:"type:varchar(20);index;default:null"`
        Password string `gorm:"type:varchar(255)"`

        models.CommonTimestampsField
    }

    up := func(migrator gorm.Migrator, DB *sql.DB) {
        migrator.AutoMigrate(&User{})
    }

    down := func(migrator gorm.Migrator, DB *sql.DB) {
        migrator.DropTable(&User{})
    }

    migrate.Add("2022_05_16_003445_add_users_table", up, down)
}