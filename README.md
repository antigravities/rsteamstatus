# /r/Steam status bot

This is the status bot we use on /r/Steam to update the sidebar. It's automatically run every 15 minutes as a cron job in /r/Steam's container.

The data is sourced from a secret place. You'll need to know the URL of the secret place in order to use this bot.

## Setup
You'll need Go to run the bot. The bot *requires* Go 1.16. The last commit to have Go 1.15 support is 8dd9b8.

```
GO111MODULE=on go get get.cutie.cafe/rsteamstatus
$EDITOR .env #described below
~/go/bin/rsteamstatus
```

## .env
.env contains the configuration for the `rsteamstatus` tool. You can either use a .env file, or specify the values directly via environment variables. Here are the variables to get you started:

```
R_STATUS_URL=
R_CLIENT_ID=
R_CLIENT_SECRET=
R_USERNAME=
R_PASSWORD=
R_SUBREDDIT=
R_USER_AGENT=
```

`R_USER_AGENT` and your outgoing IP must be whitelisted by the secret place. `R_USER_AGENT` is also the `User-Agent` the reddit bot uses (good example: `rsteamstatus/X.X (/u/your_username; your@email.tld)`).

## License
```
Copyright (C) 2020-2021 Alexandra Frock

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
