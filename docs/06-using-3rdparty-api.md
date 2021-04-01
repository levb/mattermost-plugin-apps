# Using 3rd party APIs

Mattermost Apps framework provides assistance for using remote (3rd party) OAuth2 HTTP APIs, and receiving authenticated webhook notifications from remote (3rd party) systems. There are 2 examples here to illustrate the OAuth and webhook support.

# Hello OAuth2!

Here is an example of an HTTP App ([source](/server/examples/go/hello-oauth2)),
written in Go and runnable on http://localhost:8080. 

- It contains a `manifest.json`, declares itself an HTTP application, requests
  permissions and binds itself to locations in the Mattermost UI.
- In its `bindings` function it declares 3 commands: `configure`, `connect`, and
  `send`.
- It contains a `send` function that sends a parameterized message back to the
  user. 
- It contains a `send-modal` function that forces displaying the `send` form as
  a Modal.

To install "Hello, OAuth2" on a locally-running instance of Mattermost follow
these steps (go 1.16 is required):
```sh
cd .../mattermost-plugin-apps/server/examples/go/hello-oauth2
go run . 
```

In Mattermost desktop app run:
```
/apps debug-add-manifest --url http://localhost:8080/manifest.json
/apps install --app-id hello-oauth2
```

You need to configure your [Google API
Credentials](https://console.cloud.google.com/apis/credentials) for the App. Use
`{MattermostSiteURL}/com.mattermost.apps/apps/hello-oauth2/oauth2/remote/complete`
for the `Authorized redirect URIs` field. After configuring the credentials, in Mattermost desktop app run:
```
/hello-oauth2 configure --client-id {ClientID} --client-secret {ClientSecret}
```

Now, you can connect your account to Google with `/hello-oauth2 connect` command, and then try `/hello-oauth2 send`.

## Manifest
The manifest declares App metadata, For HTTP apps like this no paths mappings
are needed. The Hello World App requests the *permission* to act as a Bot, and
to *bind* itself to the channel header, and to /commands.

```json
{
	"app_id": "helloworld",
	"display_name": "Hello, world!",
	"app_type": "http",
	"root_url": "http://localhost:8080",
	"requested_permissions": [
		"act_as_bot"
	],
	"requested_locations": [
		"/channel_header",
		"/command"
	]
}
```

## Bindings and Locations
Locations are named elements in Mattermost UI. Bindings specify how App's calls
should be displayed at, and invoked from these locations. 

The Hello App creates a Channel Header button, and adds a `/helloworld send` command.

```json
{
	"type": "ok",
	"data": [
		{
			"location": "/channel_header",
			"bindings": [
				{
					"location": "send-button",
					"icon": "http://localhost:8080/static/icon.png",
					"label":"send hello message",
					"call": {
						"path": "/send-modal"
					}
				}
			]
		},
		{
			"location": "/command",
			"bindings": [
				{
					"icon": "http://localhost:8080/static/icon.png",
					"description": "Hello World app",
					"hint":        "[send]",
					"bindings": [
						{
							"location": "send",
							"label": "send",
							"call": {
								"path": "/send"
							}
						}
					]
				}
			]
		}
	]
}
```

## Functions and Form
Functions handle user events and webhooks. The Hello World App exposes 2 functions:
- `/send` that services the command and modal.
- `/send-modal` that forces the modal to be displayed.

```go
func main() {
	// Serve its own manifest as HTTP for convenience in dev. mode.
	http.HandleFunc("/manifest.json", writeJSON(manifestData))

	// Returns the Channel Header and Command bindings for the App.
	http.HandleFunc("/bindings", writeJSON(bindingsData))

	// The form for sending a Hello message.
	http.HandleFunc("/send/form", writeJSON(formData))

	// The main handler for sending a Hello message.
	http.HandleFunc("/send/submit", send)

	// Forces the send form to be displayed as a modal.
	// TODO: ticket: this should be unnecessary.
	http.HandleFunc("/send-modal/submit", writeJSON(formData))

	// Serves the icon for the App.
	http.HandleFunc("/static/icon.png", writeData("image/png", iconData))

	http.ListenAndServe(":8080", nil)
}

func send(w http.ResponseWriter, req *http.Request) {
	c := apps.CallRequest{}
	json.NewDecoder(req.Body).Decode(&c)

	message := "Hello, world!"
	v, ok := c.Values["message"]
	if ok && v != nil {
		message += fmt.Sprintf(" ...and %s!", v)
	}
	mmclient.AsBot(c.Context).DM(c.Context.ActingUserID, message)

	json.NewEncoder(w).Encode(apps.CallResponse{})
}
```

The functions use a simple form with 1 text field named `"message"`, the form
submits to `/send`.

```json
{
	"type": "form",
	"form": {
		"title": "Hello, world!",
		"icon": "http://localhost:8080/static/icon.png",
		"fields": [
			{
				"type": "text",
				"name": "message",
				"label": "message"
			}
		],
		"call": {
			"path": "/send"
		}
	}
}
```

## Icons 
Apps may include static assets. At the moment, only icons are used.
