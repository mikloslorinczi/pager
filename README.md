# pager
Pager Duty CLI helper

pager is a small command line interface to the [PagerDuty API](https://api-reference.pagerduty.com/#!/API_Reference/get_api_reference)

In order to use pager you need a [PagerDuty API Token](https://support.pagerduty.com/docs/generating-api-keys)\
You may supply it to pager with the **--api_token** flag (shorthand: **-t**) or with the **PAGER_API_TOKEN** env var

## Usage

```SHELL
$ pager (flags) [command] [sub-command]
```

**Available Commands:**

|Command|Description|
|-|-|
|oncall|On-call information|
|team|Team Information|
|user|User Information|
|help|Help about any command|

**Flags:**

|Shorthand|Flag|Type|Description|Default|
|-|-|-|-|-|
|-t|--api_token|string|Pager Duty API Token|
|-u|--api_url|string|Pager Duty API URL|https://api.pagerduty.com|
|-c|--config|string|Config file|
|-z|--time_zone|string|Time Zone in which results will be rendered|CET/CEST|

You can set the logging level with the **LOG_LEVEL** environment variable

## Examples

List all users
```SHELL
$ pager --api_token PD_TOKEN user list
```

Do the same with shorthands and aliases
```SHELL
$ pager -t PD_TOKEN u l
```

List all teams
```SHELL
$ pager --api_token PD_TOKEN team list
```

Check who is on-call
```SHELL
$ pager --api_token PD_TOKEN oncall
```

Get Sandor's on-call hours in this month
```SHELL
$ pager --api_token PD_TOKEN oncall hours --user Sandor
```

Get Miklós's on-call  hours from last month
```SHELL
$ pager -t PD_TOKEN oncall hours --user Miklós --last
```

Get Viktor's on-call  hours since 2019.07.01 until 2019.09.30
```SHELL
$ pager -t PD_TOKEN oncall hours --user Viktor --since 2019-07-01 --until 2019-09-30
```

Get the Infrastructure team member's on-call hours from last month
```SHELL
$ pager -t PD_TOKEN oncall hours --team Infrastructure --last
```

Do the same but with shorthands and aliases
```SHELL
$ pager -t PD_TOKEN oc h --team inf --last
```

## How on-call hours are calculated

When the **oncall** command invoked with the **hours** sub-command pager lists the on-call duties that fell into the configured timeframe. And calculates how many hour was OnWorkHour and how many was OffWorkHour amongst them, and reports the sum of each.

The PagerDuty API returns whole on-call duties when queried for a specific timeframe, these duties can reach out from the **since** and **until** limits of the timeframe. pager only considers the hours of these duties that are falling into the queried timeframe, it iterates over them hour by hour, checking if they are a working day's working hour and adds them to OnWorkHours or OffWorkHours accordingly.

pager relies on constant values **workHours** and **dayExceptions** in [rules/constants.go](rules/constants.go) when checking if a hour is a working hour and if a date is on a working day.