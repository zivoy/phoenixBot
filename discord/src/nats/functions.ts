import {DiscordVerifyRequest, DiscordVerifyResponse, nRPCFunction} from "./types";
import {getUser, SendVerification} from "../functions";
import {GuildMember} from "discord.js";
import {client} from "../bot";
import {VerifyList} from "../consts";

const functions: { [name: string]: nRPCFunction } = {
    "verify-rsi": (data: DiscordVerifyRequest): Promise<DiscordVerifyResponse> => {
        console.log(`verification requested for '${data.discord_id ?? data.discord_name}' with code '${data.code}'`)
        return new Promise(async (resolve, reject) => {
            let user: GuildMember
            if (data.discord_id != undefined)
                try {
                    user = await getUser(client, data.discord_id)
                } catch {
                    reject("user not found")
                    return
                }
            else if (data.discord_name != undefined) {
                let parts = data.discord_name.split("#")
                if (parts.length < 2) {
                    if (!/^[a-z0-9._]{2,32}$/.test(data.discord_name)) {
                        reject("not a valid name")
                        return
                    }

                    parts[1] = "0"
                }
                try {
                    user = await getUser(client, parts[0], parts[1])
                } catch {
                    reject("user not found")
                    return
                }
            } else {
                reject("user identifier needed")
                return
            }

            VerifyList[user.id] = data.code
            resolve({discord_id: user.id})
            await SendVerification(user, data.code)
        })
    }
}

export default functions