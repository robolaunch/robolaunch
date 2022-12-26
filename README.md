<div align="center">
    <a href="robolaunch.io" title="robolaunch">
        <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/logos/svg/rocket.svg" width="400" height="auto"/>
    </a>
    <p align="center">
        <h1 href="robolaunch.io" title="robolaunch">
            robolaunch Cloud Robotics Platform
        </h1>
    </p>
    <p align="center">
        <a href="https://github.com/robolaunch/robolaunch/releases">
        <img src="https://img.shields.io/github/v/release/robolaunch/robolaunch" alt="release">
        </a>
        <a href="https://github.com/robolaunch/robolaunch/blob/master/LICENSE">
        <img src="https://img.shields.io/github/license/robolaunch/robolaunch" alt="license">
        </a>
        <a href="https://github.com/robolaunch/robolaunch/issues">
        <img src="https://img.shields.io/github/issues/robolaunch/robolaunch" alt="issues">
        </a>
        <a href="https://github.com/robolaunch">
        <img src="https://img.shields.io/github/stars/robolaunch?style=social" alt="issues">
        </a>
    </p>
    <p align="center">
        <a href="https://twitter.com/robolaunch">
        <img src="https://img.shields.io/twitter/follow/robolaunch?style=social" alt="issues">
        </a>
    </p>
</div>

# Overview

[robolaunch](robolaunch.io) is a Cloud-Native Robotics Platform that provides the end-to-end infrastructure, software stack and tools for developing, simulating, deploying and operating ROS/ROS2 robots at scale. 

robolaunch's end-to-end infrastructure, which we also call it "Unified Robotics", enables robotics developers to develop & simulate robotics applications on cloud via VDI infrastructure that uses GPU power, deploy robotics applications to the physical robots remotely and operate the robot's status in a run-time with management, orchestration and monitoring functionalities.

# Why robolaunch?

Robotics development contains many challenges itself. Such as it requires multi-domain knowledge to build a robot and this process is mostly dependent on expensive hardware. Therefore the process becomes hard to achieve and time consuming as a result.

With robolaunch, we want to give robotics developers the ability to develop, deploy and manage ROS/2 robots easily. One of our primary motivation in developing robolaunch is to remove barriers to entry in robotics development.

robolaunch Cloud Robotics Platform is based on Kubernetes because we believe future of robotics will be a distributed and cloudified system and Kubernetes helps us to provide tenancy, automated deployments and self configuration of work-processes.

We believe that robolaunch will provide great convenience to robotics developers in many ways. We will continue evolving robolaunch with our dedicated team. In the meantime, we hope you to contribute to the codebase so that we can improve robolaunch together.

- **Faster to Getting Started -** eliminates the time spent configuring a development environment locally on personal computer

- **Cost Efficient -** eliminates the costs required for robotics hardware and computational resources

- **Sim-to-Real -**  remotely deploy the robotics applications that developed on robolaunch to the physical robot(s)

- **Visualization -** monitor and manage your robots in run time from robolaunch dashboard for a full control over your robots

- **Flexibility -** leverage from entire platform features, or pick the components according to your needs

# Feature Set

Four non-functional outputs ([Scalability](#earthamericas-scalability), [Connectivity](#link-connectivity), [Acceleration](#zap-acceleration), [Automation](#gear-automation)) that have deriven from functional features are listed below.

#### :earth_americas: Scalability
- **Infrastructure as Code -** automated Kubernetes cluster provisioning
- **Multicluster Orchestration -** manage clouds and robots geographically
- **Container-based -** benefit from the power of virtualization, optimized for the edge

#### :link: Connectivity
- **Built-in VPN (cloud-powered mode) -** access robots securely regardless of their network environment
- **Isolated L2 Networks -** isolating tenants' data and control planes
- **Isolated Robots -** isolating robots by both ROS namespacing and DDS domains


#### :zap: Acceleration

- **Hardware Acceleration -** use hardware acceleration (GPU) for robot simulation, training and virtual desktop environment
- **Virtual Desktop Environment -** use desktop environment to run observability tools
- **Hardware-in-the-Loop Simulation -** test, train and simulate robotics applications by updating embedded software remotely

#### :gear: Automation

- **Robot-as-a-Service -** deploying and managing robots and fleets declaratively
- **On-the-Fly Development -**  software development using integrated Cloud IDE on runtime
- **Easy Updates -**  remotely update the software on your robots 
- **Robot Templates -** predefined and reusable robot templates to work with right away
- **Extended Kubernetes API -** extend robolaunch feature set according to your specific needs

Learn more on [robolaunch.io](robolaunch.io).

<div align="center">
    <h2>
    &#11088; Start contributing or support the project with a star! &#11088;
    </h2>
</div>

# Core Concepts

Architectural components of robolaunch are defined below.

- **[robolaunch Management Suite](./docs/concepts/robolaunch/management-suite/README.md)** - contains management console and components
    - **[robolaunch Management Console](./docs/concepts/robolaunch/management-suite/management-console/README.md)** - subcomponent of management suite that contains orchestrator and UI
        - **[robolaunch Orchestrator](./docs/concepts/robolaunch/management-suite/management-console/orchestrator.md)** - backend that orchestrates platform operations using workflow engine
        - **[robolaunch UI](./docs/concepts/robolaunch/management-suite/management-console/ui.md)** - frontend that enables performing platform capabilities via user interface
  - **[robolaunch Management Components](./docs/concepts/robolaunch/management-suite/management-components/README.md)** - subcomponent of management suite that contains other open-source management tools

- **[robolaunch Kubernetes Infrastructure](./docs/concepts/robolaunch/kubernetes-infrastructure/README.md)** - Kubernetes infrastructures for central cloud, regional cloud and edge cloud
    - **[robolaunch Super Cloud Instance](./docs/concepts/robolaunch/kubernetes-infrastructure/regional/super-cloud-instance.md)** - Kubernetes cluster which is the control plane of cloud instances (virtual clusters)
    - **[robolaunch Cloud Instance (Virtual Cluster)](./docs/concepts/robolaunch/kubernetes-infrastructure/regional/cloud-instance.md)** - tenant Kubernetes cluster that contains fleets of robots in a region

- **[robolaunch Kubernetes Operators](./docs/concepts/robolaunch/kubernetes-operators/README.md)**
    - **[robolaunch Robot Operator](./docs/concepts/robolaunch/kubernetes-operators/robot-operator.md)** - software that does decomposition, regional distribution, lifecycle management and configuration of ROS-based robots in cloud-only, cloud-powered or cloud-connected modes
    - **[robolaunch Robot](./docs/concepts/robolaunch/assets/robot.md)** - contains ROS components (Runtime Pod, VDI, Code Server IDE, ROS Tracker, Foxglove Studio, ROS Bridge Suite, Configurational Resources), robolaunch Robot instances can be decomposed and distributed to both cloud instances and physical robots using federation
    - **[robolaunch Fleet Operator](./docs/concepts/robolaunch/kubernetes-operators/fleet-operator.md)** - software that manages lifecycle and configuration of multiple robots and robot's connectivity layer that contains DDS Discovery Server and ROS Bridge Suite
    - **[robolaunch Fleet](./docs/concepts/robolaunch/assets/fleet.md)** - contains multiple robot deployments across multiple physical robots and cloud instances using federation

#### Glossary
- **[Fleet Namespace](./docs/concepts/glossary.md)** - corresponds to a federated Kubernetes namespace in a cloud instance which is abstractly matched with a fleet
- **[Robot Namespace](./docs/concepts/glossary.md)** - corresponds to a federated Kubernetes namespace in edge (physical robot), contains a physical member of fleet
- **[Cloud Infrastructure (AWS)](./docs/concepts/glossary.md)** - corresponds common cloud infrastructure on AWS
- **[Physical Robot](./docs/concepts/glossary.md)** - corresponds a physical robot and it's components
- **[Robot Infrastructure](./docs/concepts/glossary.md)** - corresponds common robot infrastructure

# Architectural Diagram

<div align="center">
    <img src="https://raw.githubusercontent.com/robolaunch/trademark/main/repository-media/robolaunch/architecture.png" width="100%"/>
</div>

For more detail, see [concepts](./docs/concepts/README.md).

# Getting Started

Submodules of [robolaunch Cloud Robotics Platform](https://github.com/robolaunch) are as follows:

- [robolaunch Kubernetes Infrastructure](./docs/setup/robolaunch/kubernetes-infrastructure/README.md)
- [robolaunch Kubernetes Robot Operator](./docs/setup/robolaunch/kubernetes-operators/robot-operator.md)
- [robolaunch Management Suite](./docs/setup/robolaunch/management-suite/README.md)

There are two deployment options available: [Cloud Based Deployment](./docs/setup/cloud-based.md) and [On Premise Deployment](./docs/setup/on-premise.md).

## Prerequisites

The prerequisites for each deployment option are specified separately. In each option, you can deploy single-node or multi-node Kubernetes clusters. 


#### [Cloud Based Deployment](./docs/setup/cloud-based.md)
- AWS account


#### [On Premise Deployment](./docs/setup/on-premise.md)
- Laptop, virtual machine or physical server
- MiniKube or Kubernetes

## Contact
Contact us from info@robolaunch.io for general inquiries.

## Contributing
Thank you for your interest in robolaunch and the desire to contribute! Please check out [organization repositories'](https://github.com/orgs/robolaunch/repositories) contribution guides to learn about conventions to make your changes compatible to our style.

## Community

- [Twitter](https://twitter.com/robolaunch)
- [Slack]() - *soon*
- [Discord]() - *soon*
- [robolaunch Forum]() - *soon*
