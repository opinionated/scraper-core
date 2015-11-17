# scraper-core
Core code to scrape articles from sources. 

#####fetcher  
Functions to make and schedule requests. NOTE: this will get nixed and moved over to net

#####scraper 
[Scraper](scraper/) contains scraper interfaces, source implementations and funtions to run net requests.

#####net
[Net](net/) contains restful scraping client and server code. The server reads RSS feeds. Clients GET work from the server then POST the results.
