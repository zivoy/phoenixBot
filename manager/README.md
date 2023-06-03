# Manager

### Endpoints

| Method | URL       | Arguments                                     | Returns                                                              | Description                                                                                        |
|:------:|-----------|-----------------------------------------------|----------------------------------------------------------------------|----------------------------------------------------------------------------------------------------|
|  POST  | `/verify` | **discord_name:** string<br/>**code:** string | **success:** bool<br/>**error?:** string<br/>**discord_id?:** string | sends message to the discord bot to have it DM a user given their name to verify their RSI account |

to test this run the following command in the main folder to bring up only the container for this and nats

```bash
docker compose up manager
```