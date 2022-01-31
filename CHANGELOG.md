# Changelog

### [0.5.1](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.5.0...v0.5.1) (2022-01-31)


### Bug Fixes

* use `github.com/goccy/go-json` ([8fe8711](https://www.github.com/brokeyourbike/lets-go-chat/commit/8fe8711f8af09cdc4aeb2dd4b18a0b10866621b2))

## [0.5.0](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.4.0...v0.5.0) (2022-01-17)


### Features

* make dependencies with wire ([0739aaf](https://www.github.com/brokeyourbike/lets-go-chat/commit/0739aafe9c1886d8adf669b4537df170d083d01d))


### Bug Fixes

* adjust to generated code ([8c85acc](https://www.github.com/brokeyourbike/lets-go-chat/commit/8c85acca6b674df25bebeb4b5624facffa598f33))
* run goroutine in main ([66de665](https://www.github.com/brokeyourbike/lets-go-chat/commit/66de665b56559dea9873bcd68d8d899a2ad8a833))
* use code geenration ([e238771](https://www.github.com/brokeyourbike/lets-go-chat/commit/e238771ea5c9198acef9c2cb68bd2e85fd24b6c9))

## [0.4.0](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.3.1...v0.4.0) (2022-01-14)


### Features

* make it work, I guess ([70ff5f6](https://www.github.com/brokeyourbike/lets-go-chat/commit/70ff5f6c75abe6ada835bf41712764fa5cb7c914))
* save messages to database ([625bbce](https://www.github.com/brokeyourbike/lets-go-chat/commit/625bbceaa442666b995d91cd278a514968f357d2))


### Bug Fixes

* handle errors ([6b3d06e](https://www.github.com/brokeyourbike/lets-go-chat/commit/6b3d06e498212652af46283e188435ad68401127))
* include userID ([cde5723](https://www.github.com/brokeyourbike/lets-go-chat/commit/cde57232b7a04b226d98cd80ac932a016170ca25))

### [0.3.1](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.3.0...v0.3.1) (2021-12-13)


### Bug Fixes

* decouple users handle ([8d866bc](https://www.github.com/brokeyourbike/lets-go-chat/commit/8d866bc8165ffda269a452d15500368791f61c88))
* remove mocks from source ([f949b0f](https://www.github.com/brokeyourbike/lets-go-chat/commit/f949b0fa1ccd721a680e115747b218b8d4367ed6))
* remove user from active list if request was not upgraded ([d21eba7](https://www.github.com/brokeyourbike/lets-go-chat/commit/d21eba7083a15b17140e316657033318c88c3636))
* return after error ([f5e4be9](https://www.github.com/brokeyourbike/lets-go-chat/commit/f5e4be90bf88d9ab501590a0e8f01b14e7a00ed5))
* return model after creation ([163e2a9](https://www.github.com/brokeyourbike/lets-go-chat/commit/163e2a9361e02af943832a2c20989d925560a102))
* use table tests for `HandleUserCreate` ([749bebf](https://www.github.com/brokeyourbike/lets-go-chat/commit/749bebfaa8b4b551a10c876027e6de86645aa3d2))

## [0.3.0](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.2.0...v0.3.0) (2021-11-26)


### Features

* add `/v1/user/active` endpoint ([5df4ff4](https://www.github.com/brokeyourbike/lets-go-chat/commit/5df4ff4f2e8f9e2febf1306940e500bff1af1e4e))
* add `ErrorLogger` ([b4f282f](https://www.github.com/brokeyourbike/lets-go-chat/commit/b4f282fddb079bedee78f9f2123e7287761aca94))
* add `token` model ([008c54b](https://www.github.com/brokeyourbike/lets-go-chat/commit/008c54b6e5c182d15740e64df4cef7277373fc24))
* add rate_limiter ([f9ac041](https://www.github.com/brokeyourbike/lets-go-chat/commit/f9ac041f92b79b365de6edad9ac4b157594697f0))
* add request logger ([50b2f2f](https://www.github.com/brokeyourbike/lets-go-chat/commit/50b2f2fae1e00cd15ef0ddf01c08f832ae254763))
* add tokens repo ([991566c](https://www.github.com/brokeyourbike/lets-go-chat/commit/991566c8a6e2e8e588aece05395b6276ec90294b))
* add ws endpoint ([423c5be](https://www.github.com/brokeyourbike/lets-go-chat/commit/423c5be5ecf151e877778509bc6ecaf0ae148297))
* allow users to chat ([af37af1](https://www.github.com/brokeyourbike/lets-go-chat/commit/af37af108947cab0ffb81a7dff0218780840dc03))


### Bug Fixes

* add rate limiter for all routes ([0fc248f](https://www.github.com/brokeyourbike/lets-go-chat/commit/0fc248f60906c5a3fb667db8871d2798ac73eff2))
* cast respBody to string ([c955329](https://www.github.com/brokeyourbike/lets-go-chat/commit/c9553290c4284470d3a724f8d67ee46d8294462f))
* do not allow expired tokens ([f666191](https://www.github.com/brokeyourbike/lets-go-chat/commit/f666191b1072ab53115f86f8e405df6be1224c53))
* get method ([93a731a](https://www.github.com/brokeyourbike/lets-go-chat/commit/93a731a3baf5164ab993c0a717c4557301ec6779))
* pass pointer ([da3d6c7](https://www.github.com/brokeyourbike/lets-go-chat/commit/da3d6c7183bf147964c1c0040e2235077df8a8a5))
* remove redundant code ([5e45d93](https://www.github.com/brokeyourbike/lets-go-chat/commit/5e45d931660d4a03b9ba91dd3648aa0738fa34c0))
* simplify recoverer ([9d1caec](https://www.github.com/brokeyourbike/lets-go-chat/commit/9d1caec6f9f9d4209d16ddc0003e3652cf682bb3))
* tidy ([3abc7c9](https://www.github.com/brokeyourbike/lets-go-chat/commit/3abc7c96f4350957876fdb3c4ea5c7fea4bfacba))
* use better method ([ae0a79f](https://www.github.com/brokeyourbike/lets-go-chat/commit/ae0a79f8649b08a2dc357adf9d0b243e642fbf17))
* use in house recoverer ([fe3378f](https://www.github.com/brokeyourbike/lets-go-chat/commit/fe3378f5d32bba3d1246eb8e588fb9bd7a7a2d1f))
* use logger and recoverer ([4c25b9e](https://www.github.com/brokeyourbike/lets-go-chat/commit/4c25b9e73728a4ee13ec885ff38e044a6ff322c4))
* use new specs ([30037fd](https://www.github.com/brokeyourbike/lets-go-chat/commit/30037fd42cc36bd3fd1e6c35e50967578b9fbf77))
* use proper error ([995c543](https://www.github.com/brokeyourbike/lets-go-chat/commit/995c543a5ddd439d0f46ac1f98729b5e407925a1))

## [0.2.0](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.1.0...v0.2.0) (2021-11-11)


### Features

* use database ([e96ca09](https://www.github.com/brokeyourbike/lets-go-chat/commit/e96ca09bcae85c4a52810545cf137e897696eb59))


### Bug Fixes

* return appropriate error message ([6e05a89](https://www.github.com/brokeyourbike/lets-go-chat/commit/6e05a89fce53f8e48f2ac4d2b909f9a33d47387d))
* simplify method ([7609c39](https://www.github.com/brokeyourbike/lets-go-chat/commit/7609c395e49fbea9f341a643efde807272a0d01d))

## [0.1.0](https://www.github.com/brokeyourbike/lets-go-chat/compare/v0.0.1...v0.1.0) (2021-11-08)


### Features

* add `hasher` ([68fcb88](https://www.github.com/brokeyourbike/lets-go-chat/commit/68fcb884d56ce11592ad646407052f08a4fb4893))
* it can respond to `/v1/user` ([8c71bb7](https://www.github.com/brokeyourbike/lets-go-chat/commit/8c71bb7adf746fec0c4f1a5a96f958de21669366))
* users can login ([17ec511](https://www.github.com/brokeyourbike/lets-go-chat/commit/17ec5119e0ed7f6ecf152faeb3a61965aa589af6))


### Bug Fixes

* allow to specify host ([08817f5](https://www.github.com/brokeyourbike/lets-go-chat/commit/08817f55c8d85bc2e33c5b9023ac425ea292c832))
* follow proposed structure ([12fe5aa](https://www.github.com/brokeyourbike/lets-go-chat/commit/12fe5aaec3cc7e6d500b85c5140a0ac2408ed385))
* merge require ([3a96ee2](https://www.github.com/brokeyourbike/lets-go-chat/commit/3a96ee2c7969b9d31ec3537272fcb3ca6dd54a6c))
* simplify project structure ([96cf960](https://www.github.com/brokeyourbike/lets-go-chat/commit/96cf9604433b3975834d6c79eff8da52f4de5d4f))
* test `user` ([347e608](https://www.github.com/brokeyourbike/lets-go-chat/commit/347e60804cb307e49a0d523a75d151813e44173f))
* tidy the dependencies ([bc1a5cd](https://www.github.com/brokeyourbike/lets-go-chat/commit/bc1a5cdab4cdd5177af79a7300c0f43d9be2b0d9))
* use better method ([d3f5759](https://www.github.com/brokeyourbike/lets-go-chat/commit/d3f5759c4e285a15e03379b2cc4556b7ca12d3a5))
* use port from env ([5ead1fa](https://www.github.com/brokeyourbike/lets-go-chat/commit/5ead1fa56d42b95e0a6a9ebbf51a6eedd894d60b))
* use uuid as key ([3846e1d](https://www.github.com/brokeyourbike/lets-go-chat/commit/3846e1d1ea13f3a643f7329c02ccff18b7560e55))
