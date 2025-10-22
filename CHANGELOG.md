# Changes between v4.9.92-alpha12 and v4.9.92-alpha13

[See Full Changelog](https://github.com/pydio/cells/compare/v4.9.92-alpha12...v4.9.92-alpha13)

- [#f05cc39](https://github.com/pydio/cells/commit/f05cc39cc374ef22d592d96d8ab6be69c82fdc2c): fix(logs): CLICKUP #869atbfe1 - Fix logs related flags (level, encoding, to file) in a more logical way : flags are only related to stdout, and if loggers are passed via config, flags are overriding stdout existing values.
- [#279b49e](https://github.com/pydio/cells/commit/279b49e28d1075d5ce921925bde0b3ddbae32a77): fix(idm): fix user roles ordering not persisted
- [#c8fcff0](https://github.com/pydio/cells/commit/c8fcff0557ee7770582b5a7eafc0a185fdc3faeb): Merge remote-tracking branch 'origin/v5-dev' into v5-dev
- [#d5beec0](https://github.com/pydio/cells/commit/d5beec0d3f112a7d1c783f019957112a51b81200): fix(upgrade): fix v4 to v5 upgrade glitches
- [#bb03f0b](https://github.com/pydio/cells/commit/bb03f0bec31815f217b056efe368e4a36c005074): fix(sql): raise defaults for mysql connections
- [#a5ccfdb](https://github.com/pydio/cells/commit/a5ccfdbc8423a4227eb8888fef5feb2406c8374f): fix(thumbs): WPB-21131 generate thumbs with the exif orientation (#686)
- [#b2c37dd](https://github.com/pydio/cells/commit/b2c37dd67fe06992832ecce030ca6f9e6fdc3fe8): fix(cli): silently fails when quiet flag is set
- [#fe040bb](https://github.com/pydio/cells/commit/fe040bbb88c2a1031074d58d6a755c49fd4f9f96): fix(context): pass context to Init() interface and get rid of RuntimeContext pattern.
- [#63f54bc](https://github.com/pydio/cells/commit/63f54bcd340192da011db1ef5ee53b7ea78ec7ee): fix(context): pass context to GetParametersForm() interface so that some actions implementations can use config.
