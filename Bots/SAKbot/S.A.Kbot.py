import discord 
import random
import os
from dotenv import load_dotenv


load_dotenv()

token = os.getenv("discord_token")

intents = discord.Intents.all()
client = discord.Client(intents=intents)


aggro_words = ["eimai gia tin poutsa...", "wRaIo dAmAgE", "ADERFIA MOYY!!!", "nai alla tiwra.gr"]
sad_words = ["Den exei nohma i yparksi xwris ton mpakaliaro...", "xwris ton mpakaliaro periplaniemai askopa se ena axanes sumpan...",\
"afougrazomai tis stigmes poy zisame me ton mpakaliaro kai den tis xortainw.."]

trigger_word = "mpakaliaros"

katathlipsi = False
counter = 0
MSG_LOOP = 4
song = "audio_files/Mad World - Gary Jules.mp3"

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
    
    if katathlipsi:
        if counter == MSG_LOOP:
            counter = 0
            await message.channel.send(random.choice(sad_words))
            await message.add_reaction("<:custom_emoji:778663996534161418")


        if msg.startswith("mpakaliaros"):
            voice_channel = message.author.voice.channel
            if voice_channel is not None:
                voice_client = await voice_channel.connect()

                audio_source = discord.FFmpegPCMAudio(song)

                voice_client.play(audio_source)
                
        return

    
    if message.author.id == 629658802203131904:
        if counter == MSG_LOOP:
            counter = 0
            await message.channel.send(random.choice(aggro_words))
            await message.add_reaction("<:custom_emoji:778663996534161418")

client.run(token)