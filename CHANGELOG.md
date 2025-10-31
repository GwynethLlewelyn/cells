# Changes between v4.9.92-alpha13 and v4.9.92-alpha14

[See Full Changelog](https://github.com/pydio/cells/compare/v4.9.92-alpha13...v4.9.92-alpha14)

- [#53fa497](https://github.com/pydio/cells/commit/53fa49796457a107f1c14f1b02bcf06776f51e5e): feat(i18n): more Norvegian translations
- [#cc0b020](https://github.com/pydio/cells/commit/cc0b020fd20a5caec4717c1340aabaa6ec1ec021): chore/add prettier config (#697)
- [#3f7298b](https://github.com/pydio/cells/commit/3f7298b26a0387d0257006e6d042822326960e78): tests: configure vitest and react-testing library for unit tests (#696)
- [#58f4ee6](https://github.com/pydio/cells/commit/58f4ee6824864518444b4453d0bf04570d42fed8): fix(meta): improve usermeta propagation in events, along with contexts in jobs Matchers.
- [#7770aa4](https://github.com/pydio/cells/commit/7770aa45a1e9b91fcfacf188b559aa71a70d0806): fix: WPB-19703 sanitize url for links (#695)
- [#5ee1965](https://github.com/pydio/cells/commit/5ee1965180a87ed331cdc3bf4ce4eacfe5cdf244): fix(deps): WPB-19702 - update pdfjs and react-pdf to latest version (#694)
- [#b63feae](https://github.com/pydio/cells/commit/b63feae7a41c281f48a12e6576e3f9256b27407f): fix(errors): fix error catching
- [#693](https://github.com/pydio/cells/pull/693): fix(ckeditor): Improving content sanitization for ckeditor
- [#94645b3](https://github.com/pydio/cells/commit/94645b3ef7e96265974bebc9214dd98561ab9b71): chore(ckeditor): prod build
- [#26ba8fc](https://github.com/pydio/cells/commit/26ba8fc480d1d0459df65319754a198c6870cd0f): fix(ckeditor): - sanitize file content for ckeditor - handle on change with xss sanitize
- [#ef735f4](https://github.com/pydio/cells/commit/ef735f4dbb4d63224ae52293f81aad4781044230): fix(tests): try to fix flakky test
- [#ff23e1a](https://github.com/pydio/cells/commit/ff23e1aaaeddd21947417eed349e5441fa89141e): feat(i18n): more messages + DE/FR translations
- [#691](https://github.com/pydio/cells/pull/691): test(search): Improve unit test coverage and fix failed run
- [#ddcf92f](https://github.com/pydio/cells/commit/ddcf92fe5e37f09bfe220cb95bededfd67c33972): test(search): - fix unit test - add a no match test for AND regex with multi tags
- [#add7563](https://github.com/pydio/cells/commit/add75633dbebee89c8f4a996a21c20f9257b6a0e): fix(mongo): re-introduce unique indexes now that upgrade is fixed
- [#844642f](https://github.com/pydio/cells/commit/844642fcced1f5652dc2a8364d1bcc2b1e2a4669): fix(errors): fix some errors context for proper logger display
- [#9da7f0b](https://github.com/pydio/cells/commit/9da7f0b55c64130cfece9bd07ad5dbd3ab3cb5cf): fix(storage): missing LastUpdated value on workspaces
- [#515c40b](https://github.com/pydio/cells/commit/515c40b0846b5999b5acbe94a25fead61cfde731): feat(datasource): totally remove gcs option, and hide structured datasources behind an advanced configuration.
- [#e23b9bf](https://github.com/pydio/cells/commit/e23b9bfe0ad5c96ec9a2891d3b0ae0da88482f8e): fix(mod): fix missing dep in gomod
- [#ad896a1](https://github.com/pydio/cells/commit/ad896a16324e67c4e440bf55e982b93500612f75): fix(docs-gen): mkdocs compat
- [#e6fa7b4](https://github.com/pydio/cells/commit/e6fa7b4cf2f1d5544733d1f9bef172d2e14d4fc9): fix(apiv2): adding cache-control in the PresignedGet
- [#687](https://github.com/pydio/cells/pull/687): fix(search): ensure tags filter match exactly one result
- [#685](https://github.com/pydio/cells/pull/685): chore(dev-ex): adding live reload config
- [#07f8b26](https://github.com/pydio/cells/commit/07f8b2644b15f2d83ceb140f38ea91e88199c037): chore(search): add description and examples
- [#59f9b59](https://github.com/pydio/cells/commit/59f9b5909e44ee7e64e11706a14d8c4bb7e62194): chore(search): update descriptions
- [#b03ab09](https://github.com/pydio/cells/commit/b03ab09e432efb2493dba14b59d577919b70e751): test(search): adding more test cases
- [#9212a75](https://github.com/pydio/cells/commit/9212a751337966d74b53dea1d78fbb358c18b3bb): fix(search): revert logic to combine using AND and simplify
- [#8d02645](https://github.com/pydio/cells/commit/8d026454e98d81ac2139c9a9a908267ac03720ad): fix(search): prettier
- [#63a4e34](https://github.com/pydio/cells/commit/63a4e34142e588870da2df92f09a853b2c4d2c0d): fix(search): handle various tag search combinations
- [#95a27da](https://github.com/pydio/cells/commit/95a27da6afbc3c188d064217cc94073289281108): fix(search): get exact matches for multiword tags
- [#9d8e8fd](https://github.com/pydio/cells/commit/9d8e8fdb96784cf55598f247ae3f86ea3e14bc1e): fix(search): ensure tags match exactly one result
- [#a9bc219](https://github.com/pydio/cells/commit/a9bc21941f26e0915f02f52e5c08c8341d9c3b7a): chore(dev-ex) add instructions for live reload in readme
- [#dbdc1c4](https://github.com/pydio/cells/commit/dbdc1c42a4f2fd207205caa6fc1ed71bd1931911): chore(dev-ex): specifiy folders to watch
- [#59e5188](https://github.com/pydio/cells/commit/59e5188fed4ce408b5b45da3823042ccccb869ae): chore(dev-ex): adding live reload config
