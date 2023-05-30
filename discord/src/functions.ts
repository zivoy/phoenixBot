import {ActionRowBuilder, ButtonBuilder, ButtonStyle, Client, EmbedBuilder, GuildMember, User} from "discord.js";
import {VerifyButtonID} from "./consts"

export function getUser(client: Client, id: string): Promise<GuildMember>
export function getUser(client: Client, name: string, deliminator: string): Promise<GuildMember>
export function getUser(client: Client, nameOrId: string, deliminator: string | undefined = undefined): Promise<GuildMember> {
    return new Promise<GuildMember>((resolve, reject) => {
        client.guilds.cache.forEach((guild) => {
            // check cache
            guild.members.cache.forEach(value => {
                if ((deliminator === undefined && value.user.id === nameOrId) ||
                    (value.user.username == nameOrId && value.user.discriminator == deliminator)) {
                    resolve(value);
                    return;
                }
            });
            // fetch list
            guild.members.list().then(users => {
                users.forEach(value => {
                    if ((deliminator === undefined && value.user.id === nameOrId) ||
                        (value.user.username == nameOrId && value.user.discriminator == deliminator)) {
                        resolve(value);
                        return;
                    }
                });
                reject("not found");
            });
        })
    })
}

export function SendVerification(user: User | GuildMember, code: string) {
    const button = new ButtonBuilder()
        .setStyle(ButtonStyle.Primary)
        .setLabel("Done")
        .setCustomId(VerifyButtonID)
        .setEmoji("✔")
    const buttonRow = new ActionRowBuilder<ButtonBuilder>()
        .setComponents(button);

    const embed = new EmbedBuilder()
        .setColor(0xff00ff)
        .setDescription(`Please add the ${code} to your [RSI accounts Short Bio](https://robertsspaceindustries.com/account/profile) to verify your account, then click done`);

    return user.send({embeds: [embed], components: [buttonRow]})
}