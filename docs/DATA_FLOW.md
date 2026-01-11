📊 Request Flow Example
First Request (Cache Miss)
1. Client → GET /api/countries/search?name=India
2. Middleware → Log request, generate request ID
3. Handler → Parse "name=India", validate input
4. Service → Check cache for "country:india"
5. Cache → Returns ErrCacheMiss (not found)
6. Service → Call repository
7. Repository → HTTP GET to restcountries.com/v3.1/name/India
8. Repository → Parse response, create Country object
9. Service → Store in cache with 5min TTL
10. Service → Return Country to handler
11. Handler → Format JSON response
12. Middleware → Log completion
13. Client ← {"success":true,"data":{...}}

⏱️ Time: ~100-150ms (network call)
Second Request (Cache Hit)
1. Client → GET /api/countries/search?name=India
2. Middleware → Log request
3. Handler → Parse "name=India"
4. Service → Check cache for "country:india"
5. Cache → Returns Country object (HIT!)
6. Service → Return Country to handler
7. Handler → Format JSON response
8. Client ← {"success":true,"data":{...}}

⏱️ Time: ~2-5ms (in-memory read) - 30x faster! 🚀
