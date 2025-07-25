# Simple Chat2Desk API service

Chat2Desk API functions

# Functions
* Channels(ctx context.Context, offset, limit int) (*ChannelsResponse, error)
* GetChannels(ctx context.Context, offset, limit int) (*[]ChannelItem, error)

# Used libraries
* https://github.com/ra-company/env - Simple environment library (GPL-3.0 license)
* https://github.com/ra-company/logging - Simple logging library (GPL-3.0 license)
* https://github.com/stretchr/testify - Module for tests (MIT License)

# Staying up to date
To update library to the latest version, use go get -u github.com/ra-company/ctd.

# Supported go versions
We currently support the most recent major Go versions from 1.24.5 onward.

# License
This project is licensed under the terms of the GPL-3.0 license.