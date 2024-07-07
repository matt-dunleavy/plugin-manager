### Configuration Options for redbean.ini

[redbean](https://redbean.dev/) is an open source webserver in a zip executable that runs on six operating systems.

##### Server Settings:

- `port`: The port number the server listens on (default: 8080)
- `addr`: The IP address to bind to (default: 0.0.0.0)
- `daemon`: Run as a daemon (true/false)
- `log`: Path to the log file
- `access_log`: Path to the access log file
- `error_log`: Path to the error log file
- `pid`: Path to the PID file
- `user`: User to run the server as
- `group`: Group to run the server as
- `chroot`: Directory to chroot into
- `ssl`: Enable SSL (true/false)
- `ssl_cert`: Path to SSL certificate file
- `ssl_key`: Path to SSL key file
- `ssl_password`: Password for SSL key file

1. ##### MIME Types:

   - `mime_type`: Set custom MIME types (e.g., `mime_type.xyz=application/x-xyz`)

2. ##### URL Rewriting:

   - `rewrite`: URL rewrite rules

3. ##### Directory Listings:

   - `dir_list`: Enable directory listings (true/false)
   - `dir_index`: Default index file names (comma-separated)

4. ##### CGI Settings:

   - `cgi_timeout`: Timeout for CGI scripts (in seconds)
   - `cgi_dir`: Directory for CGI scripts

5. ##### Lua Settings:

   - `lua_path`: Lua module search path
   - `lua_cpath`: Lua C module search path

6. ##### Security Settings:

   - `access_control_allow_origin`: Set CORS headers
   - `strict_transport_security`: Set HSTS header
   - `content_security_policy`: Set CSP header

7. ##### Performance Settings:

   - `workers`: Number of worker threads
   - `max_connections`: Maximum number of simultaneous connections
   - `keep_alive_timeout`: Keep-alive timeout (in seconds)
   - `gzip`: Enable gzip compression (true/false)
   - `gzip_types`: MIME types to compress (comma-separated)

8. ##### Caching:

   - `cache_control`: Set Cache-Control header
   - `etag`: Enable ETag header (true/false)

9. ##### Custom Error Pages:

   - `error_page`: Set custom error pages (e.g., `error_page.404=/custom_404.html`)

10. ##### Virtual Hosts:

    - `vhost`: Configure virtual hosts

11. ##### Proxy Settings:

    - `proxy_pass`: Configure reverse proxy settings

12. ##### WebSocket Settings:

    - `websocket`: Enable WebSocket support (true/false)

13. ##### Basic Authentication:

    - `auth_basic`: Enable basic authentication
    - `auth_basic_user_file`: Path to the user file for basic authentication

14. ##### Rate Limiting:

    - `limit_req`: Configure request rate limiting

15. ##### IP Filtering:

    - `allow`: Allow specific IP addresses or ranges
    - `deny`: Deny specific IP addresses or ranges

16. ##### File Serving:

    - `alias`: Create aliases for directories
    - `try_files`: Specify a list of files to try when serving requests

17. ##### Miscellaneous:

    - `server_tokens`: Control the emission of the Server header (on/off)
    - `client_max_body_size`: Maximum allowed size of the client request body

> [!IMPORTANT]
>
> Remember that the availability and syntax of these options may vary slightly depending on the version of redbean you're using. Always refer to the [official redbean documentation](https://redbean.dev/) for the most up-to-date and version-specific information.