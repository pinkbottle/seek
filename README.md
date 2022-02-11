### seek
seek lets you index various knowledge bases and seamlessly search it (powered by elasticsearch for now) in a unified manner

##### Example usage on indexed wiki sink
```
./seek sentence release process | head -n 18
https://en.wikipedia.org/wiki/PHP (13.181680)

Beginning on 28 June 2011, the PHP Development Team implemented a timeline for the release of new versions of PHP.[52] Under this system, at least one release should occur every month. Once per year, a minor release should occur which may include new features. Every minor release should at least be supported for two years with security and bug fixes, followed by at least one year of only security fixes, for a total of a three-year release process for every minor release. No new features, unless small and self-contained, are to be introduced into a minor release during the three-year release process.



https://en.wikipedia.org/wiki/Facebook (11.101337)

Facebook is developed as one monolithic application. According to an interview in 2012 with Facebook build engineer Chuck Rossi, Facebook compiles into a 1.5 GB binary blob which is then distributed to the servers using a custom BitTorrent-based release system. Rossi stated that it takes about 15 minutes to build and 15 minutes to release to the servers. The build and release process has zero downtime. Changes to Facebook are rolled out daily.[204]



https://en.wikipedia.org/wiki/PHP (10.579038)

Because of the major internal changes in phpng, it must receive a new major version number of PHP, rather than a minor PHP 5 release, according to PHP's release process.[52] Major versions of PHP are allowed to break backward-compatibility of code and therefore PHP 7 presented an opportunity for other improvements beyond phpng that require backward-compatibility breaks. In particular, it involved the following changes:
```

#### example CLI usage shared with WARP ❤️
https://app.warp.dev/block/7T7yQQkwk0Y31GzkTxAF9A