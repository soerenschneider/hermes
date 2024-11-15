# Changelog

## [1.4.0](https://github.com/soerenschneider/hermes/compare/v1.3.0...v1.4.0) (2024-11-12)


### Features

* auto-generate REST client and server code ([c8af1b1](https://github.com/soerenschneider/hermes/commit/c8af1b1daf00e92c2ce7cc9ee2309dfbf08c63c0))


### Bug Fixes

* check staus code of http responses ([efc8145](https://github.com/soerenschneider/hermes/commit/efc814548e46dca8889cda663063491218639eea))
* require subject as according to openapi spec ([cd122a3](https://github.com/soerenschneider/hermes/commit/cd122a34516d0c46392e716ef2c4975b7aee917b))
* use correct metric ([3c0384d](https://github.com/soerenschneider/hermes/commit/3c0384defd78392cfc5b5c18279c629abd5f2b60))

## [1.3.0](https://github.com/soerenschneider/hermes/compare/v1.2.0...v1.3.0) (2024-11-11)


### Features

* add config for dbs ([f1c66d2](https://github.com/soerenschneider/hermes/commit/f1c66d215b8ba80fc7e8e5fdf03b3fa90a8f544b))
* store retryable messages across restarts ([87a42c6](https://github.com/soerenschneider/hermes/commit/87a42c6c7ac3c97465684d3ee57107b35b041a86))

## [1.2.0](https://github.com/soerenschneider/hermes/compare/v1.1.0...v1.2.0) (2024-11-07)


### Features

* add awtrix support ([9f16a26](https://github.com/soerenschneider/hermes/commit/9f16a2649ce4ddaa4fb5ec98c7ec950a0232a01f))
* allow specifying multiple comma-separated service ids ([c05a29f](https://github.com/soerenschneider/hermes/commit/c05a29f183b2d685c4ecf8da5cea3f98ea83b938))


### Bug Fixes

* add call to build awtrix ([dd5923f](https://github.com/soerenschneider/hermes/commit/dd5923f3dcbea96f9bdaedc7780aff21af979d70))
* **deps:** bump github.com/go-playground/validator/v10 ([e44e0f2](https://github.com/soerenschneider/hermes/commit/e44e0f20d212685d6c71a4a998e298d62e0c88b4))
* **deps:** bump github.com/prometheus/client_golang ([9a3beaf](https://github.com/soerenschneider/hermes/commit/9a3beaf86e44ec5a703edbf7af80ae55b1c770fd))
* **deps:** bump golang from 1.23.1 to 1.23.2 ([740017a](https://github.com/soerenschneider/hermes/commit/740017a6e73932753b42d1cbce327c1a3481fce4))
* don't ignore error ([3fab4cc](https://github.com/soerenschneider/hermes/commit/3fab4cc956e3c9c664cf3562bce5faf8bba814ba))
* update metric after retrieving item to fix off by one metric ([fbacd0a](https://github.com/soerenschneider/hermes/commit/fbacd0afa0f52ee5ef361da34d06c3a968912bfe))

## [1.1.0](https://github.com/soerenschneider/hermes/compare/v1.0.1...v1.1.0) (2024-09-28)


### Features

* add feature to relay notifications to gotify ([df78c38](https://github.com/soerenschneider/hermes/commit/df78c389f7a0ad824a902f8cafe160d4cacb0206))
* initial support for awtrix ([8b8f72a](https://github.com/soerenschneider/hermes/commit/8b8f72a87f39cad1e2860574378bb6eb3e95037f))
* support for rabbitmq ([b108a83](https://github.com/soerenschneider/hermes/commit/b108a8349cd3eba515f10a119355dfb07e272b00))


### Bug Fixes

* close body ([0657c3c](https://github.com/soerenschneider/hermes/commit/0657c3cf4939aaf5e5a6bb4841dd325c386257c1))
* **deps:** bump github.com/cenkalti/backoff/v4 from 4.2.1 to 4.3.0 ([82a6828](https://github.com/soerenschneider/hermes/commit/82a68283ffb6ef82f389f007d4574058a3606336))
* **deps:** bump github.com/emersion/go-smtp from 0.20.2 to 0.21.3 ([6c4b5f6](https://github.com/soerenschneider/hermes/commit/6c4b5f68c2353347aefc6d05b4ca62e108520461))
* **deps:** bump github.com/go-playground/validator/v10 ([f511cc7](https://github.com/soerenschneider/hermes/commit/f511cc79dfdcec541b9f7683c42ecf0fdfe206e3))
* **deps:** bump github.com/go-playground/validator/v10 ([580a87a](https://github.com/soerenschneider/hermes/commit/580a87a70399ae95f464e3a8abdb0fce0c66be46))
* **deps:** bump github.com/nikoksr/notify from 0.41.0 to 1.0.0 ([57dcf67](https://github.com/soerenschneider/hermes/commit/57dcf678abfa2b8206a32f5bdd68e8c263da92f2))
* **deps:** bump github.com/prometheus/client_golang ([c03062b](https://github.com/soerenschneider/hermes/commit/c03062b6e47c4b3a6eca465075c61d087103fffd))
* **deps:** bump github.com/rabbitmq/amqp091-go from 1.9.0 to 1.10.0 ([43b8f8b](https://github.com/soerenschneider/hermes/commit/43b8f8b28e83a74c69e53012505fa76c06fc92d2))
* **deps:** bump github.com/rs/zerolog from 1.31.0 to 1.32.0 ([2feca8c](https://github.com/soerenschneider/hermes/commit/2feca8c5d0b28392df68c0d0f8bef2e603556023))
* **deps:** bump github.com/rs/zerolog from 1.32.0 to 1.33.0 ([385ac21](https://github.com/soerenschneider/hermes/commit/385ac217ad0a233c60f3efb2879775a4073974d2))
* **deps:** bump golang from 1.22.0 to 1.22.1 ([5d9e826](https://github.com/soerenschneider/hermes/commit/5d9e8261cc3b95a9ce2bf0bac34efbabf4c05efa))
* **deps:** bump golang from 1.22.1 to 1.22.5 ([28270cf](https://github.com/soerenschneider/hermes/commit/28270cf82018e71761bb310470f11e025290575d))
* **deps:** bump golang from 1.22.5 to 1.23.1 ([854aa07](https://github.com/soerenschneider/hermes/commit/854aa07076703b1d4957e69d45fa76caf572fd96))
* **deps:** bump golang.org/x/net from 0.21.0 to 0.23.0 ([208fcac](https://github.com/soerenschneider/hermes/commit/208fcacc51bea460dd3ce629bb82087dc3dbe6e3))
* **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 ([67615b7](https://github.com/soerenschneider/hermes/commit/67615b725598f6674e19b69118c72dbfb6f246ce))
* fix validation ([ce559cb](https://github.com/soerenschneider/hermes/commit/ce559cb5676b6d055265ee52cc1f65c6341f26e8))
* fix validation ([fd249e7](https://github.com/soerenschneider/hermes/commit/fd249e7e000786de84f7f37bfd3be1565c9cf299))

## [1.0.1](https://github.com/soerenschneider/hermes/compare/v1.0.0...v1.0.1) (2024-02-13)


### Bug Fixes

* **deps:** bump github.com/emersion/go-smtp from 0.18.1 to 0.20.2 ([2adea42](https://github.com/soerenschneider/hermes/commit/2adea421129d9b6202d5ae5707821f8f78e77001))
* **deps:** bump github.com/go-playground/validator/v10 ([d1beeab](https://github.com/soerenschneider/hermes/commit/d1beeab94073a773a8ca8de0448b8b2bc761d79e))
* **deps:** bump github.com/segmentio/kafka-go from 0.4.44 to 0.4.47 ([b69cb7c](https://github.com/soerenschneider/hermes/commit/b69cb7ccac22a06cb7e02616214f25ac02618964))
* **deps:** bump golang from 1.21.3 to 1.22.0 ([1063d1a](https://github.com/soerenschneider/hermes/commit/1063d1a2e0626571cbbde0d08386009da9285f6d))
* **deps:** bump golang.org/x/crypto from 0.14.0 to 0.17.0 ([c9d9cd1](https://github.com/soerenschneider/hermes/commit/c9d9cd113e431527d342d9309f0631f37ac4866b))

## 1.0.0 (2023-11-07)


### Features

* add new metrics ([09b9b93](https://github.com/soerenschneider/hermes/commit/09b9b935336c079ea630ad20b6d024a329d692e0))


### Bug Fixes

* **deps:** bump github.com/adrianbrad/queue from 1.2.1 to 1.3.0 ([614f6f4](https://github.com/soerenschneider/hermes/commit/614f6f40199477b53a7131adb35f8df830f3df40))
* **deps:** bump github.com/go-playground/validator/v10 ([6608106](https://github.com/soerenschneider/hermes/commit/6608106790f67e131ab9d11aadb63f1fe0690eb0))
* **deps:** bump github.com/segmentio/kafka-go from 0.4.43 to 0.4.44 ([8923798](https://github.com/soerenschneider/hermes/commit/8923798a21596154f23164eefe8619c41c6191ac))
* **deps:** bump golang from 1.21.2 to 1.21.3 ([0f08f3c](https://github.com/soerenschneider/hermes/commit/0f08f3c590ad0ec3f79732d4c9f892e948e92ca8))
* **deps:** bump golang.org/x/net from 0.11.0 to 0.17.0 ([2af5ed2](https://github.com/soerenschneider/hermes/commit/2af5ed26f523d4f7b30ab012509da670ab02fa28))
