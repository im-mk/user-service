data "aws_availability_zones" "available" {
}

resource "aws_vpc" "aws-vpc" {
  cidr_block           = "10.10.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name        = "${var.app_name}-vpc"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_subnet" "private" {
  vpc_id            = aws_vpc.aws-vpc.id
  count             = var.az_count
  cidr_block        = cidrsubnet(aws_vpc.aws-vpc.cidr_block, 8, count.index)
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name        = "${var.app_name}-private-subnet-${count.index + 1}"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.aws-vpc.id
  count                   = var.az_count
  cidr_block              = cidrsubnet(aws_vpc.aws-vpc.cidr_block, 8, var.az_count + count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name        = "${var.app_name}-public-subnet-${count.index + 1}"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_internet_gateway" "aws-igw" {
  vpc_id = aws_vpc.aws-vpc.id
  tags = {
    Name        = "${var.app_name}-igw"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_route" "public" {
  route_table_id         = aws_vpc.aws-vpc.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.aws-igw.id
}

resource "aws_eip" "gw" {
  count      = var.az_count
  depends_on = [aws_internet_gateway.aws-igw]

  tags = {
    Name        = "${var.app_name}-eip-${count.index}"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_nat_gateway" "gw" {
  count         = var.az_count
  subnet_id     = element(aws_subnet.public.*.id, count.index)
  allocation_id = element(aws_eip.gw.*.id, count.index)

  tags = {
    Name        = "${var.app_name}-nat-gw"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_route_table" "private" {
  count  = var.az_count
  vpc_id = aws_vpc.aws-vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = element(aws_nat_gateway.gw.*.id, count.index)
  }

  tags = {
    Name        = "${var.app_name}-rt-${count.index}"
    Application = var.app_name
    Project     = var.project
    Environment = var.app_environment
  }
}

resource "aws_route_table_association" "private" {
  count          = var.az_count
  subnet_id      = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}
