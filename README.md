## Conkeys

It's a configuration server that save a key typed value pair for system that needed distributed configurations.


## API

Clients can access Conckeys via http rest api.

## TODO

- [ ] Add listening port as Env variable
- [ ] Add DockerFile in order to use containers
- [ ] Add dotenvconfig in order to use better systems
- [ ] Try to use Testcontainers for both dev environment and testing
    - [ ] Memory storage in order to unit testing
- [ ] Logging System (Logrus?)
- [ ] Optimize key search and maybe remove allkeys api
- [ ] Define length of varchar fields in all postegres storage
