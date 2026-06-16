pulumi {
  required_providers {
    docker-build = {
      source  = "pulumi/docker-build"
      version = "0.1.0-alpha.0+dev"
    }
  }
}

resource "docker-build_image" "multiPlatform" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.multiPlatform"
  }
  context = {
    location = "./app"
  }
  platforms = ["plan9/amd64", "plan9/386"]
}
resource "docker-build_image" "registryPush" {
  push = false
  context = {
    location = "./app"
  }
  tags = ["docker.io/pulumibot/buildkit-e2e:example"]
  exports {
    registry = {
      oci_media_types = true
      push            = false
    }
  }
  registries {
    address  = "docker.io"
    username = "pulumibot"
    password = var.dockerHubPassword
  }
}
resource "docker-build_image" "cached" {
  push = false
  context = {
    location = "./app"
  }
  cache_to {
    local = {
      dest = "tmp/cache"
      mode = "max"
    }
  }
  cache_from {
    local = {
      src = "tmp/cache"
    }
  }
}
resource "docker-build_image" "buildArgs" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.buildArgs"
  }
  context = {
    location = "./app"
  }
  build_args = {
    "SET_ME_TO_TRUE" = "true"
  }
}
resource "docker-build_image" "extraHosts" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.extraHosts"
  }
  context = {
    location = "./app"
  }
  add_hosts = ["metadata.google.internal:169.254.169.254"]
}
resource "docker-build_image" "sshMount" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.sshMount"
  }
  context = {
    location = "./app"
  }
  ssh {
    id = "default"
  }
}
resource "docker-build_image" "secrets" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.secrets"
  }
  context = {
    location = "./app"
  }
  secrets = {
    "password" = "hunter2"
  }
}
resource "docker-build_image" "labels" {
  push = false
  context = {
    location = "./app"
  }
  labels = {
    "description" = "This image will get a descriptive label 👍"
  }
}
resource "docker-build_image" "target" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.target"
  }
  context = {
    location = "./app"
  }
  target = "build-me"
}
resource "docker-build_image" "namedContexts" {
  push = false
  dockerfile = {
    location = "./app/Dockerfile.namedContexts"
  }
  context = {
    location = "./app"
    named = {
      "golang:latest" = {
        location = "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984"
      }
    }
  }
}
resource "docker-build_image" "remoteContext" {
  push = false
  context = {
    location = "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile"
  }
}
resource "docker-build_image" "remoteContextWithInline" {
  push = false
  dockerfile = {
    inline = "FROM busybox\nCOPY hello.c ./\n"
  }
  context = {
    location = "https://github.com/docker-library/hello-world.git"
  }
}
resource "docker-build_image" "inline" {
  push = false
  dockerfile = {
    inline = "FROM alpine\nRUN echo \"This uses an inline Dockerfile! 👍\"\n"
  }
}
resource "docker-build_image" "dockerLoad" {
  push = false
  context = {
    location = "./app"
  }
  exports {
    docker = {
      tar = true
    }
  }
}
variable "dockerHubPassword" {
  type = string
}
output "platforms" {
  value = docker-build_image.multiPlatform.platforms
}
