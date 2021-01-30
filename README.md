<b>This repository is no longer maintained.</b>

`slack-to-telegram` is a simple way to forward notifications from Slack to Telegram when you are not logged into Slack. It uses Slack's [Real Time Messaging API](https://api.slack.com/rtm) to connect to specified accounts and will forward messages through Telegram's [Bot API](https://core.telegram.org/bots/api) to you.

## Why 
I currently run an Android OS on my phone without Google services. Slack for Android requires Google Cloud Messaging (GCM) in order to receive notifications. Since Telegram has its own mechanism for notifications on Android that does not rely on GCM, I decided to use this as a workaround to receive timely Slack notifications.

## Features
* Support for multiple Slack workspaces
* Get notifications for direct messages and @username callouts in channels

## One Time Initial Setup
* From the Telegram account you want to receive messages on, get your user ID by sending a message to `@get_id_bot`
* Create a telegram bot (https://core.telegram.org/bots#3-how-do-i-create-a-bot) and get the token
* Send a test message from your Telegram account to your bots username
* For each Slack account you want to forward messages, you'll need access to the RTM API which previously was possible with legacy tokens but this has been disabled now. Instead you can use [this approach](https://github.com/wee-slack/wee-slack#get-a-session-token) which doesn't require installing an app in your organization. Or if you can't get that to work, you can use [slack-rtm-token](https://github.com/dan-v/slack-rtm-token) to get an OAuth token with appropriate access, but this will require installing an app in each organization.

## Config File
Create a file named config.toml and fill in the details from initial setup above.

    [telegram]
    user = 123456789 
    token = "323456789:ABCDE_fB19OHQZUF3FPPPF43PTEEB"

    [[slack]]
    name = "workspace #1"
    token = "xoxp-xxxxxxx-xxxxxxxx-xxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

    [[slack]]
    name = "workspace #2"
    token = "xoxp-xxxxxxx-xxxxxxxx-xxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

## Installation

### Binaries
The easiest way is to download a pre-built binary from the [GitHub Releases](https://github.com/dan-v/slack-to-telegram/releases) page.

### Docker
You can also just run it as a docker container.

    docker run --restart=always -d -v $(pwd)/config.toml:/config.toml vdan/slack-to-telegram:latest

## Usage

    ./slack-to-telegram --config config.toml

## FAQ
1. <b>Should I use slack-to-telegram?</b> That's up to you. Use at your own risk.

## Powered by
* Slack API ([slack-go/slack](https://github.com/slack-go/slack))
* Telegram Bot API ([telegram-bot-api.v4](https://gopkg.in/telegram-bot-api.v4))
* TOML parser ([BurntSushi/toml](https://github.com/BurntSushi/toml))

## Build From Source

    make tools && make
