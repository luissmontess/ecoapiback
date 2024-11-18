package handlers

import (
	"context"
	"ecoApi/config"
	"ecoApi/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

func CreateCourse(c *fiber.Ctx) error {
	collection := config.MongoDatabase.Collection("courses")
	var course models.Course

	if err := c.BodyParser(&course); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse course JSON",
		})
	}

	if course.CollaboratorFBID == "" ||
	   	course.Pillar == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields for course",
		})
	}

	course.CourseID = primitive.NewObjectID()
	course.AssistantArray = []models.Assistant{}

	_, err := collection.InsertOne(context.TODO(), course)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to insert course",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(course)
}

func GetCourseByID(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	// usar funcion hex para convertir primitive
	courseObjectID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID",
		})
	}

	coursesCollection := config.MongoDatabase.Collection("courses")

	// buscar por id
	var course models.Course
	err = coursesCollection.FindOne(context.TODO(), bson.M{"CourseID": courseObjectID}).Decode(&course)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	return c.JSON(course)
}


func GetActiveCourses (c *fiber.Ctx) error {
	currentTime := primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{
		"StartTime": bson.M{"$lte": currentTime}, // less than or equal to
		"EndTime":   bson.M{"$gte": currentTime}, // greater than or equal to
	}

	courseCollection := config.MongoDatabase.Collection("courses")
	var courses []models.Course

	cursor, err := courseCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving active courses",
		})
	}

	if err := cursor.All(context.TODO(), &courses); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding active courses",
		})
	}

	return c.JSON(courses)
}

func GetCollaboratorCourses(c *fiber.Ctx) error {
	uid := c.Params("uid") 
	filter := bson.M{"CollaboratorFBID": uid}
	var courses []models.Course

	courseCollection := config.MongoDatabase.Collection("courses")
	courseCursor, err := courseCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error retrieving courses",
        })
	}
	defer courseCursor.Close(context.TODO())

	if err := courseCursor.All(context.TODO(), &courses); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding active courses",
		})
	}

	return c.JSON(courses)
}

func AddUserToCourse(c *fiber.Ctx) error {
	type Request struct {
		UserFBID string `json:"UserFBID"`
	}

	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Course ID is required",
		})
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	courseObjectID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID",
		})
	}

	coursesCollection := config.MongoDatabase.Collection("courses")
	usersCollection := config.MongoDatabase.Collection("users")

	var course models.Course
	err = coursesCollection.FindOne(context.TODO(), bson.M{"CourseID": courseObjectID}).Decode(&course)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found",
		})
	}


	if float32(len(course.AssistantArray)) >= course.Limit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Assistant Limit has been reached",
		})
	}

	// usar userFBID esta vez
	var user models.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"FBID": req.UserFBID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	filter := bson.M{"CourseID": courseObjectID}
	update := bson.M{"$addToSet": bson.M{
		"AssistantArray": bson.M{
			"UserFBID": user.FBID,
			"Username": user.Username,
			"Email": user.Email,
		},
	}}
	_, err = coursesCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add assistant to course",
		})
	}

	var medalUpdate bson.M
	var notification models.Notification

	switch course.Pillar {
	case "ENERGIA":
		medalUpdate = bson.M{"$inc": bson.M{"MedalEnergy": 1}}
		notification = models.Notification{
			NotificationType: "ENERGIA",
			Message:   "Obtuviste una medalla por participar en el curso " + course.Title + ", con pilar de energia",
		}
	case "DESECHO":
		medalUpdate = bson.M{"$inc": bson.M{"MedalDesecho": 1}}
		notification = models.Notification{
			NotificationType: "DESECHO",
			Message:   "Obtuviste una medalla por participar en el curso " + course.Title + ", con pilar de desecho",
		}
	case "TRANSPORTE":
		medalUpdate = bson.M{"$inc": bson.M{"MedalTrans": 1}}
		notification = models.Notification{
			NotificationType: "TRANSPORTE",
			Message:   "Obtuviste una medalla por participar en el curso " + course.Title + ", con pilar de Transporte",
		}
	case "CONSUMO":
		medalUpdate = bson.M{"$inc": bson.M{"MedalConsume": 1}}
		notification = models.Notification{
			NotificationType: "CONSUMO",
			Message:   "Obtuviste una medalla por participar en el curso " + course.Title + ", con pilar de consumo",
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pillar type",
		})
	}

	userFilter := bson.M{"FBID": req.UserFBID}

	_, err = usersCollection.UpdateOne(context.TODO(), userFilter, medalUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user's medal count",
		})
	}
	
	notificationUpdate := bson.M{"$push": bson.M{"NotificationArray": notification}}

	_, err = usersCollection.UpdateOne(context.TODO(), userFilter, notificationUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add notification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Assistant and course ID added successfully and medal updated and user notification added",
	})
}