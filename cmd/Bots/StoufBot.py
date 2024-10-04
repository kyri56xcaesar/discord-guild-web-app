import discord
import random
import requests
import json
import bot_utils as ut

from dotenv import dotenv_values



config = dotenv_values("../../bots.env")
responding = True
token = config["stouf_bot_token"]
id = config["stouf_user_id"]
tracker = list()

print(f"token: {token}, id: {id}")


intents = discord.Intents.all()
client = discord.Client(intents=intents)


trigger_words = ["volibear", "voli", "teo", "teo123stouf", "teo123", "stouf",
                 "VOLI", "Voli", "VOLIBEAR", "Volibear", "!Voli", "!VOLI", "!voli", "arkouda",
                 "@teo", "ARKOUDA", "Arkouda"]

pop_words = ["ola gia mia trupa ginontai", "voliiiiii", "bmw me collector",
             "ayrio pao diamon", "simerini efimerida\n exei kialo filo\to simera exo 17 clips 1vs 3 kai 1vs 4"]





def get_quote():
    response = requests.get("https://zenquotes.io/api/random")
    json_data = json.loads(response.text)
    quote = json_data[0]['q'] + " - Theodoros Stoufios"

    return (quote)


def update_pops(new_message):
    global pop_words

    pop_words.append(new_message)


def delete_pop(index):
    global pop_words

    if len(pop_words) > index:
        del pop_words[index - 1]


def print_pops():
    global pop_words

    message = ""

    for i, pop in enumerate(pop_words):
        message += f"{i+1}. {pop}.\n"

    #print(message)

    return message


def HELP():

    msg = "  $inspire: mathe ta vasika\n\
  $new: kerase neo leksilogio\n\
  $del [number]: svise tin malakia p evales\n\
  $list: deikse ta swthika moy\n\
  $voulwne: fimotro\n\
  $ksypna: peripolia\n\
  $vohtheia: auti i malakia\n\
  "

    return msg


@client.event
async def on_ready():
    print("We have logged in as {0.user}".format(client))


@client.event
async def on_message(message):

    if message == None:
        return

    #print("Received message event")
    #print("Message content: " + message.content)

    global trigger_words
    global pop_words
    global responding

    if message.author == client.user:
        return

    if message.author.id == id:
        tracker.append(message.content)

    msg = message.content

    if msg.startswith("$inspire"):
        quote = get_quote()
        await message.channel.send(quote)

    if responding:
        if any(word in msg for word in trigger_words):
            await message.channel.send(random.choice(pop_words))

    if msg.startswith("$new"):
        msgx = msg.split("$new ", 1)[1]
        update_pops(msgx)
        await message.channel.send("To dexomai bro.")

    if msg.startswith("$del"):
        if pop_words is not None:
            index = int(msg.split("$del", 1)[1])
            if index.isnumeric() == False:
                return
            if index < 1 or index > len(pop_words):
                await message.channel.send("Duskola eisai bro, ektos oriwn...")
            else:
                delete_pop(index)
                await message.channel.send("Svisa ola")

    if msg.startswith("$list"):  # encouragements = []

        await message.channel.send(print_pops())

    if msg.startswith("$voulwne"):
        if responding == False:
            await message.channel.send("voulwmeni trupa exw afentiko.")
        else:
            responding = False
            await message.channel.send("Ok to voulwnw.")

    if msg.startswith("$ksypna"):
        if responding:
            await message.channel.send("Edw eimai re clown")
        else:
            responding = True
            await message.channel.send("Ksekinaw to treksimo.")

    if msg.startswith("$vohtheia"):
        await message.channel.send(HELP())

client.run(token)
ut.save_data("data/Stouf", tracker)
