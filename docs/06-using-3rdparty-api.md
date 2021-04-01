# Using 3rd party APIs

Mattermost Apps framework provides services for using remote (3rd party) OAuth2
HTTP APIs, and receiving authenticated webhook notifications from remote
systems. There are 2 examples here to illustrate the [OAuth2](#hello-oauth2) and
[webhook](#hello-webhooks) support.

## Hello OAuth2!

Here is an example of an HTTP App ([source](/server/examples/go/hello-oauth2)),
written in Go and runnable on http://localhost:8080. 

- It contains a `manifest.json`, declares itself an HTTP application, requests
  permissions and binds itself to locations in the Mattermost UI.
- In its `bindings` function it declares 3 commands: `configure`, `connect`, and
  `send`.
- Its `send` function mentions the user by their Google name, and lists their
  Google Calendars.

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
Hello OAuth2! is an HTTP App, it requests the *permissions* to act as an Admin to change the App's OAuth2 config, as a User to connect and send. It binds itself to /commands.

https://github.com/levb/mattermost-plugin-apps/blob/81018cdcf3bb04b75a1fb20f5919c538ad14c73d/server/examples/go/hello-oauth2/manifest.json#L1-L17

## Bindings and Locations
The Hello OAuth2! creates 3 commands: `/helloworld configure|connect|send`.

```json
{
	"type": "ok",
	"data": [
		{
			"location": "/command",
			"bindings": [
				{
					"icon": "http://localhost:8080/static/icon.png",
					"description": "Hello remote (3rd party) OAuth2 App",
					"hint": "[configure | connect | send]",
					"bindings": [
						{
							"location": "configure",
							"label": "configure",
							"call": {
								"path": "/configure"
							}
						},
						{
							"location": "connect",
							"label": "connect",
							"call": {
								"path": "/connect"
							}
						},
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

### OAuth2 Call Handlers

To handle the OAuth2 "connect" flow, the app provides 2 Calls: `/oauth2/connect` that returns the URL to redirect the user to, and `/oauth2/complete` which gets invoked once the flow is finished, and the `state` parameter is verified.

```go
	// Handle an OAuth2 connect URL request.
	http.HandleFunc("/oauth2/connect", oauth2Connect)

	// Handle a successful OAuth2 connection.
	http.HandleFunc("/oauth2/complete", oauth2Complete)
```



#### configure

Sets up the Google OAuth2 credentials. Submit will require an Admin token to
make the changes.

```json
{
	"type": "form",
	"form": {
		"title": "Configures Google OAuth2 App credentials",
		"icon": "http://localhost:8080/static/icon.png",
		"fields": [
			{
				"type": "text",
				"name": "client_id",
				"label": "client-id",
				"is_required": true
			},
			{
				"type": "text",
				"name": "client_secret",
				"label": "client-secret",
				"is_required": true
			}
		],
		"call": {
			"path": "/configure",
			"expand": {
				"admin_access_token": "all"
			}
		}
	}
}
```


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
