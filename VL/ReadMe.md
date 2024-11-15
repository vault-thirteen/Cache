# Cache

A cache mechanism with several features.

## Description

The cache stores information about the N most active records, where 'N' is a 
variable number of records which depends on several factors. New records are 
placed at the top of the cache, old records are removed from the bottom of the 
cache, but not always. This reminds the _LRU_ cache but has more features.

### TTL

Each record has an expiration period, a Time to Live (TTL). If the requested 
record is expired, it is removed from the cache; if the record is alive (i.e. 
not expired) when it is requested, it is then moved to the top position of the 
cache. 

### Volume

Each record also has a volume. Volume is a size of its contents (data) measured 
in bytes. 

### Size

And, of course, the size of each record is one, just like one piece of 
something. Here we do not use "weighted" records.

### Limits

The cache has an ability to limit records not only by their TTL, but also by 
the total size of the cache and total volume of all the records in the cache.

While the TTL limiter is always enabled, two other limiters are optional and 
may be enabled or disabled in various combinations. E.g., it is possible to 
limit both size and volume of the cache. This makes the cache a bit similar to 
a classical _LRU_ cache, but provides more features and flexibility. But do not 
be fooled by the freedom of choice, disabling both size and volume will make 
the memory usage uncontrollable.

To enable a limiter, set its value above zero, and vice versa, to disable a 
limiter, set its value to zero.

It is strongly advised to use this cache for large numbers of small records 
instead of a small number of large records while the latter may provoke OOM 
exceptions. This library guarantees to limit only the final size of the cache.
During the process of addition of new records memory usage may increase above 
the limits set in the cache's settings. 

And for God's sake, do not set the TTL limiter to zero as it will totally 
disable the cache. Zero TTL may only be used if you want to use the cache as a 
table storing records' LATs, but for such purposes it is weird.

### New Records

When a new record is added into the cache, if an "incoming" record already 
exists in the cache, i.e. its unique ID, UID, is already registered in the 
cache, the record is moved from its old position in the cache to the top of the 
cache. If a new record has a unique unregistered ID, it is added to the top.

### Old Records

Old records are removed either from the bottom of the cache after the insertion 
of a new record, or when the requested record is found to be expired or "stale".

The removals are done in this "lazy" style to save CPU time. We check TTL only 
when it is necessary.

### Record Structure

Each record has an 'UID' field and a 'Data' field. 'UID' is used for reading
and indexing records. 'Data' is used to store some useful information. 

Both fields may have dynamic types, called _Generic Types_ in _Go_ programming 
language. More information can be found in the source code.

### Requesting a Record

When a user requests a record (by its UID) from the cache, we first, check its 
existence in the cache's list, and then we check the record's TTL (Time To 
Live). If the requested record exists but is outdated, it is not returned to 
the user.

## Additional Notes

Due to some white spaces in the modern state of the _Go_ programming language 
(Version 1.19), _Generics_ (generic Types) are limited in their usage, so that 
some parts of the code such as UID and data checks are not fully implemented at 
this moment. When the developers of _Go_ language fix their "bugs", this library 
will be updated.

![Golang Logotype](../img/golang-gopher-logotype.png)

## Performance

The first idea was to limit the process of LAT updates in order to reduce the 
CPU usage, so that "hot" records which had previous LAT updates less than a 
second ago, would not get a new LAT update. However, the stress tests show that 
performance is great even without such limitations.

Stress tests on quite a decent hardware show RPS rate of about 24 to 25 
million requests per second in the most heavy test (1000 records, each 
having 1MB of data). This means that no overheating protection is required. 
Probably, this cache will not ever be a bottleneck.

## Importing

Import Commands:
```
import "github.com/vault-thirteen/Cache"
```
