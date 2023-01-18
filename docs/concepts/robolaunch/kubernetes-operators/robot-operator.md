---
sidebar_label: Robot Operator
sidebar_position: 1
---
# robolaunch Robot Operator

robolaunch Kubernetes Robot Operator manages lifecycle of ROS/2 based robots and enables defining, deploying and distributing robots declaratively.

## Overview

Like every software, each package developed with Robot Operating System has a lifecycle of build and runtime operations. In the meantime, "a robot" can be defined as a collection of ROS packages and it inherits the same lifecycle of events. The diagram below briefly explains the steps of launching a robot.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-launch-lifecycle.png" width="100%"/>
</div>

When developing the software of a ROS/2 package or robot, the lifecycle of events can be described as below.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-development-lifecycle.png" width="100%"/>
</div>

Kubernetes Robot Operator aims to **atomize these steps and automate the lifecycle of robot deployment, development and simulation**. We believe that adapting robot's assets to Kubernetes leverages the benefits of ROS/2, also makes us able to give a declarative API for any ROS/2 based robot.

Besides robot's lifecycle management, robolaunch aims to give community following benefits with Kubernetes Robot Operator:

- Robot observability with ROS/2 tools or external tools
- GPU offloading in simulation or training cases
- ROS/2 node decomposition and regional distribution
- Providing interface for fleet management and DDS connectivity

## Architecture

### Robot Custom Resource
Robot is a custom resource that manages multiple subresources to operate the ROS/2 software. The diagram below shows the components of a robot custom resource.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-robot-custom-resource.png" width="100%"/>
</div>

### Operational Subresources
This is a current list of operational subresources and their duties.

- **Configuration Agent ([RobotConfig](#configuration-agent-robotconfig))** - prepares an OS environment for robot based on it's image
- **Cloning Agent ([WorkspaceManager](#cloning-agent-workspacemanager))** - clones repositories and prepares ROS/2 workspaces
- **Building Agent ([BuildManager](#building-agent-buildmanager))** - installs packages' dependencies and builds the workspaces
- **Launching Agent ([LaunchManager](#launching-agent-launchmanager))** - launches the packages

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-operational-subresources.png" width="100%"/>
</div>

#### Configuration Agent (RobotConfig)
Configuration agent is responsible for give the robot initial configuration and creating the OS environment according to the specifications of Kubernetes cluster and desired ROS distro. This agent runs one time, when the robot is newly deployed.

Configuration agent,
- **Creates the RBAC agent** - the agent that is responsible for managing roles, role bindings and service accounts for robot
- **Creates volumes of robot** - this volumes hold the main directories that is essential for robot's OS level configuration
- **Creates jobs to configure volumes** - jobs configure these volumes so that the robot can work properly when these volumes are mounted to another container. They also configure the volumes based on the ROS distro desired.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-configuration-agent.png" width="100%"/>
</div>

#### Cloning Agent (WorkspaceManager)
Cloning agent is responsible for cloning the desired repositories into workspaces. This agent runs one time, after configuration agent successfully finished it's job.

Cloning agent,
- **Creates jobs to clone repositories** - these jobs configure the workspaces by cloning any publicly available Git repositories into workspaces.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-cloning-agent.png" width="100%"/>
</div>

#### Building Agent (BuildManager)
Building agent is responsible for making the robot ready for launching. It creates jobs to install dependencies, build ROS/2 packages and make any other configurations. Jobs' operations also can be defined using Robot API customly. This agent runs after cloning agent is successfully finished it's job for the first time, then it can be executed any time by manipulating the Robot API in the ROS/2 software development loop.

Building agent,
- **Creates jobs to build packages** - these jobs installs package dependencies, build ROS/2 packages and configures the OS to make robot ready for the launching operation.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-building-agent.png" width="100%"/>
</div>

#### Launching Agent (LaunchManager)
Launching agent is responsible for launching the ROS nodes. By ROS node decomposition feature, every node can be started in seperate container. This agent runs after building agent finished it's job successfully for the first time, then it can be executed any time by manipulating the Robot API in the ROS/2 software development loop.

Launching agent,
- **Creates launch pod** - this pod contains the containers that are each resposible for executing a ROS node.

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/concepts/robot-operator-launching-agent.png" width="100%"/>
</div>

### Other Subresources
- **Tools Agent ([RobotTools]())** - manages the internal or external tools for robot such as ROS bridge or cloud IDE
- **RBAC Agent ([RobotRBAC]())** - gives the configuration agent the roles, role bindings and service accounts for tools or users
- **Data Agent ([RobotData]())** - exports the robot-related data to Kubernetes, such as ROS/2 nodes, topics, services and actions
- **Artifact Agent ([RobotArtifact]())** - holds the initial configurations of robot and can be referenced when creating a robot

## API Reference

- [One Page Robot API Reference](./custom-resources/robot.md#reference)

## Software Stack
Kubernetes Robot Operator is being developed with [Kubebuilder SDK](https://book.kubebuilder.io/) (current version is [v3.2.0](https://github.com/kubernetes-sigs/kubebuilder/releases/tag/v3.2.0)), in Go (current version is v1.16).
