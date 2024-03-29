BUILD_VERSION=$(shell cat config.json | awk 'BEGIN { FS="\""; RS="," }; { if ($$2 == "version") {print $$4} }')

release: release_armhf release_aarch64 release_i386 release_amd64 release_armv7

release_armhf:
	@docker build -t danielapatin/external-mqtt-to-local-armhf:${BUILD_VERSION} --build-arg BUILD_ARCH=armhf --build-arg BUILD_VERSION=${BUILD_VERSION} .
	@docker push danielapatin/external-mqtt-to-local-armhf:${BUILD_VERSION}

release_aarch64:
	@docker build -t danielapatin/external-mqtt-to-local-aarch64:${BUILD_VERSION} --build-arg BUILD_ARCH=aarch64 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	@docker push danielapatin/external-mqtt-to-local-aarch64:${BUILD_VERSION}

release_i386:
	@docker build -t danielapatin/external-mqtt-to-local-i386:${BUILD_VERSION} --build-arg BUILD_ARCH=i386 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	@docker push danielapatin/external-mqtt-to-local-i386:${BUILD_VERSION}

release_amd64:
	@docker build -t danielapatin/external-mqtt-to-local-amd64:${BUILD_VERSION} --build-arg BUILD_ARCH=amd64 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	@docker push danielapatin/external-mqtt-to-local-amd64:${BUILD_VERSION}

release_armv7:
	@docker build -t danielapatin/external-mqtt-to-local-armv7:${BUILD_VERSION} --build-arg BUILD_ARCH=armv7 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	@docker push danielapatin/external-mqtt-to-local-armv7:${BUILD_VERSION}
