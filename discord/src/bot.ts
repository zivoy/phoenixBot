import {Client, GatewayIntentBits} from "discord.js";

const {MessageContent, Guilds, GuildMembers} = GatewayIntentBits
export const client = new Client({intents: [Guilds, MessageContent, GuildMembers]})

