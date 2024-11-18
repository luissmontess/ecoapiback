package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID 	            primitive.ObjectID                `bson:"UserID"`
	FBID                string                            `bson:"FBID"`
	Username            string                            `bson:"Username"`
	Email               string                            `bson:"Email"`
	Cardboard           float32                           `bson:"Cardboard"`
	Glass               float32                           `bson:"Glass"`
	Tetrapack           float32                           `bson:"Tetrapack"`
	Plastic             float32                           `bson:"Plastic"`
	Paper               float32                           `bson:"Paper"`
	Metal               float32                           `bson:"Metal"`
	MedalTrans          int                               `bson:"MedalTrans"`
	MedalEnergy         int                               `bson:"MedalEnergy"`
	MedalConsume        int                               `bson:"MedalConsume"`
	MedalDesecho        int                               `bson:"MedalDesecho"`
	NotificationArray   []Notification                    `bson:"NotificationArray"`
}  
  
type Collaborator struct {  
	CollaboratorID      primitive.ObjectID                `bson:"CollaboratorID"`
	FBID                string                            `bson:"FBID"`
	Username            string                            `bson:"Username"`
	Email               string                            `bson:"Email"`
}

type Course struct {
	CourseID  	        primitive.ObjectID 		          `bson:"CourseID"`
	CollaboratorFBID    string                            `bson:"CollaboratorFBID"`
	Title               string                            `bson:"Title"`
	Pillar    	        string             		          `bson:"Pillar"`
	StartTime           primitive.DateTime                `bson:"StartTime"`
	EndTime             primitive.DateTime                `bson:"EndTime"`
	Longitude           float64                           `bson:"Longitude"`
	Latitude            float64                           `bson:"Latitude"`
	Limit               float32                           `bson:"Limit"`
	AssistantArray      []Assistant                       `bson:"AssistantArray"`
}    
  
type Recollect struct {  
	RecollectID         primitive.ObjectID                `bson:"RecollectID"`
	CollaboratorFBID    string                            `bson:"CollaboratorFBID"`
	Cardboard           float32                           `bson:"Cardboard"`
	Glass               float32                           `bson:"Glass"`
	Tetrapack           float32                           `bson:"Tetrapack"`
	Plastic             float32                           `bson:"Plastic"`
	Paper               float32                           `bson:"Paper"`
	Metal               float32                           `bson:"Metal"`
	StartTime           primitive.DateTime                `bson:"StartTime"`
	EndTime             primitive.DateTime                `bson:"EndTime"`
	Longitude           float64                           `bson:"Longitude"`
	Latitude            float64                           `bson:"Latitude"`
	Limit               float32                           `bson:"Limit"`
	DonationArray       []Donation                        `bson:"DonationArray"`
}  
  
type Medal struct {  
MedalID             	primitive.ObjectID                `bson:"MedalID"`
Pillar              	string                            `bson:"Pillar"`
}

type Notification struct {
	NotificationType    string                            `bson:"MedalType"`
	Message             string                            `bson:"Message"`
}

type Donation struct {
	UserFBID            string                            `bson:"UserFBID"`
	Username            string                            `bson:"Username"`
	Cardboard           float32                           `bson:"Cardboard"`
	Glass               float32                           `bson:"Glass"`
	Tetrapack           float32                           `bson:"Tetrapack"`
	Plastic             float32                           `bson:"Plastic"`
	Paper               float32                           `bson:"Paper"`
	Metal               float32                           `bson:"Metal"`
}

type Assistant struct {
	UserFBID            string                            `bson:"UserFBID"`
	Username            string                            `bson:"Username"`
	Email               string                            `bson:"Email"`
}

type TopTenUser struct {
	UserFBID            string                            `bson:"UserFBID"`
	Username            string                            `bson:"Username"`
	Email               string                            `bson:"Email"`
	Place               int                               `bson:"Place"`
}

type TopTenArray struct {
	UserArray           []TopTenUser                      `bson:"UserArray"`  
}