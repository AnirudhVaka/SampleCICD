provider "null" {}

resource "null_resource" "provision_self_hosted_server" {
  provisioner "remote-exec" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y nginx"
    ]

    connection {
      type        = "ssh"
      host        = "103.174.226.90"
      user        = "root"   # You can replace 'root' with your server's username
      password    = "8JL6HzP*KIvf"
      port        = 22
      timeout     = "2m"
    }
  }
}
