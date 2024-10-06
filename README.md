# goPR


## Config

Uses a config file at `~/.goPR.json` with the following format:
```json
{
  "ghtoken": "TOKEN",
  "repos": [
    {
      "org": "bebrws",
      "repo": "2DMessAround"
    }
  ]
}
```

## TODOs

Get Native OSX Notifications for GitHub Pull Requests (open PRs, comments, etc)

# TODO:
* Show multiple notifications for lots of changes?
* Document how to config

# Why not script terminal-notifier with the GH cli?
For example use terminal-notifier with a diff of `gh pr status`.
Yeah, I could have. But it is about the GoLang practice alright?


# Plist.info With STDOUT Logging
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>goPR</string>
	<key>CFBundleGetInfoString</key>
	<string>Created to notify on GH changes</string>	
	<key>CFBundleIdentifier</key>
	<string>com.bebrws.goPR</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.7</string>
	<key>CFBundleName</key>
	<string>goPR</string>
	<key>NSAppleEventsUsageDescription</key>
	<string>goPR requires sending apple events</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleVersion</key>
	<string>1000</string>
	<key>CFBundleShortVersionString</key>
	<string>Build 1000</string>
	<key>NSUserSelectedReadWriteAttribute</key>
	<array>
		<string>public.folder</string>
	</array>
</dict>
</plist>
```