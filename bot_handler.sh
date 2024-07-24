#! /bin/bash


function intro()
{
    printf "Welcome to the bot handler.\n"
    printf "Current discord bots: StoufBot - S.A.K. - StratigosBot - KokeBot\n\n"
    
}

function usage()
{
    printf "This script is responsible for handling some discord bots.\n\n"
    printf "Usage: $0\n\n"
    printf "Handles all bots\n\n"
    
    printf "Usage: $0  [OPTION]\n\n"
    printf "Options:\n\n"
    printf "  -stouf\t  Handle Stouf bot.\n"
    printf "  -koke\t  Handle Koke bot.\n"
	printf "  -strat\t  Handle Stratigos bot.\n"
    printf "  -sak\t\t  Handle S.A.K. bot.\n"
    printf "  -help\t\t This message..\n"
    printf "  -list\t\t  List current available bots.\n"
    
        
}

function bot_handler()
{
    intro

    if [ "$1" = "" ]; then  
        
        python KokeBot.py &
        python S.A.Kbot.py &
        python StoufBot.py &
        python StratigosBot.py

    elif [ "$1" = "-stouf" ]; then
        true
    elif [ "$1" = "-koke" ]; then
        true
    elif [ "$1" = "-strat" ]; then
        python StratigosBot.py
    elif [ "$1" = "-sak" ]; then
        true
    else
        usage
    fi

 
}


bot_handler $1
exit 0






