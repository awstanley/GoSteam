# GoSteam `webapi`

GoSteam `webapi` is a simple WebAPI wrapper for the Steam WebAPI.  At the moment is does not assume any knowledge of the return value, which have to be crafted by hand.

The file format recommended for responses is `<Interface>/<Method>Response.go`. (optionally with `V<Version>` prior to Response).

It is not advised to edit the files directly unless it's absolutely necessary; the auto-updater will clobber them.

## Usage

The API contains nothing but a core connection system (in `core`).  To get the API, and to update it, you need to use the utility application.

First you need to install the webapi updater:

    go install github.com/awstanley/GoSteam/webapi/apps/go-steam-webapi-updater

Then you need to update it using either a key:

    go-steam-webapi-updater --key="32CHARACTERSTEAMAPIKEYHERE"

Or a local JSON copy of `https://api.steampowered.com/ISteamWebAPIUtil/GetSupportedAPIList/v1`:

    go-steam-webapi-updater --file="/path/to/json/file.json"

Finally, import and use it as you will.  The one catch is almost no returns are currently handled; you'll need to write your own structs to handle the JSON.

**Warning**: The connection manager is designed to work without an API key, as is the updater.  If you don't pass a key it will generate the empty list.

Example apps will appear in `apps`.

## TODO

  * More response values (as I use them);
  * Helpers for various things (e.g. comma delimited lists of strings, ints, etc.).

## More information

For more information see: 

  * http://steamcommunity.com/dev
  * https://partner.steamgames.com/documentation/webapi
  * https://golang.org

## Licence

This project is licensed under the BSD 3-Clause licence (see LICENCE.md).