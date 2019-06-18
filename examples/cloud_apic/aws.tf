
data "aws_vpc" "capic_vpc" {
  tags {
    Name = "context-[${aci_vrf.vrf1.name}]-addr-[${aci_cloud_context_profile.context_profile.primary_cidr}]"
  }
  depends_on = ["aci_cloud_context_profile.context_profile"]
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

data "aws_subnet" "capic_subnet" {
  filter {
    name   = "tag:Name"
    values = ["subnet-[${aci_cloud_subnet.cloud_apic_subnet.ip}]"]
  }

  availability_zone = "us-west-1a"
  vpc_id            = "${data.aws_vpc.capic_vpc.id}"
}

resource "aws_instance" "web" {
  ami                         = "${data.aws_ami.ubuntu.id}"
  instance_type               = "t2.micro"
  subnet_id                   = "${data.aws_subnet.capic_subnet.id}"
  associate_public_ip_address = true

  tags = {
    Name = "admin-ep2"
  }
}
