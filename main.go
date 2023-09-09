package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		vpc, err := ec2.NewDefaultVpc(ctx, "default", &ec2.DefaultVpcArgs{
			Tags: pulumi.StringMap{
				"Name": pulumi.String("Default VPC"),
			},
		})
		if err != nil {
			return err
		}

		// Create an AWS security group that allows HTTP and SSH traffic

		sg, err := ec2.NewSecurityGroup(ctx, "web-secgrp", &ec2.SecurityGroupArgs{
			VpcId: vpc.ID(), // Assign default VPC here

			Description: pulumi.String("Enable HTTP and SSH access"),
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(80),
					ToPort:     pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"), // -1 stands for "all"
					FromPort:   pulumi.Int(0),       // Use 0 to allow all ports
					ToPort:     pulumi.Int(0),       // Max port number
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Create an AWS EC2 instance
		_, err = ec2.NewInstance(ctx, "web-server-www", &ec2.InstanceArgs{
			InstanceType:   pulumi.String("t2.micro"),
			SecurityGroups: pulumi.StringArray{sg.Name},
			Ami:            pulumi.String("ami-0261755bbcb8c4a84"),
			UserData: pulumi.String(`#!/bin/bash

			# Update the system
			sudo apt update
			echo "System Updated"
			
			# Install prerequisite packages
			sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
			echo "--- Prerequisite packages installed"
			
			# Add Docker GPG key
			curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
			echo "Added Docker GPG key"
			
			# Add Docker repository
			sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
			echo "Added Docker repository"
			
			# Install Docker CE
			sudo apt update
			sudo apt install -y docker-ce
			echo "Docker CE installed"
			
			# Install Docker Compose
			sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
			sudo chmod +x /usr/local/bin/docker-compose
			echo "Docker Compose installed"
			
			# Clone the application repository
			git clone https://github.com/shubhamdixit863/pulumi /home/ubuntu/app
			echo "Cloned the application repository"
			
			# Navigate to the app directory and start Docker Compose
			cd /home/ubuntu/app || { echo "Failed to navigate to app directory"; exit 1; }
			sudo docker-compose up -d
			echo "App started"
			`),
			KeyName: pulumi.String("pulumi"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("Pulumi setup"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
