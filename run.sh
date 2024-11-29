#!/usr/bin/env bash

BUGNF=false
SCP=false

while [[ $1 ]]
do
	case $1 in
		--buggy ) BUGNF=true; shift ;;
		--with-scp ) SCP=true; shift ;;
		*)
			echo "Usage: $0 [--buggy] [--with-scp]"
			exit 1
	esac
done

function terminate()
{
    echo "Terminating... bringing down Docker containers"
	if [[ $BUGNF = true ]]; then
		docker compose -f docker-compose-buggy.yaml down
	elif [[ $SCP = true ]]; then
		docker compose -f docker-compose-scp.yaml down
	else
		docker compose down
	fi
    exit 0
}
trap terminate SIGINT

if [[ $BUGNF = true ]]; then
	docker compose -f docker-compose-buggy.yaml up -d
elif [[ $SCP = true ]]; then
	docker compose -f docker-compose-scp.yaml up -d
else
	docker compose up -d
fi

if [[ $BUGNF == true || $SCP == true ]]; then
	tmux new-session "sleep 2 && docker logs -f upf" \
		\; splitw -p 75 "sleep 2 && docker logs -f amf" \
		\; splitw -p 33 "sleep 2 && docker logs -f udm" \
		\; splitw -t 0 -h "sleep 2 && docker logs -f ausf" \
		\; splitw -t 3 -h "sleep 3 && docker logs -f scp" \
		\; swap-pane -t 4 -s 2
else
	tmux new-session "sleep 2 && docker logs -f upf" \
		\; splitw -d "sleep 2 && docker logs -f amf" \
		\; splitw -h "sleep 2 && docker logs -f udm" \
		\; splitw -t 2 -h "sleep 2 && docker logs -f ausf"
fi

terminate