# Sitemap Generator
This application has 2 modes:

    - LocalConsole: This mode is meant for executing locally as a simple console application. 
        *** LocalConsole mode has to be run locally only, Docker was done explicitly for http server requests.
        To run LocalConsole just execute `go run .\main.go` within src/cmd directory. The configuration's default value is loaded as LocalConsole.
        For this mode the file will be generated and saved in the directory configured

    - http Server: This mode is for executing via http request. The Url to do this should be as following:
        http://localhost:{HTTP_PORT in .env}/?url={url}&maxDepth={max-depth}&xmlFileName={fileName}
        To run the HttpServer, check the `.env` file `EXECUTION_MODE` = 2, then run the `make start` command to start the docker container.
        For HttpServer, no file is being generated nor stored in the server, but the xml data is being returned along in the http response.

        URL Example: http://localhost:8080/?url=https://golangcode.com/get-a-url-parameter-from-a-request/&maxDepth=3&xmlFileName=sitemapv20220723

    PLEASE CHECK THE IMAGES WITHIN OUTPUT DIRECTORY

Unit tests:
    Only internal/sitereader has unit test due to time constraint from my end so please check that out and picture the same applies to the other left out functionality.

.env:
    I am aware third parties were not permitted but I cheatted a little bit here to use .env configuration, sorry >D!

    To run the HTTPServer, please create an .env file in the root and paste the data below: 
    *** configuration values are not intended to be public, more if we're talking about sensitive data, but for you guys to be able to run the app, and having in mind there are no tokens or security access credentials, I will leave them here.

    # 1 = LOCALCONSOLE / 2 = HTTPSERVER
    EXECUTION_MODE=1
    APP_NAME="Sitemaps"
    HTTP_PORT=8080
    URL_REG_EXPR=<a.*?href="(.*?)"
    OUTPUT_FILE="%s.xml"
    