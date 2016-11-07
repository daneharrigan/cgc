# cgc

Crucible Game Counter

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
