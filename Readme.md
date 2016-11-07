# cgc

Crucible Game Counter

This project is to help me track the number of crucible games I play per day.
Here is an example of me querying my stats:

```
curl "https://cgc.herokuapp.com/games/2/4611686018448728582/2305843009297394104?from=2016-10-31"

{
   "from" : "2016-10-31",
   "total" : 86,
   "to" : "2016-11-07",
   "periods" : [
      {
         "date" : "2016-11-07",
         "count" : 3
      },
      {
         "date" : "2016-11-06",
         "count" : 23
      },
      {
         "date" : "2016-11-05",
         "count" : 12
      },
      {
         "count" : 12,
         "date" : "2016-11-04"
      },
      {
         "date" : "2016-11-03",
         "count" : 12
      },
      {
         "count" : 12,
         "date" : "2016-11-02"
      },
      {
         "date" : "2016-10-31",
         "count" : 12
      }
   ]
}
```

---

```
GET /games/{membershipType}/{membershipId}/{characterId}?from={yyyy-mm-dd}
```

* membershipType: 2 for PS4, 1 for XB1
* membershipId: Your account ID
* characterId: Your character ID

All of the above values can be taken from the URL of bungie.net when you're
looking at your character.

The `from` query parameter takes the format of 2016-10-10. This value is the
start date when the crucible games should start being counted.

Private matches, Trials of Osris, and Iron Banner are excluded from the results.
