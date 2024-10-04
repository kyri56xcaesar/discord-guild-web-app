import discord
import random
import asyncio

from discord.ext import commands
from elevenlabs import save
from elevenlabs.client import AsyncElevenLabs
from dotenv import dotenv_values

config = dotenv_values("../../bots.env")

token = config["sak_bot_token"]
elevenlabs_token = config["elevenlabs_token"]

eleven_client = AsyncElevenLabs(
    api_key=elevenlabs_token
)

intents = discord.Intents.all()
intents.message_content = True
bot = commands.Bot(command_prefix='>', intents=intents)

aggro_words = ["eimai gia tin poutsa...", "wRaIo dAmAgE", "ADERFIA MOYY!!!", "nai alla tiwra.gr"]
sad_words = ["Den exei nohma i yparksi xwris ton mpakaliaro...", "xwris ton mpakaliaro periplaniemai askopa se ena axanes sumpan...",
             "afougrazomai tis stigmes poy zisame me ton mpakaliaro kai den tis xortainw.."]

katathlipsi = False
counter = 0
MSG_LOOP = 3
song = "../audio_files/Mad World - Gary Jules.mp3"

cody_text = ">say Cody Fury, the Prime, where valor meets divine. In Cody Fury's reign, strength and wisdom entwine. Prime Cody Fury, our guiding star in the endless sky. With Cody Fury, the Prime, all dreams soar high."
voices = {"Lina":"rCog6MJ305VojjZbtGWQ", "Callum":"N2lVS1w4EtoT3dr4eOWO", "Koulis":"6z1Ks05MOtac6wYNh9PJ"}
myvoice = voices["Callum"]


async def play_audio_in_channel(channel, audio):
    vc = await channel.connect()
    vc.play(discord.FFmpegPCMAudio(executable="C:/ffmpeg/bin/ffmpeg.exe", source=audio))

    while vc.is_playing():
        await asyncio.sleep(0.1)
    await vc.disconnect()



@bot.event
async def on_ready():
    print("Logged in as")
    print(bot.user.name)
    print(bot.user.id)
    print("--------")

@bot.event
async def on_message(message):
    if message.author == bot.user:
        return

    global katathlipsi
    global counter
    global MSG_LOOP

    msg = message.content
    counter += 1

    if msg.startswith("sakmode=katathlipsi"):
        katathlipsi = True
        await message.add_reaction("🫡")
    if msg.startswith("sakmode=dog"):
        katathlipsi = False
        await message.channel.send("OK eixame")
        
    if msg.startswith(">say"):
        
        text_split = msg.split(">say")[-1]
        text_voice = text_split.split("]", 1)
        text = text_voice[-1]
        voice = text_voice[0].replace("[", "")
        myvoice = voices.get(voice)
        
        if myvoice == "":
            myvoice = "Callum"
        
        print(f"text:{text}, voice:{voice}, myvoice:{myvoice}\n")

        
        voice_channel = message.author.voice.channel
        if voice_channel is None:
            await voice_channel.send("mpes se vc re kagoura")
            return

        #await ctx.message.delete()


        audio = await eleven_client.generate(
            text=text,
            voice=myvoice,
            model="eleven_multilingual_v2"
        )

        out = b''
        async for value in audio:
            out += value

        save(out, "aud.mp3")

        await play_audio_in_channel(voice_channel, "aud.mp3")

    if msg.startswith("##praise##"):
        voice_channel = message.author.voice.channel
        if voice_channel is not None:
            audio = await eleven_client.generate(
                text=cody_text,
                voice=myvoice,
                model="eleven_multilingual_v2"
            )

            out = b''
            async for value in audio:
                out += value

            save(out, "audio.mp3")
            await play_audio_in_channel(voice_channel, "audio.mp3")

    if katathlipsi:
        if counter == MSG_LOOP:
            counter = 0
            await message.channel.send(random.choice(sad_words))
            await message.add_reaction("<:custom_emoji:778663996534161418")

        if msg.startswith("mpakaliaros"):
            voice_channel = message.author.voice.channel
            if voice_channel is not None:
                await play_audio_in_channel(voice_channel, song)
        return

    if message.author.id == 629658802203131904:
        if counter % MSG_LOOP == 0:
            counter = 0
            await message.channel.send(random.choice(aggro_words))
            await message.add_reaction("<:custom_emoji:778663996534161418")

    #await bot.process_commands(message)

bot.run(token)
