group "default" {
    targets = ["compilers", "standalone", "mq", "worker", "manager"]
}

variable "TAG" {
    default = "1.0.1"
}

target "settings" {
    platforms = ["linux/amd64", "linux/arm64"]
}

target "compilers" {
    inherits = ["settings"]
    context = "."
    dockerfile = "docker/Compilers"
    tags = [
		"msc24x/showdown:${TAG}-compilers",
		"msc24x/showdown:latest-compilers"
	]
}

target "standalone" {
    inherits = ["settings"]
    context = "."	
	contexts = {
        "compilers-base" = "target:compilers"
    }
    dockerfile = "docker/Standalone"
    tags = [
		"msc24x/showdown:${TAG}-standalone",
		"msc24x/showdown:latest-standalone"
	]
}

target "mq" {
    inherits = ["settings"]
    context = "."
    dockerfile = "docker/Queue"
    tags = [
		"msc24x/showdown:${TAG}-queue",
		"msc24x/showdown:latest-queue"
	]
}

target "worker" {
    inherits = ["settings"]
    context = "."
	contexts = {
        "compilers-base" = "target:compilers"
    }
    dockerfile = "docker/Worker"
    tags = [
		"msc24x/showdown:${TAG}-worker",
		"msc24x/showdown:latest-worker"
	]
}

target "manager" {
    inherits = ["settings"]
    context = "."
    dockerfile = "docker/Manager"
    tags = [
		"msc24x/showdown:${TAG}-manager",
		"msc24x/showdown:latest-manager"
	]
}
