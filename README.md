# Istina

This set of applications log the performance (and gameplay) of Garry's Mod servers into a single database for further analyzation and informational calculations.

The fetching of data happens via a GMod LUA addon. It gets sent via a HTTP POST request to a Golang webserver that saves the data into a MariaDB (MySQL). Every gmodserver sends a serverid which should be unique per gameserver. The server application is written to support many different gameservers sending data to it.

If you want to know more information about your players and the server itself then this could benefit you.  
This is not a logging application. It only saves performance and a few gameplay informations in regular intervals.

Example infos you could gather with this:

 - Most played jobs by playtime
 - Most played jobs by job-switch-count
 - Most playtime of players
 - Most used weapons
 - Average jointime
 - How many players have successfully joined or disconnected while loading
 - Playercount over time
 - Screensize / Operatingsystem / bitversion of players
 - RAM usage to detect leaks
 - Average FPS / Ping


# Installation

Go into the `server` folder and start the Golang application with:

    go get .
    go build .
    # copy and adjust the example config:
    cp config.example.yaml config.yaml
    nano config.yaml
    # then run the server with:
    ./istina


Move the folder inside the `lua_xyz_client` directory into your gmod servers' addons folder.  
Change the URLs at the top of the serverside file if needed.  


To also log LUA errors:  
Set the "lua_error_url" in your gmodserver.cfg (or server.cfg) to e.g. the following:

    lua_error_url "http://localhost:7077/luaerror"


To also log linux resource usage:  
Compile the Golang application inside the `linuxclient` folder, move it to the desired system and start it.


# Nginx config

It is a good idea to put your Go applications behind a reverse proxy like nginx if you want to make it accessible via the internet.

Put the following in your `server` config block:

        location /monitor/ {
                include ipwhitelist.conf
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_pass http://localhost:7077/;
        }

And the following into `/etc/nginx/ipwhitelist.conf`:

    allow 8.8.8.8;
    allow 9.9.9.9;
    deny all;

# Security

The above example nginx config has an ip whitelist. This is because of multiple factors:

 - The upload webserver has no authentication. Anyone who finds the URL could spam your database with useless data.
 - The lua_error_url is replicated to the client, which means every player on your server can easily find out your server URL and the other endpoints.
 - Other people could use your upload server and database without your knowledge if they get their hands on the addon and config of another server.

The easiest fix for this is the ip whitelist, because a passwort would not suffice if the whole addon of a server gets leaked and used on another server.


# Example queries

You can build your own frontend for this application which shows different statistics and details.

Example queries for this:

```sql
--Get the average jointime in seconds and ratio of joined/connected players
SELECT CONCAT(AVG(jointime),';',SUM(connected)/COUNT(connected)) FROM joinstats WHERE serverid = "xxxxxxxx" AND ts >= DATE_SUB(NOW(), INTERVAL 1 MONTH);
```
