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
}

func NewClient(apikey string) *SSIClient {
  return &SSIClient{APIKey: apikey}
}

func (c *SSIClient) Request(req *graphql.Request, respData interface{}) error {
  client := graphql.NewClient("https://shootnscoreit.com/graphql/")
  //client.Log = func(s string) { log.Println(s) }
  ctx := context.Background()
  if err := client.Run(ctx, req, respData); err != nil {
    return err
  }
  return nil
}

func (c *EventListResponse) FilterEvents() {
  filteredEvents := []EventDetails{}

  for _, event := range c.Events {
    if (event.SubRule == "nm") && (event.Level != "l1") {
      filteredEvents = append(filteredEvents, event)
    }
  }

  c.Events = filteredEvents
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
  isoDate := time.Now().Format("2006-01-02")
  req := graphql.NewRequest(`
    query($date: String!) {
      events(starts_after: $date, status: "on", rule: "ip", firearms: "hg", region: "SWE") {
        id
        name
        starts
        ends
        state
        status
        registration_starts
        sub_rule
        ... on IpscMatchNode {
          level
        }
      }
    }
  `)

  req.Var("date", isoDate)
  req.Header.Set("Authorization", "JWT " + c.APIKey)

  var response EventListResponse
  if err := c.Request(req, &response); err != nil {
    log.Fatal(err)
  }

  response.FilterEvents()
  return response
}

func (c *SSIClient) Renew(refreshToken string) {
  req := graphql.NewRequest(`
    mutation($refreshToken: String!) {
      refresh_token(refresh_token: "$refreshToken", revoke_refresh_token: false) {
        success
        errors
        token {
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
          token
        }
        success
        errors
      }
    }
  `)

  req.Var("username", username)
  req.Var("password", password)

  var response TokenAuthResponse
  if err := c.Request(req, &response); err != nil {
    log.Fatal(err)
  }

  if response.TokenAuth.Errors != nil {
    log.Fatal(response.TokenAuth.Errors)
  }

  return response
}
