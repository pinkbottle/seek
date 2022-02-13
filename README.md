# The idea behind seek 💡
`Seek` is a set of tools aimed for indexing and searching of any data

`Seek` works in a manner of Extract -> Transform -> Index

### Collectors 📦
`Collectors` are tools used for data collection, for example, a slack bot indexing a channel, jira scrapper indexing tickets for a given project or web crawler indexing websites internal to your organization containing valuable knowledge

### Seekers 🔎
`Seekers` are tools allowing you to query already indexed data that `Collectors` are producing, for now only `CLI` tool is supported

### Indexers 🗂
`Indexers` are tools used for data collection and transformation to an appropriate format and indexing that data in a search engine

## Notes 🌱
We're aiming for seek to be highly customizable and composable

For now the stack consists of `kafka` used for communication between `Collectors` and `Indexers`

For now we're using `elasticsearch` to manage the search index but ideally it should be trivial to use any other search engine and even combine them

## Examples 🧐
`SeekCLI` produces very simple output, source of the data matching your query, score (the higher the score, the better the match) and part of the indexed content

The idea is to use `seek` to search for the link to the source and do the exploring there

**Examples below show `seek` working only on a single data source (Wikipedia)**

> seek by sentence
```
./seek sentence release process | head -n 18
https://en.wikipedia.org/wiki/PHP (13.181680)

Beginning on 28 June 2011, the PHP Development Team implemented a timeline for the release of new versions of PHP.[52] Under this system, at least one release should occur every month. Once per year, a minor release should occur which may include new features. Every minor release should at least be supported for two years with security and bug fixes, followed by at least one year of only security fixes, for a total of a three-year release process for every minor release. No new features, unless small and self-contained, are to be introduced into a minor release during the three-year release process.



https://en.wikipedia.org/wiki/Facebook (11.101337)

Facebook is developed as one monolithic application. According to an interview in 2012 with Facebook build engineer Chuck Rossi, Facebook compiles into a 1.5 GB binary blob which is then distributed to the servers using a custom BitTorrent-based release system. Rossi stated that it takes about 15 minutes to build and 15 minutes to release to the servers. The build and release process has zero downtime. Changes to Facebook are rolled out daily.[204]



https://en.wikipedia.org/wiki/PHP (10.579038)

Because of the major internal changes in phpng, it must receive a new major version number of PHP, rather than a minor PHP 5 release, according to PHP's release process.[52] Major versions of PHP are allowed to break backward-compatibility of code and therefore PHP 7 presented an opportunity for other improvements beyond phpng that require backward-compatibility breaks. In particular, it involved the following changes:
```

> seek by keyword
```
./seek word debian | head -n 18
https://en.wikipedia.org/wiki/Debian (11.800510)

A large number of forks and derivatives have been built upon Debian over the years. Among the more notable are Ubuntu, developed by Canonical Ltd. and first released in 2004, which has surpassed Debian in popularity with desktop users;[249] Knoppix, first released in the year 2000 and one of the first distributions optimized to boot from external storage; and Devuan, which gained attention in 2014 when it forked in disagreement over Debian's adoption of the systemd software suite, and has been mirroring Debian releases since 2017.[250][251]



https://en.wikipedia.org/wiki/Debian (11.720029)

Debian was first announced on August 16, 1993, by Ian Murdock, who initially called the system "the Debian Linux Release".[10][11] The word "Debian" was formed as a portmanteau of the first name of his then-girlfriend (later ex-wife) Debra Lynn and his own first name.[12] Before Debian's release, the Softlanding Linux System (SLS) had been a popular Linux distribution and the basis for Slackware.[13] The perceived poor maintenance and prevalence of bugs in SLS motivated Murdock to launch a new distribution.[14]



https://en.wikipedia.org/wiki/Debian (9.510215)

Debian distribution codenames are based on the names of characters from the Toy Story films. Debian's unstable trunk is named after Sid, a character who regularly destroyed his toys.[9]

```

> example CLI usage shared with WARP ❤️

https://app.warp.dev/block/7T7yQQkwk0Y31GzkTxAF9A

https://app.warp.dev/block/tVcg4bLv15VeDp6LwDcyhT