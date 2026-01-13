
# SSI Calendar

This is a project aimed at 🇸🇪 Swedish IPSC shooters (pun intended) 🔫

The service exposes a calendar ics file with both registration and competition events from "Shoot-N-Score It" so that shooters can easily see when an events occurs or when registration opens in their calendar as well as get built-in notifications from their calendar apps.

There are a few filters applied at the moment:
 - Level > 1, meaning it will not show level 1 events
 - Has to be a handgun event
 - Has to be in Sweden

 **Note** this is simply an aid to remember when an event occurs. You still have to sign up yourself!
# How to use

This project is intended to be used together with your smartphone calendar.  
You should subscribe to it in your calendar.

Example iOS / iPhone:  
Settings -> Apps -> Calendar -> Accounts -> Add account -> Subscribed calendar

You can achieve the same thing on Android with your calendar of choice.  
The phone will periodically query the calendar and add it and your default notification policy should alert automatically if enabled.
## Endpoints

| Endpoint | Description                       |
| :-------- | :-------------------------------- |
| `/calendar.ics`      | Combined competitions + registrations |
| `/competitions.ics`      | Competitions only (dates of competition) |
| `/registrations.ics`      | Registration dates of competititons only |


## Live test

The calendar is currently live at:  
https://ssi-calendar.home.sigkill.me/calendar.ics

Please note that any abuse will result in service termination.

The service may stop working at any time should I decide to stop hosting it for any reason.


## Roadmap

- Per-user calendar with custom settings (IPSC levels, Distance etc.)
- Include squadding dates for events you've signed up to
- Use users own API keys instead of a single one


## Feedback

If you have any feedback, please visit Bromma PK discord at https://discord.gg/jqV9kf8wUR

