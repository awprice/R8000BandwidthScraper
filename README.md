# R8000BandwidthScraper

A Golang command line tool for web scraping the device and bandwidth information from the Netgear R8000's web interface.

## Usage:

`./r8000bandwidthscraper <ROUTER IP> <AUTH KEY>`

The router ip is pretty self explanatory. It is usually the default gateway in your network.

The auth key is the username and password of your router's login information that has been base64 encoded. An example
of this:

1. Username/password combo: `admin:password`
2. Base64 encoded: `YWRtaW46cGFzc3dvcmQ=`

Note the use of the colon between the username and password. The key can also be found by logging into the admin
interface and inspecting the traffic between you and the router.

## Things to note

 - The R8000's web interface software only allows a request roughly every 2 seconds. The script will try and send as
 many requests as possible but is ultimately restricted by the router.
 - I have included a rather rudimentary method for clearing the terminal using the terminal's `clear` command. Pretty
  ugly, but it works.

