import {ButtonInteraction} from "discord.js";
import {VerifyButtonID, VerifyList} from "../consts";
import {request} from "http";

export default async function (interaction: ButtonInteraction) {
    switch (interaction.customId) {
        case VerifyButtonID:
            console.log(`${interaction.user.username} finished auth`)
            let u = new URL(process.env.VERIFYURL || "http://localhost:80/reload-rsi")
            let req = request({
                host: u.host,
                port: u.port,
                path: u.pathname,
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                }
            }, (r) => {
                console.log("callback: ", r.statusCode);
            });

            let verified = {verified: interaction.user.id, code: ""}
            if (verified.verified in VerifyList) {
                verified.code = VerifyList[verified.verified]
            }
            req.write(JSON.stringify(verified))
            req.end();
            await interaction.update({components: []})
            break;
    }
}