package client

import (
  "log"
  "context"
  "time"
  "reflect"
  "github.com/machinebox/graphql"
)

type SSIClient struct {
  APIKey string
  Token string
  RefreshToken string
  TokenExpiry time.Time
}

func NewClient(apikey string, username string, password string) *SSIClient {
  client := &SSIClient{
    APIKey: apikey,
  }

  authResponse := client.Auth(username, password)
  client.Token = authResponse.TokenAuth.Token.Token
  client.RefreshToken = authResponse.TokenAuth.RefreshToken.Token
  client.TokenExpiry = authResponse.TokenAuth.Token.Payload.Exp

  return client
}

func (c *SSIClient) refreshIfNeeded() {
  if time.Now().Add(30 * time.Second).After(c.TokenExpiry) {
    response := c.Refresh(c.RefreshToken)
    c.Token = response.Data.Token.Token
    c.TokenExpiry = response.Data.Token.Payload.Exp
    c.RefreshToken = response.Data.RefreshToken.Token
  }
}

func (c *SSIClient) Request(req *graphql.Request, respData interface{}) error {
  client := graphql.NewClient("https://shootnscoreit.com/graphql/")
  client.Log = func(s string) { log.Println(s) }
  ctx := context.Background()
  if err := client.Run(ctx, req, respData); err != nil {
    return err
  }
  return nil
}

func (c *EventDetails) IsEqualTo(event EventDetails) bool {
  aValue := reflect.ValueOf(*c)
  bValue := reflect.ValueOf(event)

  for i := 0; i < aValue.NumField(); i++ {
    fieldType := aValue.Type().Field(i)
    fieldName := fieldType.Name

    if fieldName == "UpdatedAt" {
      continue
    }

    if !reflect.DeepEqual(aValue.Field(i).Interface(), bValue.Field(i).Interface()) {
      return false
    }
  }
  
  return true
}

func (c *SSIClient) GetEvents() EventListResponse {
  c.refreshIfNeeded()

  isoDate := time.Now().Format("2006-01-02")
  req := graphql.NewRequest(`
    query($date: String!) {
      ip: events(starts_after: $date, status: "on", rule: "ip", firearms: "hg", region: "SWE") {
        id
        name
        starts
        ends
        state
        status
        registration_starts
        get_full_absolute_url
        sub_rule
        ... on IpscMatchNode {
          level
        }
      }
      sc: events(starts_after: $date, rule: "sc", firearms: "hg", region: "SWE") {
        id
        name
        starts
        ends
        state
        status
        registration_starts
        get_full_absolute_url
        sub_rule
        ... on SteelMatchNode {
          level
        }
      }
    }
  `)

  req.Var("date", isoDate)
  req.Header.Set("Authorization", "JWT " + c.Token)
  req.Header.Set("x-api-key", c.APIKey)

  var response EventListResponse
  if err := c.Request(req, &response); err != nil {
    log.Fatal(err)
  }

  return response
}

func (c *SSIClient) Refresh(refreshToken string) RefreshTokenResponse {
  req := graphql.NewRequest(`
    mutation($refreshToken: String!) {
      refresh_token(refresh_token: $refreshToken, revoke_refresh_token: false) {
        success
        errors
        token {
          payload {
            exp
            origIat
            username
          }
          token
        }
        refresh_token {
          token
          expires_at
        }
      }
    }
  `)

  req.Var("refreshToken", refreshToken)
  req.Header.Set("x-api-key", c.APIKey)

  var response RefreshTokenResponse
  if err := c.Request(req, &response); err != nil {
    log.Fatal(err)
  }

  if response.Data.Errors != nil {
    log.Fatal(response.Data.Errors)
  }

  return response
}

func (c *SSIClient) Auth(username string, password string) TokenAuthResponse {
  req := graphql.NewRequest(`
    mutation($username: String!, $password: String!) {
      token_auth(email: $username, password: $password) {
        refresh_token {
          token
          created
          expires_at
        }
        token {
          payload {
            exp
            username
            origIat
          }
          token
        }
        success
        errors
      }
    }
  `)

  req.Var("username", username)
  req.Var("password", password)
  req.Header.Set("x-api-key", c.APIKey)

  var response TokenAuthResponse
  if err := c.Request(req, &response); err != nil {
    log.Fatal(err)
  }

  if response.TokenAuth.Errors != nil {
    log.Fatal(response.TokenAuth.Errors)
  }

  return response
}
