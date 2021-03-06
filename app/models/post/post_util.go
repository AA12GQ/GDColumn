package post

import (
     "GDColumn/pkg/app"
     "GDColumn/pkg/database"
     "GDColumn/pkg/paginator"
     "gorm.io/gorm/clause"

    "github.com/gin-gonic/gin"
)

func Get(idstr string) (post Post) {
    database.DB.Preload(clause.Associations).Where("id", idstr).First(&post)
    return
}

func GetImage(idstr string) (imageModel Image) {
    database.DB.Where("id", idstr).First(&imageModel)
    return
}

func GetAll(idstr string)(post []Post){
    database.DB.Where("column_id = ?", idstr).Find(&post)
    return
}

func GetBy(field, value string) (post Post) {
    database.DB.Where("? = ?", field, value).First(&post)
    return
}

func All() (posts []Post) {
    database.DB.Find(&posts)
    return 
}

func IsExist(field, value string) bool {
    var count int64
    database.DB.Model(Post{}).Where(" = ?", field, value).Count(&count)
    return count > 0
}

func Paginate(c *gin.Context, perPage int,id string) (posts []Post, paging paginator.Paging) {
    paging = paginator.Paginate(
        c,
        database.DB.Where("column_id = ?",id).Model(Post{}),
        &posts,
        app.V1URL(database.TableName(&Post{})),
        perPage,
    )
    return
}