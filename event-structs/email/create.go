package eventstructs

import "time"

type UserCreate struct {
	Event struct {
		SessionVariables struct {
			XHasuraRole string `json:"x-hasura-role"`
		} `json:"session_variables"`
		Op   string `json:"op"`
		Data struct {
			Old interface{} `json:"old"`
			New struct {
				Email                interface{} `json:"email"`
				DisplayName          string      `json:"display_name"`
				Active               bool        `json:"active"`
				UpdatedAt            time.Time   `json:"updated_at"`
				SecretTokenExpiresAt time.Time   `json:"secret_token_expires_at"`
				SecretToken          string      `json:"secret_token"`
				CreatedAt            time.Time   `json:"created_at"`
				ID                   string      `json:"id"`
				AvatarURL            interface{} `json:"avatar_url"`
				DefaultRole          string      `json:"default_role"`
			} `json:"new"`
		} `json:"data"`
	} `json:"event"`
	CreatedAt    time.Time `json:"created_at"`
	ID           string    `json:"id"`
	DeliveryInfo struct {
		MaxRetries   int `json:"max_retries"`
		CurrentRetry int `json:"current_retry"`
	} `json:"delivery_info"`
	Trigger struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}
