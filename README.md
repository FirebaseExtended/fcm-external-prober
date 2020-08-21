# FCM External Prober

## About

FCM External Prober is a tool that can be used to monitor the performance of FCM by way of repeatedly sending messages to an emulated Android app, which logs its reception of the messages, and calculating the time it took to receive the message, providing measurements of both availability and latency.

## How to Run:

In the `Controller/src` directory, call `go run main.go -config="<configPath>"` where `configPath` is the path to your configuration file.

## How to Stop:

This program can be teriminated using `^C`. If invoked before VMs are created, the program will terminate normally. If invoked after VMs are created, the prober will allow any outstanding probes to be resolved, and will then delete any regional VMs created during its runtime.

## Requirements to Run:

### Google Cloud Platform

This tool is designed to run on GCP. It utilizes Compute Engine VM instances, the gcloud CLI tool, IAM service accounts, and exports logs to Cloud Logger. Thus, a GCP account is required.

#### Compute Engine

This tool requires a VM instance on which to run the controller, and an image from which to generate the VM instances on which probes will run.

#### Image

The image for regional VMs should have the following dependencies installed and their executables available via PATH:
* Android SDK
    * `emulator`
    * `sdkmanager`
    * `avdmanager`
    * `adb` (Android Debug Bridge)
* Java JDK
* Golang
* Git
* Protocol Buffer compiler with gRPC plugin installed

In addition, nested virtualization needs to be enabled on this image in order to run the Android app. A guide for how to enable this can be found [here](https://cloud.google.com/compute/docs/instances/enable-nested-virtualization-vm-instances).

If you want to use FCMExternalProberTarget as your logging app, the package name will need to be changed in order to register the app with FCM, and the google-services.json file should be located on the image or retrieved by the startup script.
In order to install your custom logging app on the emulated device, an apk of the logging app should be present on the image, or can be retrieved from an external source with the startup script.

If you are planning to run this tool with multiple regional VMs, you must generate a new Android Virtual Device during the execution of the startup script. If an AVD saved to an image is utilized by multiple regional VMs, it will create a conflict when registering devices with FCM when the app starts.

#### Startup Script

A startup script is used to install dependencies and begin probing once a regional VM has been created. A generic startup script is included in this repository, but the script is ultimately app- and configuration-specific, so changes are likely needed in order for this tool to function properly. The general requirements for the startup script are as follows:
* Acquire the most recent version of this repository
* Acquire a version of a logging app to install on the emulated device
* Create a new AVD with a desired version of Android
* Compile Protocol Buffers in `Controller/src/controller`
* Initiate probing by running `main.go` in `Probe/src`

#### IAM Service Account

This tool utilizes service accounts to handle authentication to other GCP services. Thus, a service account needs to be created with the following permissions:
* Compute Engine - Compute Admin
* Logging - Log writer
* Firebase Products - Firebase Cloud Messaging Admin

## Operating Costs

The operating cost of the prober is correlated to the number of regions in which probes are designated. Emulating the Android app is relatively expensive, so an n1-standard-4 machine type is used. A VM with this configuration costs roughly $100.00 to run on a per-monthly basis. The VM on which the controller runs needs fewer resources, so an n1-standard-2 machine would likely suffice.


## Source Code Headers

Every file containing source code must include copyright and license
information. This includes any JS/CSS files that you might be serving out to
browsers. (This is to help well-intentioned people avoid accidental copying that
doesn't comply with the license.)

Apache header:

    Copyright 2020 Google LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
