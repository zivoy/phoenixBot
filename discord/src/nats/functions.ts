import {DiscordVerifyRequest, DiscordVerifyResponse, nRPCFunction} from "./types";
import {getUser, SendVerification} from "../functions";
import {GuildMember} from "discord.js";
import {client} from "../bot";

const functions: { [name: string]: nRPCFunction } = {
    "verify-rsi": (data: DiscordVerifyRequest): Promise<DiscordVerifyResponse> => {
        return new Promise(async (resolve, reject) => {
            let user: GuildMember
            if (data.discord_id != undefined)
                try {
                    user = await getUser(client, data.discord_id)
                } catch {
                    reject("user not found")
                    return
                }
            else {
                let parts = data.discord_name!.split("#")
                if (parts.length < 2) {
                    reject("not a valid name")
                    return
                }
                try {
                    user = await getUser(client, parts[0], parts[1])
                } catch {
                    reject("user not found")
                    return
                }
            }

            resolve({discord_id: user.id})
            SendVerification(user, data.code)
        })
    }
}

export default functions