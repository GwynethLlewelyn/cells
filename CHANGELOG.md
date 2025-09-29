# Changes between v4.9.92-alpha09 and v4.9.92-alpha11

[See Full Changelog](https://github.com/pydio/cells/compare/v4.9.92-alpha09...v4.9.92-alpha11)

- [#a24b68c](https://github.com/pydio/cells/commit/a24b68cbdc136efd0d345a52703377b05e70fe0b): chore(deps): bump ClickHouse lib
- [#629d735](https://github.com/pydio/cells/commit/629d735b5bfb14172d955b9806ebf4d85d1fcc3b): chore(deps): fix go.sum merge
- [#acea3c7](https://github.com/pydio/cells/commit/acea3c7ea71d4e151902635c0a59dc6d6e0c2ff1): Merge remote-tracking branch 'origin/v5-dev' into v5-dev
- [#4d72af5](https://github.com/pydio/cells/commit/4d72af58d199a29b239a1ef5fded676b3bcbab69): chore(deps): bump most dependencies, including caddy but excluding ory/hydra. Switch to go1.25.
- [#f887f55](https://github.com/pydio/cells/commit/f887f55077b75cfa013ab45ca7f2ceacfffe657c): build(tools): chart beta 15
- [#af22869](https://github.com/pydio/cells/commit/af228696585525aa8f7b3da36e788578b3f8f488): build(tools): chart beta 14
- [#fad6fbe](https://github.com/pydio/cells/commit/fad6fbed3deae90df31857def3a2a74bfa818bf5): build(tools): chart beta 14
- [#15498ea](https://github.com/pydio/cells/commit/15498ea702d54da6e7f1651664872c903c261d3f): Revert "fix(api): prevent setting unary space tag value"
- [#a85ff9b](https://github.com/pydio/cells/commit/a85ff9ba13c1df5fd1084096e11cb5677cfc92a9): fix(search): fix search config propagation inside indexer.
- [#d077ec6](https://github.com/pydio/cells/commit/d077ec63a2b3255fd67340c0370fc5bfac674a53): fix(cors): remove logging
- [#ece94c3](https://github.com/pydio/cells/commit/ece94c367f97ae324b18151f1949a5129e243186): feat(i18n): more messages + DE/FR translations
- [#81ef987](https://github.com/pydio/cells/commit/81ef98744a9d0e51f70b6d9221fda4b916863089): fix(start): multi node issue
- [#d21301a](https://github.com/pydio/cells/commit/d21301ac3aca40810a9c30079c6bc268de9f3d75): fix(ux): FilePreview processing positioning
- [#1fc54eb](https://github.com/pydio/cells/commit/1fc54eb06ffdc9f81c26f128a3a3ccbd536813ec): feat(bnote): hint to insert ToC
- [#4219d36](https://github.com/pydio/cells/commit/4219d361c8ed39641e4bfb4a230881e536f0f47e): feat(bnote): properly install bnote namespaces at first run (retry if ns service is not yet migrated)
- [#ba263a1](https://github.com/pydio/cells/commit/ba263a148fa270c722bc310d20bc46f8f83dfdac): fix(ux): List V2 - fix column blinking when showing/hiding extension, fix thumbs/masonry gutters for "small" mode.
- [#a37f4c7](https://github.com/pydio/cells/commit/a37f4c716c86600a03623f6f68e98335220d8c0d): fix(ux): fix grouping headers in new list
- [#adbaabe](https://github.com/pydio/cells/commit/adbaabe8d67133feea0aeb6381247d654a52785b): feat(audio): fix soundmanager positioning issue
- [#3a17abe](https://github.com/pydio/cells/commit/3a17abe19f4e0265c0e3cec5707f15528621d80a): fix(sql): refix user_tree migration for PG
- [#0e37e5f](https://github.com/pydio/cells/commit/0e37e5fc77d0f3d0fbbd5007e67d115796096274): feat(ux): generalize and style CustomDragLayer for modern list
- [#2547db0](https://github.com/pydio/cells/commit/2547db01c1f1d9b6f4192872cf97ceed55cf0f6c): feat(pages): create default namespaces for Pages feature.
- [#4a2025f](https://github.com/pydio/cells/commit/4a2025fbad60bed69ef5e5808ca25feed2e94d0d): fix(sql): fix collation mechanism for index tables: leave table to db default and patch name, mpath and hash columns. This should prevent migration issues. Fix also migration of new policies (and insertion of new default for api v2).
- [#98bcc13](https://github.com/pydio/cells/commit/98bcc133740e4655885b3a52f4d931c742032046): fix(api): prevent setting unary space tag value
- [#2145861](https://github.com/pydio/cells/commit/214586105e08ebe93dc01b83097a31cdbf6ede89): feat(search): improve homepage search engine as a "show more" button was missing there.
- [#05af7d5](https://github.com/pydio/cells/commit/05af7d58344d4b512f95ccb44aa26f42f4dddeae): fix(log): lower level for api v2 logs
- [#5752504](https://github.com/pydio/cells/commit/5752504d91ef2059e589499e542a90daa692afa9): fix(sql): Fix GetNodeChildrenCounts that must return different results whether it's recursive or not.
- [#cea6603](https://github.com/pydio/cells/commit/cea6603e9fd0e35cf7602ba9f81e69fa9167eafd): fix(tests): fix unit tests after last changes
- [#72e5167](https://github.com/pydio/cells/commit/72e5167139c83a3d2f3bbf9214dfcd2de0b62856): fix(pprof): debug endpoint should not be enabled by default.
- [#9b947bb](https://github.com/pydio/cells/commit/9b947bbd6b656a8a71a2c34e04d050ee4d0c756e): fix(leak): fix goroutine leak (and possible underlying mem leak with subscriptions) in SearchNodes by pooling NsProvider initialization.
- [#8689d6e](https://github.com/pydio/cells/commit/8689d6ef18774381239e2fabafa52d2ac75c62b5): fix(logs): make datasource health check logs less verbose, consume rest.install handlers as core context.
- [#aadbb2a](https://github.com/pydio/cells/commit/aadbb2a54bffa49d14cb5c72e961761cdfa426fd): fix(policies): Add v4=>v5 migration to enable access to api v2.
- [#ca7389a](https://github.com/pydio/cells/commit/ca7389abf13a0ddef075eb6f5e1fce0c4a789114): chore(i18n): update admin page title
- [#22109f2](https://github.com/pydio/cells/commit/22109f2e7135f189deb100b52863b39dc3d08659): feat(sites): Improve log error when site's external URL is not properly matching incoming request.
- [#61fb37b](https://github.com/pydio/cells/commit/61fb37beaae43778c84572cfbf9eb9759b0fd7f2): fix(restore): Properly clean RecycleRestore metadata (original source path) during restoration from recycle bin task.
