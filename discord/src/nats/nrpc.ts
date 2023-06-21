import {connect, JSONCodec, type NatsConnection} from "nats";
import functions from "./functions";
import * as process from "process";

let nc: NatsConnection;
const jc = JSONCodec()

export async function ConnectNatsListener(natsAddress:string) {
    try {
        nc = await connect({
            name: "Discord bot",
            pingInterval: 60*1000, // ping once a minute
            servers: natsAddress
        });
    }
    catch  {
        console.error("error connecting to nats")
        process.exit(1)
    }
    console.log(`Connected to NATS at ${nc.getServer()}`)

    const sub = nc.subscribe("discord.>");
    for await (const m of sub) {
        const func = (m.subject.match(/^discord\.function\.(.+)$/) ?? [m.subject, undefined])[1]
        if (func != undefined && func in functions) {
            try {
                functions[func](jc.decode(m.data)).then(resp=>{
                    m.respond(jc.encode(resp))
                }).catch(e=>{
                    m.respond(jc.encode({
                        error: e
                    }));
                })
            } catch (e) {
                m.respond(jc.encode({
                    error: e.string()
                }));
            }
        }
    }
    console.log("subscription closed");
    await nc.drain();
}