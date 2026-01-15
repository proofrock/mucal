# μCal

I need a web application, very minimal, that connects to one or more CalDAV calendars and displays a read-only "task list" of the calendar items that are in a week, grouped by day. The current day ('today') should of course be highlighted. There should be a calendar view to select the week to display, with a mark on the days that have items. This calendar view will probably show the whole month, so these marks should be made for the whole visible month.

Both recurring and one-shot items must be viewed, with hours; and they must be in chronological order, with "all day" events at the top of a day. Items for different calendars must be "marked" in different colors, as configured; a left border stripe, for example, is ok. Also, an one-shot item that is "current" (i.e. the current date/time is inside the item's scheduled time span) should be subtly highlighted.

The Web UI must be responsive, aka being able to be rendered both on desktop and on mobile. The calendar thus implemented must be auto-refreshed, configuring the number of seconds.

As said, this is a read only view. No adding or changing calendar items. No reminders, also.

Errors must be logged to the stderr, and also shown in a popup in the Web UI.

## Architecture

There will be, of course, a frontend served in the browser, that will be in Svelte 5; for the UI, Bootstrap should be used. This frontend will connect to a backend written in Go, via REST methods. The frontend's resources must be embedded in the Go file.

CalDAV auth can be Basic.

No database should be needed, it just go to the caldav server. The server is local, so the network communication is fast.

Port must be fixed to 8080, HTTP. No need for HTTPS.

## Building and deploying

The way to install this must be via docker. A docker image must be provided.

CI using Github Actions must be provided, when a git tag is applied on the main branch on the GIT repo. The git tag is the version, e.g. 'v0.1.2'. The CI must build a image on the repo's ghcr.io docker space, with a docker tag equal to the git tag. Also, the same image with a docker tag of `latest` must be loaded.

Also, a way to build and test locally must be provided, via `make`; including a `cleanup` target to remove compilation caches/artifacts/dependencies, and a `update` target to upgrade the dependencies.

Semantic versioning must be used; the version - from the git tag - must be both in the docker image tag, as said, and shown in the Web UI.

## Configuration

The configuration must happen with a YAML file, passed read-only to the server via docker mapping. Example:

```yaml
time_zone: "Europe/Rome"
auto_refresh: 60 # Seconds
calendars:
  - name: "Birthdays"
    url: "https://.../birthdays"
    user_id: "mano"
    password_file: "/secrets/mypass1.txt"
    color: "#AABBCC"
  - name: "Personal"
    url: "https://.../personal"
    user_id: "mano"
    password_file: "/secrets/mypass2.txt"
    color: "#BBCCDD"
```

The color refers to what I said before, "Items for different calendars must be "marked" in different colors, as configured".

The password for the CalDAV communication is taken from a file, as shown. The text file will contain only the password, no newlines, whitespaces or structure of any kind.
