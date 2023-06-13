# Phoenix management system
a system of programs and bots that help manage community events

connected using NATS

### [Discord bot](/discord)
written in js
- [ ] manage events
- [x] ~~verify users~~

### [Manger](/manager)
written in go

api end points for 3rd party integrations

### [Teamspeak bot](/teamspeak)
- [ ] create channels for events
- [ ] assign roles for varying command structures to allow for cross channel communication via authority
- [ ] move people based on commands to allow for set up and chatting before the event

## development
run 
make sure that there is a .env file that has the discord token and run to bring up the entire system
```bash
docker compose up
```

for deployment and updating use
```bash
docker compose up -d
```


## todo
- look into [nRPC](https://github.com/nats-rpc/nrpc/tree/master) rather to make the requests more structured