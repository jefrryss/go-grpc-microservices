package model

import "github.com/google/uuid"

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

type User struct {
	ID                  uuid.UUID
	Login               string
	PasswordHash        string
	Email               string
	NotificationMethods []NotificationMethod
}
