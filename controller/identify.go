package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pavanmuttange/Bitespeed/model"
	"github.com/pavanmuttange/Bitespeed/pkg/config"
)

func Identify(c *gin.Context) {
	db := config.DB
	var request model.IdentifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if request.Email == nil && request.PhoneNumber == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phoneNumber required"})
		return
	}

	var contacts []model.Contact
	if err := db.Where("email = ? OR phone_number = ?", request.Email, request.PhoneNumber).Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	var primaryContact *model.Contact
	secondaryContactIDs := []uint{}
	emails := make(map[string]struct{})
	phoneNumbers := make(map[string]struct{})

	for _, contact := range contacts {
		if contact.LinkPrecedence == "primary" {
			primaryContact = &contact
		}
		if contact.Email != nil {
			emails[*contact.Email] = struct{}{}
		}
		if contact.PhoneNumber != nil {
			phoneNumbers[*contact.PhoneNumber] = struct{}{}
		}
		if contact.LinkedID != nil {
			secondaryContactIDs = append(secondaryContactIDs, contact.ID)
		}
	}

	if primaryContact == nil && len(contacts) > 0 {
		primaryContact = &contacts[0]
		primaryContact.LinkPrecedence = "primary"
		db.Save(primaryContact)
		for i := 1; i < len(contacts); i++ {
			contacts[i].LinkedID = &primaryContact.ID
			contacts[i].LinkPrecedence = "secondary"
			db.Save(&contacts[i])
		}
	}

	if primaryContact == nil {
		// Create new primary contact
		newContact := model.Contact{
			Email:          request.Email,
			PhoneNumber:    request.PhoneNumber,
			LinkPrecedence: "primary",
		}
		db.Create(&newContact)
		primaryContact = &newContact
	} else {
		// Create secondary contact if new email or phone number is encountered
		if request.Email != nil && (primaryContact.Email == nil || *primaryContact.Email != *request.Email) {
			newContact := model.Contact{
				Email:          request.Email,
				PhoneNumber:    primaryContact.PhoneNumber,
				LinkedID:       &primaryContact.ID,
				LinkPrecedence: "secondary",
			}
			db.Create(&newContact)
			secondaryContactIDs = append(secondaryContactIDs, newContact.ID)
		}

		if request.PhoneNumber != nil && (primaryContact.PhoneNumber == nil || *primaryContact.PhoneNumber != *request.PhoneNumber) {
			newContact := model.Contact{
				Email:          primaryContact.Email,
				PhoneNumber:    request.PhoneNumber,
				LinkedID:       &primaryContact.ID,
				LinkPrecedence: "secondary",
			}
			db.Create(&newContact)
			secondaryContactIDs = append(secondaryContactIDs, newContact.ID)
		}
	}

	if request.Email != nil {
		emails[*request.Email] = struct{}{}
	}
	if request.PhoneNumber != nil {
		phoneNumbers[*request.PhoneNumber] = struct{}{}
	}

	emailList := []string{}
	for email := range emails {
		emailList = append(emailList, email)
	}
	phoneNumberList := []string{}
	for phoneNumber := range phoneNumbers {
		phoneNumberList = append(phoneNumberList, phoneNumber)
	}

	response := model.ContactResponse{
		PrimaryContactID:    primaryContact.ID,
		Emails:              emailList,
		PhoneNumbers:        phoneNumberList,
		SecondaryContactIDs: secondaryContactIDs,
	}

	c.JSON(http.StatusOK, gin.H{"contact": response})
}
