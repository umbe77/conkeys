## Conkeys

It's a configuration server that save a key typed value pair for system that needed distributed configurations.


## API

Clients can access Conckeys via http rest api.

## TODO

- ~~[ ] Embed etcd in order to create a distributed cluster?~~
- [ ] Logging System (Logrus?)
- [ ] Optimize key search and maybe remove allkeys api
- [ ] Define length of varchar fields in all postegres storage
- [ ] Memory storage in order to unit testing
- [ ] Should we move to tcp communication?
