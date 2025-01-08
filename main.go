package main

import (
  "io"
  "encoding/json"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
  "github.com/arran4/golang-ical"
  "ssi-calendar/client"
  "ssi-calendar/storage"
)


/*
	Make a graphql client with supported methods
	like auth, renewing token, get events

	Process:
	1. Check if DB has auth creds
	1.1 If Not, Fetch via auth endpoint
	1.2 Store in DB
	2. Check if token expired, if so renew with refresh token
	3. Get all events
	4. Update their times (event id, event start date, event end date, reg time, squad time, event name) in DB

	Make a HTTP server that can serve an ical file
	upon request the ical should have *all* of the events, if current date > event end, do not include it
	Either mark expired events at this stage or have a cron to do this check so we dont update the DB when fetching cal

  Some sort of cron/scheduling is needed to periodically fetch events and update the DB
  The actual HTTP request that generates a calendar should just do it through the data in the DB
*/

var db *storage.Storage
var ssiClient *client.SSIClient

func updateEvents() {
  eventsResponse := ssiClient.GetEvents()
  for _, event := range eventsResponse.Events {
    db.UpdateEvent(event)
  }
}

func startPeriodicTask(interval time.Duration, task func()) {
  ticker := time.NewTicker(interval)
  go func() {
    for {
      select {
      case <-ticker.C:
        task()
      }
    }
  }()
}

func getRoot(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "SSI Calendar 1.0\n")
}

func getEvents(w http.ResponseWriter, r *http.Request) {
  events := db.GetEvents()
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(events)
}

func doUpdate(w http.ResponseWriter, r *http.Request) {
  updateEvents()
  io.WriteString(w, "OK\n")
}

func getCalendar(w http.ResponseWriter, r *http.Request) {
  events := db.GetEvents()
  cal := ics.NewCalendar()
  cal.SetMethod(ics.MethodRequest)
  for _, event := range events {
    if time.Now().After(event.Ends) {
      continue
    }
    cevent := cal.AddEvent(event.Id)
    cevent.SetCreatedTime(event.Starts)
    cevent.SetDtStampTime(event.Starts)
    cevent.SetModifiedAt(event.UpdatedAt)
    cevent.SetStartAt(event.Starts)
    cevent.SetEndAt(event.Ends)
    cevent.SetSummary(event.Name)
    cevent.SetURL("https://shootnscoreit.com/event/22/" + event.Id + "/")

    if time.Now().Before(event.RegistrationStarts) {
      revent := cal.AddEvent("reg" + event.Id)
      revent.SetCreatedTime(event.RegistrationStarts)
      revent.SetDtStampTime(event.RegistrationStarts)
      revent.SetModifiedAt(event.UpdatedAt)
      revent.SetStartAt(event.RegistrationStarts)
      revent.SetEndAt(event.RegistrationStarts.Add(15 * time.Minute))
      revent.SetSummary("Registration opens " + event.Name)
      revent.SetURL("https://shootnscoreit.com/event/22/" + event.Id + "/")
    }
  }
  w.Header().Set("Content-Type", "text/calendar")
  io.WriteString(w, cal.Serialize())
}

func main() {
  // Set up storage
  db = storage.NewStorage()
  defer db.Close()

  // Set up SSI client
  key := os.Getenv("SSI_APIKEY")
  ssiClient = client.NewClient(key)

  // Periodically fetch new events
  updateEvents() // Do an initial fetch
  startPeriodicTask(3*time.Hour, func() {
    updateEvents()
  })

  // Set up HTTP server
  // Also handle SIGTERM so we can let the defers do their thing
  http.HandleFunc("/", getRoot)
  http.HandleFunc("/events", getEvents)
  http.HandleFunc("/update", doUpdate)
  http.HandleFunc("/calendar.ics", getCalendar)

  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

  go func() {
    port := os.Getenv("PORT")
    http.ListenAndServe(":" + port, nil)
  }()

  <-stop
}
