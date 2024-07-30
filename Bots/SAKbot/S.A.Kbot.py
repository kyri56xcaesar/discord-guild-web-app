import discord 
import random
import os
import asyncio
from discord.ext import commands

from elevenlabs import save
from elevenlabs.client import AsyncElevenLabs
from dotenv import load_dotenv


load_dotenv()

token = os.getenv("discord_token")
elevenlabs_token = os.getenv("evenlabs_token")

eleven_client = AsyncElevenLabs(
    api_key=elevenlabs_token
)

intents = discord.Intents.all()
intents.message_content = True
client = discord.Client(intents=intents)
bot = commands.Bot(command_prefix='>', intents=intents)


aggro_words = ["eimai gia tin poutsa...", "wRaIo dAmAgE", "ADERFIA MOYY!!!", "nai alla tiwra.gr"]
sad_words = ["Den exei nohma i yparksi xwris ton mpakaliaro...", "xwris ton mpakaliaro periplaniemai askopa se ena axanes sumpan...",\
"afougrazomai tis stigmes poy zisame me ton mpakaliaro kai den tis xortainw.."]


katathlipsi = True
counter = 0
MSG_LOOP = 4
song = "../audio_files/Mad World - Gary Jules.mp3"
# audio = "praise_cody.mp3"

cody_text = ">say Cody Fury, the Prime, where valor meets divine. In Cody Fury's reign, strength and wisdom entwine. Prime Cody Fury, our guiding star in the endless sky. With Cody Fury, the Prime, all dreams soar high."


async def play_audio_in_channel(channel, audio):
    vc = await channel.connect()
    vc.play(discord.FFmpegPCMAudio(executable="C:/ffmpeg/bin/ffmpeg.exe", source=audio))

    while vc.is_playing():
        await asyncio.sleep(.1)
    await vc.disconnect()


@client.event
async def on_ready():
    print("Logged in as")
    print(client.user.name)
    print(client.user.id)
    print("--------")


@client.event
async def on_message(message):

    if message == None:
        return
    if message.author == client.user:      
        return

    global katathlipsi
    global counter
    global MSG_LOOP

    msg = message.content
    counter += 1


    if msg.startswith("sakmode=katathlipsi"):
        katathlipsi=True
        await message.add_reaction("🫡")
    if msg.startswith("sakmode=dog"):
        katathlipsi=False
        await message.channel.send("OK eixame")
        
    if msg.startswith("##praise##"):
        voice_channel = message.author.voice.channel
        if voice_channel is not None:
            audio = await eleven_client.generate(
            text = cody_text,
            voice="Callum",
            model="eleven_multilingual_v2"
        )
    
        out = b''
        async for value in audio:
            out += value
        
    
        save(out, "audio.mp3")
    
            
        await play_audio_in_channel(voice_channel, audio="audio.mp3")
            
            
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
        if counter == MSG_LOOP:
            counter = 0
            await message.channel.send(random.choice(aggro_words))
            await message.add_reaction("<:custom_emoji:778663996534161418")

client.run(token)