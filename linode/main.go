package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	flag.Parse()

	switch flag.Arg(0) {
	case "Bringup":
		dcs, err := GetDCs()
		if err != nil {
			log.Fatal(err)
		}

		for _, dc := range dcs {
			log.Printf("making linode in %v", dc.Abbr)
			id, err := CreateLinode(dc)
			if err != nil {
				log.Print("Error making linode", err)
				continue
			}
			log.Printf("Made linode %d in %v ", id, dc.Abbr)
			break
		}
	case "DeleteAll":
		linodes, err := ListLinodes()
		if err != nil {
			log.Print("error making linodes")
		}
		for _, id := range linodes {
			log.Print("Deleting ", id)
			err := DeleteLinode(id)
			if err != nil {
				log.Print(err)
			}
		}

	case "Plans":
		out, err := MakeLinodeRequest("avail.linodeplans")
		log.Print(string(out), err)

	case "Kernels":
		out, err := MakeLinodeRequest("avail.kernels")
		log.Print(string(out), err)

	case "Distros":
		out, err := MakeLinodeRequest("avail.distributions")
		log.Print(string(out), err)

	case "IPs":
		ips, err := GetAllIPs()
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range ips {
			fmt.Print(i, ",")
		}

	case "Push":
		ips, err := GetAllIPs()
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range ips {
			fmt.Printf("scp ~/gowork/bin/demesure root@%s:/root/demesure\n", i)
			fmt.Printf("ssh root@%s localectl set-locale LANG=en_US.UTF-8\n", i)
			fmt.Printf("ssh root@%s locale-gen\n", i)
			fmt.Printf("ssh root@%s pacman -Sy --noconfirm tmux\n", i)
			fmt.Printf("ssh root@%s tmux new-session -d -s demesure /root/demesure -listen :8080\n", i)
		}

	default:
		log.Printf("No idea what to do for %q", flag.Arg(1))

	}
}
