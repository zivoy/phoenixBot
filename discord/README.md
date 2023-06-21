# Discord bot
### functions
none atm


## NATS endpoints

| ENDPOINT                    | Argument                                      | Returns                            | Description                                          |
|-----------------------------|-----------------------------------------------|------------------------------------|------------------------------------------------------|
| discord.function.verify-rsi | **discord_name:** string<br/>**code:** string | **discord_id**: Discord id of user | requests a verification request to be sent to a user |

## Testing
To test this run the following command in the main folder to bring up only the container for this and nats
```bash
docker compose up discord-bot
```


## todo
- make context command for adding picture to gallery 
