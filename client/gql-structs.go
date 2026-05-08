package client

import (
  "time"
)

type RefreshTokenResponse struct {
  Data RefreshTokenBody `json:"refresh_token"`
}

type RefreshTokenBody struct {
  Token        TokenDetails        `json:"token"`
  RefreshToken RefreshTokenDetails `json:"refresh_token"`
  Errors       ExpectedError       `json:"errors,omitempty"`
  Success      bool                `json:"success"`
}

type ExpectedError map[string][]FieldError
type FieldError struct {
  Message string `json:"message"`
  Code string    `json:"code"`
}

type TokenAuthResponse struct {
  TokenAuth TokenAuthDetails `json:"token_auth"`
}

type TokenAuthDetails struct {
  RefreshToken RefreshTokenDetails  `json:"refresh_token"`
  Token TokenDetails                `json:"token"`
  Success bool                      `json:"success"`
  Errors *[]GraphQLError            `json:"errors"`
}

type GraphQLError struct {
  Message string            `json:"message"`
  Locations []ErrorLocation `json:"locations"`
}

type ErrorLocation struct {
  Line int    `json:"line"`
  Column int  `json:"column"`
}

type RefreshTokenDetails struct {
  Token string      `json:"token"`
  Created string    `json:"created"`
  ExpiresAt string  `json:"expires_at"`
}

type TokenDetails struct {
  Token   string       `json:"token"`
  Payload TokenPayload `json:"payload"`
}

type TokenPayload struct {
  Exp      time.Time  `json:"exp"`
  Username string     `json:"username"`
  OrigIat  time.Time  `json:"origIat"`
}

type EventListResponse struct {
  IPSCEvents []EventDetails           `json:"ip"`
  SteelChallengeEvents []EventDetails `json:"sc"`
}

type EventDetails struct {
  Id string                     `json:"id"`
  Name string                   `json:"name"`
  Starts time.Time              `json:"starts"`
  Ends *time.Time               `json:"ends"`
  State string                  `json:"state"`
  Status string                 `json:"status"`
  RegistrationStarts time.Time  `json:"registration_starts"`
  URL string                    `json:"get_full_absolute_url"`
  SubRule string                `json:"sub_rule"`
  Level string                  `json:"level"`
  UpdatedAt time.Time           `json:"updated_at"`
}
