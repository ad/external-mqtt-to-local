BUILD_VERSION=0.0.3

release: release_armhf release_aarch64 release_i386 release_amd64 release_armv7

release_armhf:
	docker build -t danielapatin/external-mqtt-to-local-armhf --build-arg BUILD_ARCH=armhf --build-arg BUILD_VERSION=${BUILD_VERSION} .
	docker push danielapatin/external-mqtt-to-local-armhf

release_aarch64:
	docker build -t danielapatin/external-mqtt-to-local-aarch64 --build-arg BUILD_ARCH=aarch64 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	docker push danielapatin/external-mqtt-to-local-aarch64

release_i386:
	docker build -t danielapatin/external-mqtt-to-local-i386 --build-arg BUILD_ARCH=i386 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	docker push danielapatin/external-mqtt-to-local-i386

release_amd64:
	docker build -t danielapatin/external-mqtt-to-local-amd64 --build-arg BUILD_ARCH=amd64 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	docker push danielapatin/external-mqtt-to-local-amd64

release_armv7:
	docker build -t danielapatin/external-mqtt-to-local-armv7 --build-arg BUILD_ARCH=armv7 --build-arg BUILD_VERSION=${BUILD_VERSION} .
	docker push danielapatin/external-mqtt-to-local-armv7
