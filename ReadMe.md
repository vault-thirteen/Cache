# Cache

A cache mechanism with several features.

## Description

The cache stores information about the N most active records, where 'N' is a
variable number of records which depends on several factors. New records are
placed at the top of the cache, old records are removed from the bottom of the
cache, but not always. This reminds the _LRU_ cache but has more features.

## Variants

Due to limitations in Go programming language, there is no way to calculate 
memory usage of a generic variable. This leads to separation of cache models 
into following two variants:

* A cache with volume calculation which supports only `string` and `[]byte` variable types;
* A cache without volume calculation which supports `any` variable type.

The first variant is the most practical as it allows to control memory usage. 

The second variant allows to use different variable types for cached records, 
but it lacks the ability for controlling memory usage except for the indirect 
limitation by records' number.

The first variant is the base variant recommended for most situations where 
memory is limited. The second variant may be used if you are sure that size of 
each record will not exceed the expected size. 

Documentation is provided only for the first variant, as the second variant is 
a degraded version of it.
