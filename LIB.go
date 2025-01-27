package main

import (
	
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)
var Db *gorm.DB
type BookIssueDetail struct{
	ID uint8 `gorm:"primaryKey"`
	IssuerName string `gorm:"size:100;not null`
	IssuerDay uint8 `gorm:"not null"`
	IssuerMonth uint8 `gorm:"not null"`
	IssuerYear uint16 `gorm:"not null"`

}

func connectToDB (c*gin.Context){
	var err error
	
	dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	Db,err =gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err!=nil{
		c.IndentedJSON(http.StatusServiceUnavailable,gin.H{"message":"Database connection Failed!"})
	}
	c.IndentedJSON(http.StatusOK,gin.H{"message":"Connection Successful"})
}
func migrateDB(c*gin.Context){
	var err error
	if Db==nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{
			"message":"Connection not established!",
		})

	}
	err=Db.AutoMigrate(&BookIssueDetail{})
	if err!=nil{
		c.IndentedJSON(400,gin.H{
			"message":"table not created",
		})
		return
	}
	c.IndentedJSON(http.StatusAccepted,gin.H{
		"message":"table  create",
	})
}
func InsertDetails(c*gin.Context){
	var newIssue BookIssueDetail
	err :=c.BindJSON(&newIssue)
	if err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{
			"message":"Error Data Do not match with the schema",
		})
		return
	}
	if err:=Db.Create(&newIssue).Error;err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{
			"message":"Unable to insert Data!!",
			
		})
		return
	}
	c.IndentedJSON(http.StatusAccepted,gin.H{
		"message":"Data added succefully",
	})


}
func update(c*gin.Context){
	id:=c.Param("id")
	newname :=c.Param("newname")
	if err:=Db.Model(&BookIssueDetail{}).Where("id = ?",id).Update("IssuerName",newname).Error;err!=nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{
			"message":"Failed to Update details",
		})
		return
	}
	c.IndentedJSON(http.StatusAccepted,gin.H{
		"message":"Name Updated Successfully",
	})

}
func ReadData(c*gin.Context){
	var data BookIssueDetail
	id:=c.Param("id")
	if err:=Db.First(&data,id).Error;err!=nil{
		c.IndentedJSON(http.StatusAccepted,gin.H{
			"message":"User not found",
		})
		return

	}
	
	c.IndentedJSON(http.StatusAccepted,data)

}
func delete(c*gin.Context){
	id:=c.Param("id")

	if err:=Db.Delete(&BookIssueDetail{},id).Error;err!=nil{
		c.IndentedJSON(http.StatusNotFound,gin.H{
			"message":"Unable to delete",
		})
		return
	}
	c.IndentedJSON(http.StatusAccepted,gin.H{
		"message":"Deleted Successfully",
	})
}
func main(){
	router:=gin.Default()

	router.GET("/ConnectTODB",connectToDB)
	router.GET("/BOOKISSUEDETAILS",migrateDB)
	router.POST("/BOOKISSUEDETAILS/Create",InsertDetails)
	router.GET("/BOOKISSUEDETAILS/Update/:id/:newname",update)
	router.GET("/BOOKISSUEDETAILS/Read/:id",ReadData)
	router.GET("/BOOKISSUEDETAILS/Delete/:id",delete)



	router.Run("localhost:8080")
}

//Step1 :=GET: http://localhost:8080/ConnectToDB (To connect to DB)
//Step2 :=GET: http://localhost:8080/BOOKISSUEDETAILS (To Migrate to DB)
//Step3 :=POST: http://localhost:8080/BOOKISSUEDETAILS/Create/ 
//				{DATA...} (TO Add data in table)
//Step4 :=GET: http://localhost:8080/BOOKISSUEDETAILS/Update/2/Himesh (TO Update Detail using ID)
//Step5 :=GET: http://localhost:8080/BOOKISSUEDETAILS/Read/2 (To read data using ID)
//Step6 :=GET: http://localhost:8080/BOOKISSUEDETAILS/Delete/1 (To Delete data using Id) 