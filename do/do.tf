variable "do_token" {}

provider "digitalocean" {
    token = "${var.do_token}"
}

resource "digitalocean_ssh_key" "rene" {
    name = "rene"
    public_key = "${file("/home/colin/.ssh/id_ed25519.pub")}"
}

variable "zones" {
    type = "list"
    default = [
        "nyc1", "nyc2", "nyc3",
        "sfo1", "sfo2",
        "ams1", "ams2",
        "sgp1",
        "lon1",
        "fra1",
        "tor1",
        "blr1",
    ]
}

resource "digitalocean_droplet" "demesure" {
    image  = "ubuntu-16-04-x64"
    name   = "${element(var.zones, count.index)}"
    region = "${element(var.zones, count.index)}"
    size   = "512mb"
    ssh_keys = ["${digitalocean_ssh_key.rene.id}"]
    connection{
        type = "ssh"
        user = "root"
    }
    provisioner "file" {
        source = "~/gowork/bin/demesure"
        destination = "/root/demesure"
    }
    provisioner "remote-exec" {
        inline = [
            "chmod +x /root/demesure",
            "apt-get install tmux -y",
            "tmux new-session -d -s demesure /root/demesure -listen :8080",
        ]
    }
    count = "${length(var.zones)}"
}

output "ips" {
    value = "${digitalocean_droplet.demesure.*.ipv4_address}"
}
