### hooThat

A small application to analyse www logs to determine number of unique IP hits and country of origin. The purpose is to be able to remove analytics from the site while being able to keep an eye on these basic stats.

1. Retrieve IPs from www access log.
1. Run the IPs through a geo IP DB.
1. Direct the parsed to a basic html page.
1. Repeat for X amount of logs so there's view of X days of traffic.

Version 1 is very basic and my first attempt at Golang, I hope to improve the code as I work my way through learning the language.

Items I'd like to add in the future:

1. Some form of graph, even if it's just ASCII that would be fine.
1. Retain historic data, including a useful way to display.
1. Report on stats of specific pages of the site.
1. Improve code.