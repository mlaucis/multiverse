#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

INSTANCE_ID=`wget -qO- http://instance-data/latest/meta-data/instance-id`
REGION=`wget -qO- http://instance-data/latest/meta-data/placement/availability-zone | sed 's/.$//'`
INSTANCE_TAGS="{"`aws ec2 describe-tags --region ${REGION} --filter "Name=resource-id,Values=$INSTANCE_ID" --output=text | sed -r 's/TAGS\t(.*)\t.*\t.*\t(.*)/"\1": "\2",/g'`"\"blank_tag\":\"blank_tag\"}"
DISTRIBUTOR_DEPLOY_TARGET=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["distributor_target"]'`
INSTALLER_CHANNEL=`echo ${INSTANCE_TAGS} | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["installer_channel"]'`

logger -t distributor_installer got INSTANCE_ID: ${INSTANCE_ID}
logger -t distributor_installer got REGION: ${REGION}
logger -t distributor_installer got INSTANCE_TAGS: ${INSTANCE_TAGS}
logger -t distributor_installer got INSTALLER_CAHNNEL: ${INSTALLER_CHANNEL}
logger -t distributor_installer got DISTRIBUTOR_DEPLOY_TARGET: ${DISTRIBUTOR_DEPLOY_TARGET}

if [ -z "${DISTRIBUTOR_DEPLOY_TARGET}" ]; then
    logger -t distributor_installer installer target not found
    exit 1
fi

mkdir -p ~/releases/distributor/${DISTRIBUTOR_DEPLOY_TARGET}
cd ~/releases/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/

logger -t distributor_installer installing certificate
echo '-----BEGIN CERTIFICATE-----
MIIFcjCCA1oCCQCTAAojbwACaDANBgkqhkiG9w0BAQsFADB7MQswCQYDVQQGEwJE
RTEPMA0GA1UECAwGQmVybGluMQ8wDQYDVQQHDAZCZXJsaW4xEDAOBgNVBAoMB1Rh
cGdsdWUxFjAUBgNVBAMMDSoudGFwZ2x1ZS5jb20xIDAeBgkqhkiG9w0BCQEWEXRv
b2xzQHRhcGdsdWUuY29tMB4XDTE1MDkxMjIxNTAwMloXDTI1MDkwOTIxNTAwMlow
ezELMAkGA1UEBhMCREUxDzANBgNVBAgMBkJlcmxpbjEPMA0GA1UEBwwGQmVybGlu
MRAwDgYDVQQKDAdUYXBnbHVlMRYwFAYDVQQDDA0qLnRhcGdsdWUuY29tMSAwHgYJ
KoZIhvcNAQkBFhF0b29sc0B0YXBnbHVlLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAO5vH5wwW5K9hTij6rzPYFZ8UXvpuNmYRugIbCvKnUhQGvBn
a4Xff/SU+lDm7EHFq9k5K8dv0dw4hRl2MPsTYAM7w1ycxgdYfBdxfp3BGzsaXGVy
uGHAQPdVrwAXf1fxBAfCYvwAAOzAGaivfTPRfYBba2Vv74u35zPGS4/0jyGHGHFE
00r7/MfxX4ahx0WATyeJLpXUesxcCS1B14smPZKfm82PQktJmRdsQO2fe7n3PT/x
NnoGuKs+wHlnz0rIaDISUXJ/ggyoO31wRalEL3i3w9kskWrPO1N9TSQk9AA7pEFi
iIFHP4pk9mzioT2UmXmyZAiaOcBm9ZXxUEv6yOoMu4EPMxtjHd9ti6oqynz8luM9
DK8t7VQ5rHx5eY1PcjHyuPSfWiQrLYyPSUcgCdCrLfADodMMakJzt32lwk/Kt+n0
RkLv4c6XeAvsfvxu6AQCm9CE++GEoeDaP3nnVkJWJL3wYGeEKEmV0nemb6P31HAw
kVpsmpdKYIs0OutCdjvcQbVAKibgZICk9fLxGSYRm5tjdP8njAJYmvdNlAZb4c19
Zumfc3pq/aZnQVxD2pVKjrlFioHsNmvWdj42WPpxNA+yJkHo0mTsD2umFBUKr2z6
EkQElSC4k0MU5mSAzJ3xmgRBP3DoHKI2TGN0mcWJoICsEMBJOEQD9bM0wrN5AgMB
AAEwDQYJKoZIhvcNAQELBQADggIBABWAkm/sXM6FGo6w5p8mHHaR7XwVl+dPdzGp
0dxfILXJQKCeXcvD3C6ChfHkZgGnNqX80Jq8b6mwSP702y9J6VnB+HXiI4FVmufK
m3o6Ep/Ca0uvh3ku+KYvCSakeL2S83RRocTU2wsuOldnuYq5KPjDWexv/HtKwAw9
MuJtUigh94tzljBG6VChdjOTxSNS6jwoBq2ndX77zkh2Qfh3OGPaUgA6GkIN5nEB
VUZfdHu9n3d0SoSRyk49lpYryXyqYipCm3GvKlEX34EmoR6MLgNcPS4zF/VimODC
zTwqgQNeK7CqFJJrql0ShUB2vW7LFoYPA/B2YxB3vPdu+TcX0y4WuJHX6vSI2iKq
hF2509iQAEQYd5ucZbG/s6Oc06xSP5Z+5OLjstHiIu5Z0rtJmhbKSOCt/g/x3BnM
BzRZjSnkp3GTWCXaV3o16S5RoZ1QHzlMIU8DGuDUY5o8rWmSIv5HYX0bW7n5LYqI
vb0Y0Q/z/AA4YOO2jEzxxgKxk9Qf8l3Mp6Qwyvc5nPaxhTDFjPb3/Vl7nsYk0UIk
FSYfZNgtGf9G5+LS6vs8QesLkC1MZivGn4uGAwOKosvsJlM184+8MhYER2b+m7MI
vbPtYg32n4BkRMUEsGy5G8O1oMJWtfSF5/cb9mVmNvqKPIgVxBirLxaUPlYWafLp
BmNddyWm
-----END CERTIFICATE-----' > ./self.crt

logger -t distributor_installer installing private key
echo '-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----' > ./self.key

logger -t distributor_installer installing cloudflare certificate
echo '-----BEGIN CERTIFICATE-----
MIIGBjCCA/CgAwIBAgIIV5G6lVbCLmEwCwYJKoZIhvcNAQENMIGQMQswCQYDVQQG
EwJVUzEZMBcGA1UEChMQQ2xvdWRGbGFyZSwgSW5jLjEUMBIGA1UECxMLT3JpZ2lu
IFB1bGwxFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xEzARBgNVBAgTCkNhbGlmb3Ju
aWExIzAhBgNVBAMTGm9yaWdpbi1wdWxsLmNsb3VkZmxhcmUubmV0MB4XDTE1MDEx
MzAyNDc1M1oXDTIwMDExMjAyNTI1M1owgZAxCzAJBgNVBAYTAlVTMRkwFwYDVQQK
ExBDbG91ZEZsYXJlLCBJbmMuMRQwEgYDVQQLEwtPcmlnaW4gUHVsbDEWMBQGA1UE
BxMNU2FuIEZyYW5jaXNjbzETMBEGA1UECBMKQ2FsaWZvcm5pYTEjMCEGA1UEAxMa
b3JpZ2luLXB1bGwuY2xvdWRmbGFyZS5uZXQwggIiMA0GCSqGSIb3DQEBAQUAA4IC
DwAwggIKAoICAQDdsts6I2H5dGyn4adACQRXlfo0KmwsN7B5rxD8C5qgy6spyONr
WV0ecvdeGQfWa8Gy/yuTuOnsXfy7oyZ1dm93c3Mea7YkM7KNMc5Y6m520E9tHooc
f1qxeDpGSsnWc7HWibFgD7qZQx+T+yfNqt63vPI0HYBOYao6hWd3JQhu5caAcIS2
ms5tzSSZVH83ZPe6Lkb5xRgLl3eXEFcfI2DjnlOtLFqpjHuEB3Tr6agfdWyaGEEi
lRY1IB3k6TfLTaSiX2/SyJ96bp92wvTSjR7USjDV9ypf7AD6u6vwJZ3bwNisNw5L
ptph0FBnc1R6nDoHmvQRoyytoe0rl/d801i9Nru/fXa+l5K2nf1koR3IX440Z2i9
+Z4iVA69NmCbT4MVjm7K3zlOtwfI7i1KYVv+ATo4ycgBuZfY9f/2lBhIv7BHuZal
b9D+/EK8aMUfjDF4icEGm+RQfExv2nOpkR4BfQppF/dLmkYfjgtO1403X0ihkT6T
PYQdmYS6Jf53/KpqC3aA+R7zg2birtvprinlR14MNvwOsDOzsK4p8WYsgZOR4Qr2
gAx+z2aVOs/87+TVOR0r14irQsxbg7uP2X4t+EXx13glHxwG+CnzUVycDLMVGvuG
aUgF9hukZxlOZnrl6VOf1fg0Caf3uvV8smOkVw6DMsGhBZSJVwao0UQNqQIDAQAB
o2YwZDAOBgNVHQ8BAf8EBAMCAAYwEgYDVR0TAQH/BAgwBgEB/wIBAjAdBgNVHQ4E
FgQUQ1lLK2mLgOERM2pXzVc42p59xeswHwYDVR0jBBgwFoAUQ1lLK2mLgOERM2pX
zVc42p59xeswCwYJKoZIhvcNAQENA4ICAQDKDQM1qPRVP/4Gltz0D6OU6xezFBKr
LWtDoA1qW2F7pkiYawCP9MrDPDJsHy7dx+xw3bBZxOsK5PA/T7p1dqpEl6i8F692
g//EuYOifLYw3ySPe3LRNhvPl/1f6Sn862VhPvLa8aQAAwR9e/CZvlY3fj+6G5ik
3it7fikmKUsVnugNOkjmwI3hZqXfJNc7AtHDFw0mEOV0dSeAPTo95N9cxBbm9PKv
qAEmTEXp2trQ/RjJ/AomJyfA1BQjsD0j++DI3a9/BbDwWmr1lJciKxiNKaa0BRLB
dKMrYQD+PkPNCgEuojT+paLKRrMyFUzHSG1doYm46NE9/WARTh3sFUp1B7HZSBqA
kHleoB/vQ/mDuW9C3/8Jk2uRUdZxR+LoNZItuOjU8oTy6zpN1+GgSj7bHjiy9rfA
F+ehdrz+IOh80WIiqs763PGoaYUyzxLvVowLWNoxVVoc9G+PqFKqD988XlipHVB6
Bz+1CD4D/bWrs3cC9+kk/jFmrrAymZlkFX8tDb5aXASSLJjUjcptci9SKqtI2h0J
wUGkD7+bQAr+7vr8/R+CBmNMe7csE8NeEX6lVMF7Dh0a1YKQa6hUN18bBuYgTMuT
QzMmZpRpIBB321ZBlcnlxiTJvWxvbCPHKHj20VwwAz7LONF59s84ZsOqfoBv8gKM
s0s5dsq5zpLeaw==
-----END CERTIFICATE-----' > ./origin-pull-ca.pem

aws s3 cp s3://tapglue-builds/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/releases.json ./

releaseVersion=`cat ./releases.json | python -c 'import sys,json;data=json.loads(sys.stdin.read()); print data["current_release"]'`
execName=distributor_${DISTRIBUTOR_DEPLOY_TARGET}_${releaseVersion}

aws s3 cp s3://tapglue-builds/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/distributor_${DISTRIBUTOR_DEPLOY_TARGET}.${releaseVersion}.tar.gz ./
aws s3 cp s3://tapglue-builds/distributor/${DISTRIBUTOR_DEPLOY_TARGET}/config.json ./

tar -zxvf distributor_${DISTRIBUTOR_DEPLOY_TARGET}.${releaseVersion}.tar.gz

echo '#!/usr/bin/env bash

exec 1> >(logger -t $(basename ${0})) 2>&1

logger -t distributor runner received run command for distributor ${1} version ${2}

./distributor_${1}_${2} -target '${DISTRIBUTOR_DEPLOY_TARGET}'
' > run.sh

chmod +x ./run.sh

logger -t distributor_installer deployed distributor_${DISTRIBUTOR_DEPLOY_TARGET}_${releaseVersion}

rm -f distributor_${DISTRIBUTOR_DEPLOY_TARGET}.${releaseVersion}.tar.gz

kill -9 `ps aux | grep distributor_${DISTRIBUTOR_DEPLOY_TARGET}_ | awk -F" " '{print $2}'`

nohup ./run.sh ${DISTRIBUTOR_DEPLOY_TARGET} ${releaseVersion} &
