# SSI Calendar
This service automatically provides a calendar with IPSC events and their corresponding registration dates, A user may subscribe to the calendar for automatic updates.

## How it works
The generated calendar will only include shootnscoreit (SSI) events for Handgun IPSC events that are level 2 or higher located in Sweden.

## How to use
Use any calendar either on your phone or computer and subscribe to the generated ICS file, it is served via the /calendar.ics endpoint

## V2 implementation plan
### Auth
We need an authentication system to allow users to log in where they can store their api key and maybe even preferences

### Per-User API key
We would need to store an API key for every user and using that we could then query for events they have signed up for in order to either filter out (if the data is already present) or fetch the necessary data in order to just show the relevant data

### Unique link
The user should get a unique link that is linked to the user which returns a tailed ical file.

### Frontend
A basic frontend that can interact with the backend will be required in order to provide a UI to login and set API key

### Data to be shown
Event registration, show all.
Event date, Show if registered.
Event squadding, if registered.

### Future
Perhaps more customization like filter by locations or allow other events like shotgun or rifle matches.

