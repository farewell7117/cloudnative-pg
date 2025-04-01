#
# Copyright The CloudNativePG Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

variable "environment" {
  default = "production"
}

variable "registry" {
  default = "cr.yandex/crpdjkel3ptrc8hlqvfn"
}

variable "insecure" {
  default = "false"
}

variable "latest" {
  default = "false"
}

variable "tag" {
  default = "dev"
}

variable "buildVersion" {
  default = "dev"
}

variable "revision" {
  default = ""
}

suffix = (environment == "testing") ? "-testing" : ""

title = "CloudNativePG Operator"
description = "This Docker image contains CloudNativePG Operator."
authors = "The CloudNativePG Contributors"
url = "https://github.com/cloudnative-pg/cloudnative-pg"
documentation = "https://cloudnative-pg.io/documentation/current/"
license = "Apache-2.0"
now = timestamp()

distros = {
  distroless = {
    baseImage = "gcr.io/distroless/static-debian12:nonroot@sha256:6ec5aa99dc335666e79dc64e4a6c8b89c33a543a1967f20d360922a80dd21f02",
    tag = ""
  }
  ubi = {
    baseImage = "registry.access.redhat.com/ubi9/ubi-micro:latest@sha256:7e85855f6925e03f91b5c51f07886ff1c18c6ec69b5fc65491428a899da914a2",
    tag = "-ubi9"
  }
}

target "default" {
  platforms = ["linux/amd64", "linux/arm64"]
  tags = [
    "cr.yandex/crpdjkel3ptrc8hlqvfn/cloudnative-pg-local:${tag}"
  ]

  dockerfile = "Dockerfile"

  context = "."

  output = [
    "type=image,registry.insecure=${insecure}",
  ]
}

target "distroless" {
  target = "distroless"
  inherits = ["default"]

  args = {
    BASE = "${distros["distroless"].baseImage}"
  }

  output = [
    "type=image,registry.insecure=${insecure}",
  ]

  labels = {
    "org.opencontainers.image.created" = "${now}",
    "org.opencontainers.image.url" = "${url}",
    "org.opencontainers.image.source" = "${url}",
    "org.opencontainers.image.version" = "${buildVersion}",
    "org.opencontainers.image.revision" = "${revision}",
    "org.opencontainers.image.vendor" = "${authors}",
    "org.opencontainers.image.title" = "${title}",
    "org.opencontainers.image.description" = "${description}",
    "org.opencontainers.image.documentation" = "${documentation}",
    "org.opencontainers.image.authors" = "${authors}",
    "org.opencontainers.image.licenses" = "${license}",
    "org.opencontainers.image.base.name" = "${distros["distroless"].baseImage}",
    "org.opencontainers.image.base.digest" = "${digest(distros["distroless"].baseImage)}",
    "name" = "${title}",
    "maintainer" = "${authors}",
    "vendor" = "${authors}",
    "version" = "${buildVersion}",
    "release" = "1",
    "description" = "${description}",
    "summary" = "${description}",
  }
}

target "ubi" {
  target = "ubi"
  inherits = ["default"]

  args = {
    BASE = "${distros["ubi"].baseImage}"
  }

  labels = {
    "org.opencontainers.image.url" = "${url}",
    "org.opencontainers.image.source" = "${url}",
    "org.opencontainers.image.version" = "${buildVersion}",
    "org.opencontainers.image.revision" = "${revision}",
    "org.opencontainers.image.vendor" = "${authors}",
    "org.opencontainers.image.title" = "${title}",
    "org.opencontainers.image.description" = "${description}",
    "org.opencontainers.image.documentation" = "${documentation}",
    "org.opencontainers.image.authors" = "${authors}",
    "org.opencontainers.image.licenses" = "${license}",
    "org.opencontainers.image.base.name" = "${distros["ubi"].baseImage}",
    "org.opencontainers.image.base.digest" = "${digest(distros["ubi"].baseImage)}",
    "name" = "${title}",
    "maintainer" = "${authors}",
    "vendor" = "${authors}",
    "version" = "${buildVersion}",
    "release" = "1",
    "description" = "${description}",
    "summary" = "${description}",
  }
}

function digest {
  params = [ imageNameWithSha ]
  result = index(split("@", imageNameWithSha), 1)
}

function latest {
  params = [ image, latest ]
  result = (latest == "true") ? "${image}:latest" : ""
}
