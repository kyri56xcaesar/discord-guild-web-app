import database_handler as dh
import os

from dotenv import load_dotenv

load_dotenv()

koke_trigger_words = ["Koke", "koke", "KOKE", "ela", "0/100", "3 cs", "ELA", "OUELA", "xontre", "XONTRE", "gamiesai", "gamiese", "GAMIESAI", "GAMIESE",
                 "elaaaa", "elaa", "elaaa", "o palios", "O PALIOS", "pouleriko", "lb", "Ela", "Xontre", "Gamiesai", "Pouleriko", "poyleriko", "Poyleriko"]
koke_lines = ["TIIIIIIIII", "geia sou mikro pouleriko", "modia gamiesai", "geia sou dhmhtrh", "gamo to soi moy", "ksypnate reeeeeiii", "ksafnika lupothimisa", "epesa pano sto pomolo", "exw faei camp", "arneitai na paiksei",
              "DES TI KANEI O JUNGLER MOY!", "noliiii, pame kana duo??", "den exw aderfia", "modia. I POUTANA I manas", "aggele pame duo???", "glistrisa se kati ladia kai anoiksa to pigouni moy bro", "otan moy vazoyn duskola antapokrinomai"]



stouf_trigger_words = ["volibear", "voli", "teo", "teo123stouf", "teo123", "stouf",
                 "VOLI", "Voli", "VOLIBEAR", "Volibear", "!Voli", "!VOLI", "!voli", "arkouda",
                 "@teo", "ARKOUDA", "Arkouda"]
stouf_lines = ["ola gia mia trupa ginontai", "voliiiiii", "bmw me collector",
             "ayrio pao diamon", "simerini efimerida\n exei kialo filo\to simera exo 17 clips 1vs 3 kai 1vs 4"]


vassano_trigger_words = ["yorick", "praktoras"]
vassano_lines = ["Kalhspera kai kali vradia edw me ton aderfo fofota", "gamo ton karioli ton OTE apo tis 7 to prwi mexri twra den eixa internet gamo tinpanagia tous prwi prwi",
         "gamw tis manes tous", "simera anastithike o bines", "me trexoyn oi traktores", "ksekinise to praktorio", "ff x9 matchfixers/griefers/hostagers", "epesa PLAT apo master 230 lp se mia nuxta", "pame vraxous gt me trexoyn ta mpastarda", "koke gamw tinmana sou vlaka",
         "me treksan tapromo mexri diamond 0 lp", "legit ama den paw master to kleinw kai 3 meres naparei ", "na se paw mia volta aptis kuries?", "o mikros apo tin pisw porta", "me trexei o lazer"]

sak_trigger_words = []
sak_sad_lines = ["eimai gia tin poutsa...", "wRaIo dAmAgE", "ADERFIA MOYY!!!", "nai alla tiwra.gr"]
sak_normal_lines = ["Den exei nohma i yparksi xwris ton mpakaliaro...", "xwris ton mpakaliaro periplaniemai askopa se ena axanes sumpan...",
             "afougrazomai tis stigmes poy zisame me ton mpakaliaro kai den tis xortainw.."]


dh.insert_trigger_words_in_table(os.getenv('DATABASE_NAME')+".db")

# ('StoufBOT')     -- id 1
# ('KokeBOT')      -- id 2
# ('SAKBot')       -- id 3
# ('VasanoBOT')    -- id 4
# ('StratigosBOT') -- id 5
# ('HayateBOT')    -- id 6