# Cache

A Cache Mechanism with several Features.

## Description

The Cache stores Information about the N most active Records, where 'N' is a 
variable Number of Records which depends on several Factors. New Records are 
placed at the Top of the Cache, old Records are removed from the Bottom of the 
Cache, but not always. This reminds the LRU Cache but has more Features.

### TTL

Each Record has an Expiration Period, a Time to Live (TTL). If the requested 
Record is expired, it is removed from the Cache; if the Record is alive (i.e. 
not expired) when it is requested, it is then moved to the Top Position of the 
Cache. 

### Volume
Each Record also has a Volume. Volume is a Size of its Contents (Data) measured 
in Bytes. 

### Size

And, of course, the Size of each Record is One, just like One Piece of 
Something. Here we do not use "weighted" Records.

### Limits

The Cache has an ability to limit Records not only by their TTL, but also by 
the total Size of the Cache and total Volume of all the Records in the Cache.

While the TTL Limiter is always enabled, two other Limiters are optional and 
may be enabled or disabled in various Combinations. E.g., it is possible to 
limit both Size and Volume of the Cache. This makes the Cache a bit similar to 
a classical LRU Cache, but provides more Features and Flexibility. But do not 
be fooled by the Freedom of Choice, disabling both Size and Volume will make 
the Memory Usage uncontrollable.

To enable a Limiter, set its Value above Zero, and vice versa, to disable a 
Limiter, set its Value to Zero.

It is strongly advised to use this Cache for large Numbers of small Records 
instead of a small Number of gigantic Records while the latter may provoke OOM 
Exceptions. This Library guarantees to limit only the final Size of the Cache.
During the Process of Addition of new Records Memory Usage may increase above 
the Limits set in the Cache's Settings. 

And for God's Sake, do not set the TTL Limiter to Zero as it will totally 
disable the Cache. Zero TTL may only be used if you want to use the Cache as a 
Table storing Records' LATs, but for such Purposes it is weird.

### New Records

When a new Record is added into the Cache, if an "incoming" Record already 
exists in the Cache, i.e. its unique ID, UID, is already registered in the 
Cache, the Record is moved from its old Position in the Cache to the Top of the 
Cache. If a new Record has a unique unregistered ID, it is added to the Top.

### Old Records

Old Records are removed either from the Bottom of the Cache after the Insertion 
of a new Record, or when the requested Record is found to be expired or "stale".

The Removals are done in this "lazy" Style to save CPU Time. We check TTL only 
when it is necessary.

### Record Structure

Each Record has an 'UID' Field and a 'Data' Field. 'UID' is used for reading
and indexing Records. 'Data' is used to store some useful Information. 

Both Fields may have Dynamic types, called Generic Types in Go programming 
Language. More Information may be found in the Source Code.

### Requesting a Record

When a User requests a Record (by its UID) from the Cache, we first, check its 
Existence in the Cache's List, and then we check the Record's TTL (Time To 
Live). If the requested Record exists but is outdated, it is not returned to 
the User.

## Additional Notes

Due to some white Spaces in the modern State of the Go programming Language 
(Version 1.19), Generics (generic Types) are limited in their Usage, so that 
some Parts of the Code such as UID and Data Checks are not fully implemented at 
this Moment. When the Developers of Go Language fix their "bugs", this Library 
will be updated.

![Golang Logotype](img/golang-gopher-logotype.png)

Some additional Features may be added to the Cache in the Future, such as a 
Protection against Overheating, when same Records are moved to the Top Position 
too many Times in the same Second creating a useless load on the CPU.

## Importing

Import Commands:
```
import "github.com/vault-thirteen/Cache"
```
